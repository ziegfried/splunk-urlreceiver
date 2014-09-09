package main

import . "modinputs"
import "fmt"
import "encoding/xml"
import "os"
import "os/signal"
import "net/http"
import "strconv"
import "io/ioutil"
import "strings"

const DEBUG = false

func confiureInputHandler(stanza InputStanza, stream Stream, servers map[int]http.ServeMux) {
	port, err := strconv.Atoi(stanza.GetParamValue("port", "80"))
	if err != nil {
		Fatal("Invalid port number specified in stanza=%s: %s", stanza.Name, err)
		return
	}
	Info("Running webserver for stanza %s on port %d", stanza.Name, port)

	readBody := NormalizeBoolean(stanza.GetParamValue("read_body", "true"), true)
	readFormField := stanza.GetParamValue("form_field", "")
	useClientip := NormalizeBoolean(stanza.GetParamValue("host_from_clientip", "false"), false)

	server := getServer(port, servers)

	server.Handle(stanza.GetParamValue("path", "/"), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if DEBUG {
			Info("Request: %s %s from clientip=%s", r.Method, r.RequestURI, r.RemoteAddr)
		}
		event := stanza.MakeEvent()
		if useClientip {
			event.Host = strings.Split(r.RemoteAddr, ":")[0]
		}
		if readBody {
			Debug("Reading request body")
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				Warn("Error reading request body: %s", err)
				return
			}
			if len(body) > 0 {
				event.Data = string(body)
			} else {
				Warn("Received empty request body")
				http.Error(w, "Empty request body", 400)
				return
			}
		} else if readFormField != "" {
			Debug("Reading form field=%s", readFormField)
			value := r.FormValue(readFormField)
			if value != "" {

				event.Data = value
				stream.Send(event)
			} else {
				Warn("Received empty form field=%s", readFormField)
				http.Error(w, fmt.Sprintf("Missing form field %s", readFormField), 400)
				return
			}
		} else {
			http.Error(w, "No form field configured for this handler", 500)
			return
		}
		stream.Send(event)
		fmt.Fprint(w, stanza.GetParamValue("response", "OK"))
	}))
}

func getServer(port int, servers map[int]http.ServeMux) http.ServeMux {
	server, found := servers[port]
	if !found {
		server = *http.NewServeMux()
		servers[port] = server
	}
	return server
}

func startListening(servers map[int]http.ServeMux) {
	for port, server := range servers {
		go http.ListenAndServe(fmt.Sprintf(":%d", port), &server)
	}
}

var scheme = Scheme{
	Title:                 "URL-Receiver",
	Description:           "Receive and index data via URLs (such as Webhooks)",
	UseExternalValidation: false,
	StreamingMode:         "xml",
	UseSingleInstance:     true,
	EndpointArguments: []EndpointArg{
		EndpointArg{
			Name:             "port",
			Title:            "TCP Port",
			DataType:         "integer",
			Validation:       "is_avail_tcp_port('port')",
			Description:      "TCP port to listen for HTTP requests",
			RequiredOnCreate: true,
			RequiredOnEdit:   false,
		},
		EndpointArg{
			Name:             "path",
			Title:            "URL Path",
			Description:      "URL path to match against when receiving data",
			RequiredOnCreate: true,
			RequiredOnEdit:   false,
		},
		EndpointArg{
			Name:        "read_body",
			Title:       "Read request body",
			Description: "Index the request body of incoming HTTP requests",
			DataType:    "boolean",
		},
		EndpointArg{
			Name:        "form_field",
			Description: "Specify the form field to index (only if read_body is disabled)",
		},
		EndpointArg{
			Name:        "host_from_clientip",
			Title:       "Use client IP address for host field value",
			Description: "If enabled the input will pass the IP address of the client sending the HTTP request as the host field value for each event",
			DataType:    "boolean",
		},
	},
}

func printScheme() {
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "  ")
	if err := enc.Encode(scheme); err != nil {
		fmt.Printf("error: %v\n", err)
	}
	fmt.Println()
}

func main() {
	if len(os.Args) > 1 {
		arg := os.Args[1]
		switch arg {
		case "--scheme":
			printScheme()
			os.Exit(0)
		case "--validate-arguments":
			Fatal("Arg validation not implemented")
			os.Exit(1)
		default:
			Fatal("Unsupported command line flag: %s", os.Args[1])
			os.Exit(2)
		}
	}

	config := InputConfig{}
	if err := xml.NewDecoder(os.Stdin).Decode(&config); err != nil {
		Fatal("Error reading configuration: %s", err)
		os.Exit(3)
	}

	stream := NewStream(1000)
	stream.Start()

	servers := make(map[int]http.ServeMux, 100)

	for _, stanza := range config.Stanzas {
		confiureInputHandler(stanza, stream, servers)
	}
	startListening(servers)

	waitForSignal()

	// Wait until stream is stopped
	<-stream.Stop()
}

func waitForSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	if DEBUG {
		Debug("Received signal: %s", sig)
	}
}

func NormalizeBoolean(val string, def bool) bool {
	switch strings.ToLower(strings.TrimSpace(val)) {
	case "true", "yes", "on", "1":
		return true
	case "false", "no", "off", "0":
		return false
	default:
		return def
	}
}
