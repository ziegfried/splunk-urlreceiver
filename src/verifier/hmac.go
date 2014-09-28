package verifier

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
	"io/ioutil"
	. "modinputs"
	"net/http"
	. "scheme"
	"strings"
	. "utils"
)

const DEFAULT_SIGNATURE_HEADER = "X-Hub-Signature"

func HmacRequestVerifier(stanza InputStanza) Verifier {
	secret := stanza.GetParamValue(PARAM_SIGNATURE_SECRET, "")
	header := stanza.GetParamValue(PARAM_SIGNATURE_HEADER, DEFAULT_SIGNATURE_HEADER)
	debug := NormalizeBoolean(stanza.GetParamValue(PARAM_DEBUG, "false"), false)
	return func(w http.ResponseWriter, r *http.Request) VerifierResult {
		if verifyRequestHmac(r, secret, header, debug) {
			return VERIFY_OK
		} else {
			return VERIFY_ERROR
		}
	}
}

func verifyRequestHmac(r *http.Request, secret, header string, debug bool) bool {
	bodyWrapper, err := WrapRequestBody(r)
	if err != nil {
		Log.Error("Error reading request body: %s", err)
		return false
	}
	body, _ := ioutil.ReadAll(bodyWrapper)
	bodyWrapper.Reset()

	hash, expected := extractHash(r.Header[header][0])
	if hash == nil {
		return false
	}

	return verifyHmac(body, secret, expected, hash)
}

func extractHash(expected string) (func() hash.Hash, string) {
	if strings.HasPrefix(expected, "sha1=") {
		return sha1.New, expected[5:]
	} else if strings.HasPrefix(expected, "sha256=") {
		return sha256.New, expected[7:]
	} else if strings.HasPrefix(expected, "sha512=") {
		return sha512.New, expected[7:]
	} else if strings.HasPrefix(expected, "md5=") {
		return md5.New, expected[4:]
	} else {
		return nil, ""
	}
}

func verifyHmac(body []byte, secret, expected string, hash func() hash.Hash) bool {
	mac := hmac.New(hash, []byte(secret))
	mac.Write(body)
	expectedBytes, _ := hex.DecodeString(expected)
	actualMac := mac.Sum(nil)
	return hmac.Equal(actualMac, expectedBytes)
}
