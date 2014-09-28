package main

import (
	"encoding/xml"
	. "formatter"
	. "modinputs"
	"net/http"
	"os"
	"os/signal"
	. "scheme"
	"server"
	"strconv"
	"strings"
	. "utils"
	. "verifier"
)

func main() {
	if len(os.Args) > 1 {
		arg := os.Args[1]
		switch arg {
		case "--scheme":
			PrintScheme(SCHEME)
			os.Exit(0)
		case "--spec":
			PrintSpec(SCHEME, "urlreceiver")
			os.Exit(0)
		case "--validate-arguments":
			Log.Fatal("Arg validation not implemented")
			os.Exit(1)
		case "--debug":
			Log.Level = LEVEL_DEBUG
			Log.Debug("Enabled debug logging")
		default:
			Log.Fatal("Unsupported command line flag: %s", os.Args[1])
			os.Exit(2)
		}
	}

	config := InputConfig{}
	if err := xml.NewDecoder(os.Stdin).Decode(&config); err != nil {
		Log.Fatal("Error reading configuration: %s", err)
		os.Exit(3)
	}

	stream := NewStream(1000)
	stream.Start()

	for _, stanza := range config.Stanzas {
		confiureInputHandler(stanza, config, stream)
	}
	server.StartListening()

	waitForSignal()
	// Wait until stream is stopped
	<-stream.Stop()
}

func confiureInputHandler(stanza InputStanza, config InputConfig, stream Stream) {
	Log.Debug("Configuring input=%s", stanza.Name)
	debug := NormalizeBoolean(stanza.GetParamValue(PARAM_DEBUG, "false"), false)

	port, err := strconv.Atoi(stanza.GetParamValue(PARAM_PORT, "80"))
	if err != nil || port < 1 {
		Log.Error("[%s] Invalid port number specified: %s", stanza.Name, err)
		return
	}
	path := stanza.GetParamValue("path", "/")

	Log.Info("[%s] Accepting requests on port=%d and path=%s", stanza.Name, port, path)

	verifier := GetVerifier(stanza.GetParamValue(PARAM_VERIFY_TYPE, "-"), stanza, debug)
	formatter := GetFormatter(stanza.GetParamValue(PARAM_DATA_FORMAT, "-"), stanza, debug)

	useClientip := NormalizeBoolean(stanza.GetParamValue("host_from_clientip", "true"), true)
	responseText := stanza.GetParamValue("response", "OK")

	server, err := server.Get(port, stanza, config)
	if err != nil {
		Log.Error("[%s] Error configuring webserver: %s", stanza.Name, err)
		return
	}
	server.Register(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if debug {
			Log.Info("[%s] Handling request: %s %s from clientip=%s", stanza.Name, r.Method, r.RequestURI, r.RemoteAddr)
		}

		if verifier != nil {
			switch verifier(w, r) {
			case VERIFY_OK:
				if debug {
					Log.Info("[%s] Verfier %s successfully passed", stanza.Name, stanza.GetParamValue(PARAM_VERIFY_TYPE, "-"))
				}

			case VERIFY_HANDLED:
				if debug {
					Log.Info("[%s] Request has been handled by verifier", stanza.Name)
				}
				return
			case VERIFY_ERROR:
				fallthrough
			default:
				Log.Warn("[%s] Unable to verify request, responding with status 400", stanza.Name)
				http.Error(w, "Verification error", 400)
				return
			}
		}

		event := stanza.MakeEvent()
		if useClientip {
			event.Host = strings.Split(r.RemoteAddr, ":")[0]
		}

		data, err := formatter(r, stanza, debug)

		if err != nil {
			http.Error(w, err.Error(), 400)
		} else if data == "" {
			Log.Warn("Received empty data")
			http.Error(w, "Empty data", 400)
		} else {
			event.Data = data
			stream.Send(event)
			http.Error(w, responseText, 200)
		}
	}))
}

func waitForSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	Log.Debug("Received signal: %s", sig)
}
