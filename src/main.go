package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"

	"github.com/Fonzeca/Api-Ruta3/src/model"
	"github.com/Fonzeca/Api-Ruta3/src/utils"
)

// @title        Api-Ruta3
// @version      1.0
// @description  This is a ApiGatway.
func main() {
	configViper()

	//Obtenemos la configuracion
	configRouter := model.Config{}
	err := viper.Unmarshal(&configRouter)
	if err != nil {
		panic(err)
	}

	//Logger
	logger, err := utils.NewApiGatewayLogger()
	if err != nil {
		panic(err)
	}

	//Router
	r := mux.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Log(r)

			next.ServeHTTP(w, r)
		})
	})
	Router(r, configRouter)

	//Middleware for middlewares
	m := http.NewServeMux()
	m.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		path := req.URL.Path

		for _, s := range configRouter.Services {
			match, err := regexp.MatchString("^"+s.Prefix, path)
			if err != nil {
				panic(err)
			}
			if match {
				r.ServeHTTP(w, req)
				return
			}
		}

		match, _ := regexp.MatchString("^/auth", path)
		if match {
			r.ServeHTTP(w, req)
			return
		}

		req.URL.Path = "/carmind" + path

		r.ServeHTTP(w, req)
	})

	fmt.Println("░█████╗░██████╗░██╗░░░░░░██████╗░██╗░░░██╗████████╗░█████╗░██████╗░\n██╔══██╗██╔══██╗██║░░░░░░██╔══██╗██║░░░██║╚══██╔══╝██╔══██╗╚════██╗\n███████║██████╔╝██║█████╗██████╔╝██║░░░██║░░░██║░░░███████║░█████╔╝\n██╔══██║██╔═══╝░██║╚════╝██╔══██╗██║░░░██║░░░██║░░░██╔══██║░╚═══██╗\n██║░░██║██║░░░░░██║░░░░░░██║░░██║╚██████╔╝░░░██║░░░██║░░██║██████╔╝\n╚═╝░░╚═╝╚═╝░░░░░╚═╝░░░░░░╚═╝░░╚═╝░╚═════╝░░░░╚═╝░░░╚═╝░░╚═╝╚═════╝░")
	log.Fatal(http.ListenAndServe(":8082", m))
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
