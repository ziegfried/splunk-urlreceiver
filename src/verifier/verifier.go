package verifier

import "net/http"
import . "modinputs"

const (
	VERIFY_OK = iota
	VERIFY_HANDLED
	VERIFY_ERROR
)

type VerifierResult int

type Verifier func(http.ResponseWriter, *http.Request) VerifierResult
type VerifierFunc func(InputStanza) Verifier

func GetVerifier(t string, stanza InputStanza, debug bool) Verifier {
	switch t {
	case "signature":
		return HmacRequestVerifier(stanza)
	case "meraki":
		return MerakiRequestVerifier(stanza)
	case "-":
		return nil
	default:
		Log.Error("[%s] Unknown verifier_type: \"%s\"", stanza.Name, t)
		return nil
	}
}
