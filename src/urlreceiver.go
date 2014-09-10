package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	. "modinputs"
	"net/http"
	"os"
	"os/signal"
	"server"
	"strconv"
	"strings"
	. "utils"
)

const (
	DEBUG = true

	PARAM_PORT               = "port"
	PARAM_PATH               = "path"
	PARAM_DATA_RETRIEVAL     = "data_retrieval"
	PARAM_FORM_FIELD         = "form_field"
	PARAM_HOST_FROM_CLIENTIP = "host_from_clientip"
	PARAM_RESPONSE           = "response"
	PARAM_DEBUG              = "debug"

	DATA_RAW_BODY     = "raw_body"
	DATA_FORM_FIELD   = "form_field"
	DATA_FORM_KV      = "form_kv"
	DATA_FULL_REQUEST = "full_request"
)

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

	for _, stanza := range config.Stanzas {
		confiureInputHandler(stanza, stream)
	}
	server.StartListening()

	waitForSignal()
	// Wait until stream is stopped
	<-stream.Stop()
}

var scheme = Scheme{
	Title:                 "URL-Receiver",
	Description:           "Receive and index data via URLs (such as Webhooks)",
	UseExternalValidation: false,
	StreamingMode:         "xml",
	UseSingleInstance:     true,
	EndpointArguments: []EndpointArg{
		EndpointArg{
			Name:             PARAM_PORT,
			Title:            "TCP Port",
			Validation:       "is_port('port')",
			Description:      "TCP port to listen for HTTP requests",
			RequiredOnCreate: true,
			RequiredOnEdit:   false,
		},
		EndpointArg{
			Name:             PARAM_PATH,
			Title:            "URL Path",
			Description:      "URL path to match against when receiving data",
			RequiredOnCreate: true,
			RequiredOnEdit:   false,
		},
		EndpointArg{
			Name:             PARAM_DATA_RETRIEVAL,
			Title:            "Data retrieval",
			Description:      "Define how data is extracted from the incoming HTTP request",
			RequiredOnCreate: true,
			RequiredOnEdit:   false,
		},
		EndpointArg{
			Name:             PARAM_FORM_FIELD,
			Description:      "Specify the form field to index (only if data retrieval is form_field)",
			RequiredOnCreate: false,
		},
		EndpointArg{
			Name:             PARAM_HOST_FROM_CLIENTIP,
			Title:            "Use client IP address for host field value",
			Description:      "If enabled the input will pass the IP address of the client sending the HTTP request as the host field value for each event",
			DataType:         "boolean",
			RequiredOnCreate: false,
		},
		EndpointArg{
			Name:             PARAM_RESPONSE,
			Title:            "Response text",
			Description:      "Text the webserver responds with after successfully receiving data (Default is \"OK\")",
			RequiredOnCreate: false,
		},
		EndpointArg{
			Name:             PARAM_DEBUG,
			Title:            "Debug logging",
			Description:      "Enable debug logging for this input (logged to splunkd.log)",
			DataType:         "boolean",
			RequiredOnCreate: false,
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

func confiureInputHandler(stanza InputStanza, stream Stream) {
	debug := NormalizeBoolean(stanza.GetParamValue(PARAM_DEBUG, "false"), false)
	port, err := strconv.Atoi(stanza.GetParamValue(PARAM_PORT, "80"))
	if err != nil || port < 1 {
		Fatal("[%s] Invalid port number specified: %s", stanza.Name, err)
		return
	}
	path := stanza.GetParamValue("path", "/")

	if DEBUG {
		Debug("[%s] Accepting requests on port=%d and path=%s", stanza.Name, port, path)
	}
	var formatter Formatter = nil

	switch stanza.GetParamValue("data_retrieval", "n/a") {
	case DATA_RAW_BODY:
		if DEBUG || debug {
			Debug("[%s] Will index raw request body", stanza.Name)
		}
		formatter = readRawBody
	case DATA_FORM_FIELD:
		if DEBUG || debug {
			Debug("[%s] Will index form field %s", stanza.Name, stanza.GetParamValue(PARAM_FORM_FIELD, "data (default)"))
		}
		formatter = readFormField
	case DATA_FORM_KV:
		if DEBUG || debug {
			Debug("[%s] Will read form key-values", stanza.Name)
		}
		formatter = formFieldsKV
	case DATA_FULL_REQUEST:
		if DEBUG || debug {
			Debug("[%s] Will dump full request details", stanza.Name)
		}
		formatter = dumpRequest
	default:
		Warn("[%s] Ignoring unknown data_retrieval type \"%s\" - indexing request body", stanza.Name, stanza.GetParamValue("data_retrieval", ""))
		formatter = readRawBody
	}

	useClientip := NormalizeBoolean(stanza.GetParamValue("host_from_clientip", "true"), true)
	responseText := stanza.GetParamValue("response", "OK")

	server := server.Get(port)
	server.Register(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if DEBUG || debug {
			Debug("[%s] Handling request: %s %s from clientip=%s", stanza.Name, r.Method, r.RequestURI, r.RemoteAddr)
		}
		event := stanza.MakeEvent()
		if useClientip {
			event.Host = strings.Split(r.RemoteAddr, ":")[0]
		}

		data, err := formatter(r, stanza, DEBUG || debug)

		if err != nil {
			http.Error(w, err.Error(), 400)
		} else if data == "" {
			Warn("Received empty data")
			http.Error(w, "Empty data", 400)
		} else {
			event.Data = data
			stream.Send(event)
			http.Error(w, responseText, 200)
		}
	}))
}

type Formatter func(*http.Request, InputStanza, bool) (string, error)

func readRawBody(r *http.Request, stanza InputStanza, debug bool) (string, error) {
	if debug {
		Debug("[%s] Reading request body", stanza.Name)
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Warn("[%s] Error reading request body: %s", stanza.Name, err)
		return "", fmt.Errorf("IO Error")
	}
	if len(body) > 0 {
		return string(body), nil
	} else {
		Warn("[%s] Received empty request body", stanza.Name)
		return "", fmt.Errorf("Empty request body")
	}
}

func readFormField(r *http.Request, stanza InputStanza, debug bool) (string, error) {
	field := stanza.GetParamValue(PARAM_FORM_FIELD, "data")
	if debug {
		Debug("[%s] Reading form field=%s", stanza.Name, field)
	}
	return r.FormValue(field), nil
}

func formFieldsKV(r *http.Request, stanza InputStanza, debug bool) (string, error) {
	if debug {
		Debug("[%s] Reading all form fields", stanza.Name)
	}

	if r.Form == nil {
		r.ParseMultipartForm(32 << 20)
	}

	var buffer bytes.Buffer
	for param, values := range r.Form {
		for _, value := range values {
			buffer.WriteString(param)
			buffer.WriteString("=\"")
			buffer.WriteString(strings.Replace(value, "\"", "\\\"", -1))
			buffer.WriteString("\" ")
		}
	}
	return buffer.String(), nil
}

func dumpRequest(r *http.Request, stanza InputStanza, debug bool) (string, error) {
	if debug {
		Debug("[%s] Dumping request", stanza.Name)
	}

	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("%s %s HTTP/%d.%d\n", r.Method, r.RequestURI, r.ProtoMajor, r.ProtoMinor))

	for name, values := range r.Header {
		for _, value := range values {
			buffer.WriteString(fmt.Sprintf("%s: %s\n", name, value))
		}
	}

	if r.ContentLength != 0 {
		if r.Form == nil {
			r.ParseMultipartForm(32 << 20)
		}

		if r.Form != nil && len(r.Form) > 0 {
			buffer.WriteString("\nForm Data:\n")
			for name, values := range r.Form {
				for _, value := range values {
					buffer.WriteString(fmt.Sprintf("\n\t%s=%s", name, value))
				}
			}
		} else {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				Warn("[%s] Error reading request body: %s", stanza.Name, err)
				return "", fmt.Errorf("IO Error")
			}
			buffer.WriteString("\n")
			buffer.WriteString(string(body))
		}
	}
	return buffer.String(), nil
}

func waitForSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	if DEBUG {
		Debug("Received signal: %s", sig)
	}
}
