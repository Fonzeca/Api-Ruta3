package utils

import (
	"net/http"
)

var (
	allowList = map[string]bool{
		"http://vps-1791261-x.dattaweb.com": true,
		"http://carmind-app.com":            true,
		"http://localhost:4200":             true,
	}
)

func EnableCORS(r *http.Response) {
	if origin := r.Request.Header.Get("Origin"); allowList[origin] {
		if r.Request.Method == http.MethodOptions {
			r.Header.Set("Access-Control-Allow-Origin", origin)
			r.Header.Set("Access-Control-Allow-Credentials", "true")
			r.Header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			r.Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		} else {
			r.Header.Set("Access-Control-Allow-Origin", origin)
		}
	}
}

func WriteEnableCORS(w http.ResponseWriter, r *http.Request) {
	if origin := r.Header.Get("Origin"); allowList[origin] {
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		} else {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
	}
}
