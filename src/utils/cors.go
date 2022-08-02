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
		r.Header.Set("Access-Control-Allow-Origin", origin)
		r.Header.Set("Access-Control-Allow-Credentials", "true")
		r.Header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		r.Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}

}
