package utils

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func CopyHeaders(origin http.Header, destino http.Header) {
	//Seteamos los headers
	for key, v := range origin {
		destino.Set(key, string(v[0]))
	}
}

func DeepCopyRequest(r *http.Request) (*http.Request, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	r2 := r.Clone(r.Context())
	// clone body
	r.Body = ioutil.NopCloser(bytes.NewReader(body))
	r2.Body = ioutil.NopCloser(bytes.NewReader(body))

	// parse r1, proxy r2
	r.ParseForm()
	return r2, nil
}
