package formatter

import (
	. "modinputs"
	"net/http"
	. "scheme"
)

const (
	DATA_RAW_BODY     = "raw_body"
	DATA_FORM_FIELD   = "form_field"
	DATA_FORM_KV      = "form_kv"
	DATA_FULL_REQUEST = "full_request"
)

type Formatter func(*http.Request, InputStanza, bool) (string, error)

func GetFormatter(t string, stanza InputStanza, debug bool) Formatter {
	switch t {
	case DATA_RAW_BODY:
		if debug {
			Log.Info("[%s] Will index raw request body", stanza.Name)
		}
		return readRawBody
	case DATA_FORM_FIELD:
		if debug {
			Log.Info("[%s] Will index form field %s", stanza.Name, stanza.GetParamValue(PARAM_FORM_FIELD, "data (default)"))
		}
		return readFormField
	case DATA_FORM_KV:
		if debug {
			Log.Info("[%s] Will read form key-values", stanza.Name)
		}
		return formFieldsKV
	case DATA_FULL_REQUEST:
		if debug {
			Log.Info("[%s] Will dump full request details", stanza.Name)
		}
		return dumpRequest
	default:
		Log.Warn("[%s] Ignoring unknown data_retrieval type \"%s\" - indexing request body", stanza.Name, t)
		return readRawBody
	}
}
