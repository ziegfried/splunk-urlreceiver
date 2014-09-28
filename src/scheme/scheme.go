package scheme

import . "modinputs"

const (
	PARAM_PORT = "port"
	PARAM_PATH = "path"

	PARAM_SSL             = "ssl"
	PARAM_SSL_CERT_PATH   = "ssl_cert_path"
	PARAM_SSL_KEY_PATH    = "ssl_key_path"
	PARAM_SSL_CERT        = "ssl_cert"
	PARAM_SSL_KEY         = "ssl_key"
	PARAM_SSL_SELF_SIGNED = "ssl_self_signed"
	PARAM_SSL_HOST        = "ssl_host"
	PARAM_SSL_KEY_BITS    = "ssl_key_bits"

	PARAM_DATA_FORMAT = "data_format"
	PARAM_FORM_FIELD  = "form_field"

	PARAM_HOST_FROM_CLIENTIP = "host_from_clientip"

	PARAM_RESPONSE = "response"

	PARAM_DEBUG = "debug"

	PARAM_VERIFY_TYPE      = "verify_type"
	PARAM_SIGNATURE_SECRET = "signature_secret"
	PARAM_SIGNATURE_HEADER = "signature_header"
	PARAM_MERAKI_SECRET    = "meraki_secret"
	PARAM_MERAKI_VERIFIER  = "meraki_verifier"
)

var SCHEME = Scheme{
	Title:                 "URL-Receiver",
	Description:           "Receive and index data via URLs (such as Webhooks)",
	UseExternalValidation: false,
	StreamingMode:         "xml",
	UseSingleInstance:     true,
	EndpointArguments: []EndpointArg{
		EndpointArg{
			Name:             PARAM_PORT,
			Validation:       "is_port('port')",
			Description:      "TCP port to listen for HTTP requests",
			RequiredOnCreate: true,
		},
		EndpointArg{
			Name:             PARAM_PATH,
			Description:      "URL path to match against when receiving data",
			RequiredOnCreate: true,
		},
		EndpointArg{
			Name:        PARAM_SSL,
			Description: "Enable SSL for the webserver port of this input",
			DataType:    "boolean",
			Default:     "false",
		},
		EndpointArg{
			Name:        PARAM_SSL_CERT_PATH,
			Title:       "Path to SSL certificate file",
			Description: "",
		},
		EndpointArg{
			Name:        PARAM_SSL_KEY_PATH,
			Title:       "Path to SSL private key file",
			Description: "",
		},
		EndpointArg{
			Name:        PARAM_SSL_CERT,
			Title:       "SSL certificate",
			Description: "Content of the SSL certificate in Base64 encoded DER format",
		},
		EndpointArg{
			Name:        PARAM_SSL_KEY,
			Title:       "SSL private key",
			Description: "Content of the SSL private key in Base64 encoded DER format",
		},
		EndpointArg{
			Name:  PARAM_SSL_SELF_SIGNED,
			Title: "Generate self-signed certificate",
			Description: "If enabled, a new self-signed certificate is generated each time the webserver is started.\n" +
				"(only recommend for testing purposes, the webserver is restarted each time any of the input configurations are changed, added or deleted)",
			DataType: "boolean",
			Default:  "false",
		},
		EndpointArg{
			Name:        PARAM_SSL_HOST,
			Title:       "Host value for self-signed certificate",
			Description: "",
			Default:     "<value of splunk mgmt host>",
		},
		EndpointArg{
			Name:        PARAM_SSL_KEY_BITS,
			Title:       "Key size for self-signed certificate",
			Description: "",
			DataType:    "number",
			Default:     "2048",
		},
		EndpointArg{
			Name:  PARAM_DATA_FORMAT,
			Title: "Data Format",
			Description: "Define how data is extracted from the incoming HTTP request and formatted as text for the event to be indexed by Splunk\n" +
				"Possible values are:\n" +
				"   - raw_body: Read the raw (unparsed) body of the HTTP request\n" +
				"   - form_field: Read a particular form field (either POST or GET param, also supports multipart)\n" +
				"   - form_kv: Build key-value pairs from form fields of the HTTP request (GET/POST params)\n" +
				"   - full_request: Dump full request details (Request line, headers and decoded parameters)",
			EnumValues:       []string{"raw_body", "form_field", "form_kv", "full_request"},
			Default:          `"raw_body"`,
			RequiredOnCreate: true,
		},
		EndpointArg{
			Name:        PARAM_FORM_FIELD,
			Description: "Specify the form field to index (only if data retrieval is form_field)",
		},
		EndpointArg{
			Name:        PARAM_HOST_FROM_CLIENTIP,
			Title:       "Use client IP address for host field value",
			Description: "If enabled the input will pass the IP address of the client sending the HTTP request as the host field value for each event",
			DataType:    "boolean",
			Default:     "false",
		},
		EndpointArg{
			Name:        PARAM_RESPONSE,
			Title:       "Response text",
			Description: "Text the webserver responds with after successfully receiving data",
			Default:     `"OK"`,
		},
		EndpointArg{
			Name:        PARAM_DEBUG,
			Title:       "Debug logging",
			Description: "Enable debug logging for this input (logged to splunkd.log)",
			DataType:    "boolean",
			Default:     "false",
		},
		EndpointArg{
			Name:       PARAM_VERIFY_TYPE,
			Title:      "Verify incoming requests",
			EnumValues: []string{"signature", "meraki"},
			Description: "Select the type of verification for incoming requests. By default requests are not verified.\n" +
				"Possible values are:\n" +
				"   - signature: Verify HMAC of the request body to header value (eg. for Github webhooks)\n" +
				"   - meraki: Verify requests from Cisco Meraki CMX API requests\n" +
				"By default requests are not verified.",
			Default: "-",
		},
		EndpointArg{
			Name:  PARAM_SIGNATURE_SECRET,
			Title: "Github Secret Key",
			Description: "The secret key configured for the webhook. The secret is used to compute the HMAC of the request" +
				" body and compare it to the supplied header value to verify the origin of the request.\n" +
				"(only applies if verify_type = signature)",
		},
		EndpointArg{
			Name:        PARAM_SIGNATURE_HEADER,
			Title:       "Signature Header",
			Description: "The name of the HTTP header in which the signature value is passed",
			Default:     `"X-Hub-Signature"`,
		},
		EndpointArg{
			Name:        PARAM_MERAKI_VERIFIER,
			Title:       "Meraki Verifier",
			Description: "(only applies if verify_type = meraki)",
		},
		EndpointArg{
			Name:        PARAM_MERAKI_SECRET,
			Title:       "Meraki Secret",
			Description: "(only applies if verify_type = meraki)",
		},
	},
}
