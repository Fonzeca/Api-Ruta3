package utils

import (
	"io/ioutil"
	"net/http"
)

func SendHtppRequest(url string, r *http.Request, addHeaders map[string]string) (*http.Response, error) {

	//Armamos el request
	rq, err := http.NewRequest(r.Method, url, r.Body)
	if err != nil {
		return nil, err
	}

	//Se copian los headers
	CopyHeaders(r.Header, rq.Header)

	//Se agregan los headers que se configurarion
	for k, v := range addHeaders {
		rq.Header.Set(k, v)
	}

	//Hacemos la llamada
	res, err := http.DefaultClient.Do(rq)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func CopyHeaders(origin http.Header, destino http.Header) {
	//Seteamos los headers
	for key, v := range origin {
		destino.Set(key, string(v[0]))
	}
}

func CopyBody(res *http.Response, w http.ResponseWriter) {
	//Leemos el body
	by, err := ioutil.ReadAll(res.Body)
	if err != nil {
		//TODO: log error
		return
	}
	w.Write(by)
}
