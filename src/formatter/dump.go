package formatter

import (
	"bytes"
	"fmt"
	"io/ioutil"
	. "modinputs"
	"net/http"
)

func dumpRequest(r *http.Request, stanza InputStanza, debug bool) (string, error) {
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
				Log.Warn("[%s] Error reading request body: %s", stanza.Name, err)
				return "", fmt.Errorf("IO Error")
			}
			buffer.WriteString("\n")
			buffer.WriteString(string(body))
		}
	}
	return buffer.String(), nil
}
