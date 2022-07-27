package utils

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/afiskon/promtail-client/promtail"
)

type apigatwayLogger struct {
	loki promtail.Client
}

func NewApiGatewayLogger() (apigatwayLogger, error) {
	labels := "{source=\"api-ruta3\",job=\"job_api-ruta3\"}"

	conf := promtail.ClientConfig{
		PushURL:            "http://vps-2721477-x.dattaweb.com:3100/loki/api/v1/push",
		Labels:             labels,
		BatchWait:          5 * time.Second,
		BatchEntriesNumber: 10000,
		SendLevel:          promtail.INFO,
		PrintLevel:         promtail.ERROR,
	}

	loki, err := promtail.NewClientProto(conf)
	if err != nil {
		return apigatwayLogger{}, err
	}

	logger := &apigatwayLogger{
		loki: loki,
	}

	return *logger, nil
}

func (logger *apigatwayLogger) Log(r *http.Request) {
	go logger.log(r)
}

func (logger *apigatwayLogger) log(r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	logger.loki.Infof("Method: %s, Path: %s, Query: %s, Body: %s", r.Method, r.URL.Path, r.URL.RawQuery, string(bodyBytes))
}
