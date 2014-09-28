package formatter

import (
	"bytes"
	. "modinputs"
	"net/http"
	. "scheme"
	"strings"
)

func readFormField(r *http.Request, stanza InputStanza, debug bool) (string, error) {
	field := stanza.GetParamValue(PARAM_FORM_FIELD, "data")
	if debug {
		Log.Debug("[%s] Reading form field=%s", stanza.Name, field)
	}
	return r.FormValue(field), nil
}

func formFieldsKV(r *http.Request, stanza InputStanza, debug bool) (string, error) {
	if debug {
		Log.Debug("[%s] Reading all form fields", stanza.Name)
	}

	if r.Form == nil {
		r.ParseMultipartForm(32 << 20)
	}

	var buffer bytes.Buffer
	for param, values := range r.Form {
		for _, value := range values {
			buffer.WriteString(param)
			buffer.WriteString("=\"")
			buffer.WriteString(strings.Replace(value, "\"", "\\\"", -1))
			buffer.WriteString("\" ")
		}
	}
	return buffer.String(), nil
}
