# URL Receiver Modular Input

Simple modular input written in Go which starts a webserver and listens for incomming requests. Data received in the request body or form params is the indexed by Splunk. This can be used to receive data from Webhooks.

## Usage

Install the app in Splunk Manager and go to Settings -> Data Inputs -> URL Receiver to add a new input.

### Input Settings

- `port` (int, **required**)
TCP Port to list for incoming HTTP requests.
- `path` (string, **required**)
URL path to match against to associate an incoming request to a configured input. Requests not matching the path of any input will receive a HTTP 404 response and will not be indexed.
- `read_body` (boolean, default `true`)
Read and index the raw request body of incoming HTTP requets.
- `form_field` (string, default `""`)
Read and index the value of a form field of incoming HTTP requests. The form field can either be the value of a param of the decoded POST/PUT body (either url-encoded or multipart) or the value of the query string for a GET request. Note: This only works if `read_body` is set to `false`.
- `host_from_clientip` (boolean, default `false`)
If enabled use the IP address of the client sending the HTTP request as the `host` field of the indexed event.

## Build from source

1. Install the [Go](http://golang.org/) toolchain and set up the Go compiler for all target platforms/architectures
    ```
    cd $GOROOT/src
    GOOS=darwin GOARCH=amd64 ./make.bash --noclean
    GOOS=darwin GOARCH=386 ./make.bash --noclean
    GOOS=linux GOARCH=amd64 ./make.bash --noclean
    GOOS=linux GOARCH=386 ./make.bash --noclean
    GOOS=windows GOARCH=amd64 ./make.bash --noclean
    GOOS=windows GOARCH=386 ./make.bash --noclean
    ```
2. Enable all desired platforms/architectures in build.sh
    For example, if you only want to build the binaries for linux, comment the build_binary calls for OSX (darwin) and Windows:
    ```
    # build_binary darwin amd64 darwin_x86_64 urlreceiver
    # build_binary darwin 386 darwin_x86 urlreceiver
    build_binary linux amd64 linux_x86_64 urlreceiver
    build_binary linux 386 linux_x86 urlreceiver
    # build_binary windows amd64 windows_x86_64 urlreceiver.exe
    # build_binary windows 386 windows_x86 urlreceiver.exe
    ```
3. Run build.sh
    ```
    sh build.sh
    ```
    The package will be created in the `dist` folder in the project directory.

