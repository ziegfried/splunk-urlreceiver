[urlreceiver://<name>]

port = <value>
* TCP port to listen for HTTP requests

path = <value>
* URL path to match against when receiving data

data_retrieval = <value>
* Define how data is extracted from the incoming HTTP request
* Possible values are:
*   - raw_body: Read the raw (unparsed) body of the HTTP request
*   - form_field: Read a particular form field (either POST or GET param, also supports multipart)
*   - form_kv: Build key-value pairs from form fields of the HTTP request (GET/POST params)
*   - full_request: Dump full request details (Request line, headers and decoded parameters)

form_field = <value>
* "Specify the form field to index (only if data retrieval is form_field)

host_from_clientip = <value>
* Use client IP address for host field value. If enabled the input will pass the IP address of the client
* sending the HTTP request as the host field value for each event

response = <value>
* Text the webserver responds with after successfully receiving data (Default is "OK")

debug = <value>
* Enable debug logging for this input (logged to splunkd.log)
