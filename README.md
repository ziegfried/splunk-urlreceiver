# URL Receiver Modular Input

Simple modular input written in Go which starts a webserver and listens for incomming requests. Data received in the request body or form params is the indexed by Splunk. This can be used to receive data from Webhooks.

## Dependencies

- Splunk 5+ (tested with 6.1)
- Linux, MacOSX or Windows (Windows is untested but should work)

## Usage

Install the app in Splunk Manager and go to Settings -> Data Inputs -> URL Receiver to add a new input.

### Input Settings

- `port` (int, **required**)

    TCP Port to list for incoming HTTP requests.

- `path` (string, **required**)

    URL path to match against to associate an incoming request to a configured input. Requests not matching the path of any input will receive a HTTP 404 response and will not be indexed.

- `data_retrieval` (string, default `"raw_body"`)

    "Define how data is extracted from incoming HTTP requests. Possible values are:
    - `"raw_body"` - Index the raw body content of the HTTP request
    - `"form_field"` - Index the content of a particular form field (see `form_field`)
    - `"form_kv"` - Index generated key-value pairs for all POST- or GET params
    - `"full_request"` - Dump full request details (including request line, headers and form data)

- `form_field` (string, default `""`)

    (For `data_retrieval=form_field`) Read and index the value of a form field of incoming HTTP requests. The form field can either be the value of a param of the decoded POST/PUT body (either url-encoded or multipart) or the value of the query string for a GET request.

- `host_from_clientip` (boolean, default `false`)

    If enabled use the IP address of the client sending the HTTP request as the `host` field of the indexed event.

- `debug` (boolean, default `false`)

    Enable debug logging for this input (logs to splunkd.log)

- `response` (string, default "OK")

    Text send back to the client after successfully receiving data.

## Build from source

##### 1. Install the [Go](http://golang.org/) toolchain and set up the Go compiler for all target platforms/architectures

```sh
cd $GOROOT/src
GOOS=darwin GOARCH=amd64 ./make.bash --noclean
GOOS=darwin GOARCH=386 ./make.bash --noclean
GOOS=linux GOARCH=amd64 ./make.bash --noclean
GOOS=linux GOARCH=386 ./make.bash --noclean
GOOS=windows GOARCH=amd64 ./make.bash --noclean
GOOS=windows GOARCH=386 ./make.bash --noclean
```

##### 2. Enable all desired platforms/architectures in build.sh

If you only want to, for example, build the binaries for linux, comment the build_binary calls for OSX (darwin) and Windows:

```sh
# build_binary darwin amd64 darwin_x86_64 urlreceiver
# build_binary darwin 386 darwin_x86 urlreceiver
build_binary linux amd64 linux_x86_64 urlreceiver
build_binary linux 386 linux_x86 urlreceiver
# build_binary windows amd64 windows_x86_64 urlreceiver.exe
# build_binary windows 386 windows_x86 urlreceiver.exe
```

##### 3. Run build.sh

```sh
sh build.sh
```

The package will be created in the `dist` folder in the project directory.

## License

This add-on is licensed unter the terms of [Creative Commons 3.0 (CC BY 3.0)](http://creativecommons.org/licenses/by/3.0/)
