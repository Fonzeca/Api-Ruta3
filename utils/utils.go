package utils

import (
	"net/http"
)

func CopyHeaders(origin http.Header, destino http.Header) {
	//Seteamos los headers
	for key, v := range origin {
		destino.Set(key, string(v[0]))
	}
}
