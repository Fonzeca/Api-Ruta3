package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/afiskon/promtail-client/promtail"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	"github.com/spf13/viper"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/Fonzeca/Api-Ruta3/docs"
	"github.com/Fonzeca/Api-Ruta3/utils"
)

// @title        Api-Ruta3
// @version      1.0
// @description  This is a ApiGatway.
func main() {
	configViper()

	labels := "{source=\"api-ruta3\",job=\"job_api-ruta3\"}"
	conf := promtail.ClientConfig{
		PushURL:            "http://vps-2721477-x.dattaweb.com:3100/loki/api/v1/push",
		Labels:             labels,
		BatchWait:          5 * time.Second,
		BatchEntriesNumber: 10000,
		SendLevel:          promtail.INFO,
		PrintLevel:         promtail.ERROR,
	}

	loki, err := promtail.NewClientJson(conf)
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	Router(r)

	mux.CORSMethodMiddleware(r)
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			utils.OpsProcessed.Inc()
			loki.Infof(r.URL.String())
			next.ServeHTTP(w, r)
		})
	})
	r.Handle("/metrics", promhttp.Handler())

	handler := cors.AllowAll().Handler(r)

	log.Fatal(http.ListenAndServe(":8082", handler))
}

func configViper() {
	viper.SetConfigName("config.json")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

// Handler
// @description Api test
// @id test
// @router /test [GET]
func ArticlesCategoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Category: %v\n", vars["category"])
}
