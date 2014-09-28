package formatter

import (
	"fmt"
	"io/ioutil"
	. "modinputs"
	"net/http"
)

func readRawBody(r *http.Request, stanza InputStanza, debug bool) (string, error) {
	if debug {
		Log.Debug("[%s] Reading request body", stanza.Name)
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Log.Warn("[%s] Error reading request body: %s", stanza.Name, err)
		return "", fmt.Errorf("IO Error")
	}
	if len(body) > 0 {
		return string(body), nil
	} else {
		Log.Warn("[%s] Received empty request body", stanza.Name)
		return "", fmt.Errorf("Empty request body")
	}
}
