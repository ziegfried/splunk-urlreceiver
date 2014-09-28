package verifier

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	. "modinputs"
	"net/http"
	. "scheme"
	. "utils"
)

func MerakiRequestVerifier(stanza InputStanza) Verifier {
	verifier := stanza.GetParamValue(PARAM_MERAKI_VERIFIER, "")
	secret := stanza.GetParamValue(PARAM_MERAKI_SECRET, "n/a")
	debug := NormalizeBoolean(stanza.GetParamValue(PARAM_DEBUG, "false"), false)
	return func(w http.ResponseWriter, r *http.Request) VerifierResult {
		if r.Method == "GET" {
			if debug {
				Log.Info("[%s] Responding with meraki verifier %s to GET request", stanza.Name, verifier)
			}
			fmt.Fprint(w, verifier)
			return VERIFY_HANDLED
		} else if r.Method == "POST" {
			if checkMerakiSecret(r, secret, stanza, debug) {
				return VERIFY_OK
			} else {
				return VERIFY_ERROR
			}
		} else {
			Log.Warn("[%s] Unexpected method=%s", stanza.Name, r.Method)
			return VERIFY_ERROR
		}
	}
}

func checkMerakiSecret(r *http.Request, secret string, stanza InputStanza, debug bool) bool {
	bodyWrapper, err := WrapRequestBody(r)
	if err != nil {
		Log.Error("[%s] Error reading request body: %s", stanza.Name, err)
		return false
	}
	body, _ := ioutil.ReadAll(r.Body)
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		Log.Error("[%s] Error parsing JSON data: %s", stanza.Name, err)
		bodyWrapper.Reset()
		return false
	} else {
		bodySecret, found := data["secret"]
		if found && secret == bodySecret {
			if debug {
				Log.Info("[%s] Secret successfully verified", stanza.Name)
			}
			delete(data, "secret")
			newBody, _ := json.Marshal(data)
			NewRequestBody(r, newBody)
			return true
		} else {
			Log.Warn("[%s] Provided meraki secret \"%s\" does not match configured one \"%s\"", stanza.Name, secret, bodySecret)
			bodyWrapper.Reset()
			return false
		}
	}
}
