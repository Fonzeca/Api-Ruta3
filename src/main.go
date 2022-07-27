package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/Fonzeca/Api-Ruta3/src/docs"
	"github.com/Fonzeca/Api-Ruta3/src/utils"
)

// @title        Api-Ruta3
// @version      1.0
// @description  This is a ApiGatway.
func main() {
	configViper()

	logger, err := utils.NewApiGatewayLogger()
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	Router(r)
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			utils.OpsProcessed.Inc()

			logger.Log(r.Clone(r.Context()))

			next.ServeHTTP(w, r)
		})
	})
	r.Handle("/metrics", promhttp.Handler())

	utils.EnableCORS(r)

	log.Fatal(http.ListenAndServe(":8082", r))
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
