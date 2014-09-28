package utils

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

type BodyWrapper struct {
	r *bytes.Reader
}

func (w BodyWrapper) Close() error {
	w.Reset()
	return nil
}

func (w BodyWrapper) Read(b []byte) (int, error) {
	return w.r.Read(b)
}

func (w BodyWrapper) Reset() {
	// Reset body reader of the request so it can be read again
	w.r.Seek(0, 0)
}

func WrapRequestBody(r *http.Request) (*BodyWrapper, error) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return NewRequestBody(r, b), nil
}

func NewRequestBody(r *http.Request, b []byte) *BodyWrapper {
	wrapper := new(BodyWrapper)
	wrapper.r = bytes.NewReader(b)
	r.Body = *wrapper
	return wrapper
}
