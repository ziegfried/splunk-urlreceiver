package verifier

import "testing"
import "net/http"
import "utils"
import "modinputs"
import "io/ioutil"

func TestVerifyMerakiSecret(t *testing.T) {
	r := new(http.Request)
	var stanza modinputs.InputStanza = modinputs.InputStanza{
		Name: "Test",
	}
	utils.NewRequestBody(r, []byte(`{"secret":"asdf","foo":"bar"}`))
	if !checkMerakiSecret(r, "asdf", stanza, false) {
		t.Error()
	}
}

func TestVerifyInvalidSecret(t *testing.T) {
	r := new(http.Request)
	var stanza modinputs.InputStanza = modinputs.InputStanza{
		Name: "Test",
	}
	utils.NewRequestBody(r, []byte(`{"secret":"asdf","foo":"bar"}`))
	if checkMerakiSecret(r, "invalid secret", stanza, false) {
		t.Error()
	}
}

func TestVerifyInvalidJSON(t *testing.T) {
	r := new(http.Request)
	var stanza modinputs.InputStanza = modinputs.InputStanza{
		Name: "Test",
	}
	utils.NewRequestBody(r, []byte(`\\//`))
	if checkMerakiSecret(r, "...", stanza, false) {
		t.Error()
	}
	if b, _ := ioutil.ReadAll(r.Body); b != nil {
		if string(b) != `\\//` {
			t.Error()
		}
	} else {
		t.Error()
	}
}
