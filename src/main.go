package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"

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
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Log(r)

			next.ServeHTTP(w, r)
		})
	})
	Router(r)

	fmt.Println("░█████╗░██████╗░██╗░░░░░░██████╗░██╗░░░██╗████████╗░█████╗░██████╗░\n██╔══██╗██╔══██╗██║░░░░░░██╔══██╗██║░░░██║╚══██╔══╝██╔══██╗╚════██╗\n███████║██████╔╝██║█████╗██████╔╝██║░░░██║░░░██║░░░███████║░█████╔╝\n██╔══██║██╔═══╝░██║╚════╝██╔══██╗██║░░░██║░░░██║░░░██╔══██║░╚═══██╗\n██║░░██║██║░░░░░██║░░░░░░██║░░██║╚██████╔╝░░░██║░░░██║░░██║██████╔╝\n╚═╝░░╚═╝╚═╝░░░░░╚═╝░░░░░░╚═╝░░╚═╝░╚═════╝░░░░╚═╝░░░╚═╝░░╚═╝╚═════╝░")
	log.Fatal(http.ListenAndServe(":8082", r))
}

func configViper() {
	viper.SetConfigName("config.json")
	viper.SetConfigType("json")
	viper.AddConfigPath("../")
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
