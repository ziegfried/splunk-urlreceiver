[urlreceiver://<name>]

port = <string>
* (required)
* TCP port to listen for HTTP requests

path = <string>
* (required)
* URL path to match against when receiving data

ssl = [true|false]
* (optional, default: false)
* Enable SSL for the webserver port of this input

ssl_cert_path = <string>
* Path to SSL certificate file (optional)

ssl_key_path = <string>
* Path to SSL private key file (optional)

ssl_cert = <string>
* SSL certificate (optional)
* Content of the SSL certificate in Base64 encoded DER format

ssl_key = <string>
* SSL private key (optional)
* Content of the SSL private key in Base64 encoded DER format

ssl_self_signed = [true|false]
* Generate self-signed certificate (optional, default: false)
* If enabled, a new self-signed certificate is generated each time the webserver is started.
* (only recommend for testing purposes, the webserver is restarted each time any of the input configurations are changed, added or deleted)

ssl_host = <string>
* Host value for self-signed certificate (optional, default: <value of splunk mgmt host>)

ssl_key_bits = <string>
* Key size for self-signed certificate (optional, default: 2048)

data_format = [raw_body|form_field|form_kv|full_request]
* Data Format (required, default: "raw_body")
* Define how data is extracted from the incoming HTTP request and formatted as text for the event to be indexed by Splunk
* Possible values are:
*    - raw_body: Read the raw (unparsed) body of the HTTP request
*    - form_field: Read a particular form field (either POST or GET param, also supports multipart)
*    - form_kv: Build key-value pairs from form fields of the HTTP request (GET/POST params)
*    - full_request: Dump full request details (Request line, headers and decoded parameters)

form_field = <string>
* (optional)
* Specify the form field to index (only if data retrieval is form_field)

host_from_clientip = [true|false]
* Use client IP address for host field value (optional, default: false)
* If enabled the input will pass the IP address of the client sending the HTTP request as the host field value for each event

response = <string>
* Response text (optional, default: "OK")
* Text the webserver responds with after successfully receiving data

debug = [true|false]
* Debug logging (optional, default: false)
* Enable debug logging for this input (logged to splunkd.log)

verify_type = [signature|meraki]
* Verify incoming requests (optional, default: -)
* Select the type of verification for incoming requests. By default requests are not verified.
* Possible values are:
*    - signature: Verify HMAC of the request body to header value (eg. for Github webhooks)
*    - meraki: Verify requests from Cisco Meraki CMX API requests
* By default requests are not verified.

signature_secret = <string>
* Github Secret Key (optional)
* The secret key configured for the webhook. The secret is used to compute the HMAC of the request body and compare it to the supplied header value to verify the origin of the request.
* (only applies if verify_type = signature)

signature_header = <string>
* Signature Header (optional, default: "X-Hub-Signature")
* The name of the HTTP header in which the signature value is passed

meraki_verifier = <string>
* Meraki Verifier (optional)
* (only applies if verify_type = meraki)

meraki_secret = <string>
* Meraki Secret (optional)
* (only applies if verify_type = meraki)

