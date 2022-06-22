package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Fonzeca/Api-Ruta3/model"
	"github.com/Fonzeca/Api-Ruta3/utils"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func Router(r *mux.Router) {
	//Obtenemos la configuracion
	configRouter := model.Config{}
	err := viper.Unmarshal(&configRouter)
	if err != nil {
		panic(err)
	}

	//Armamos las rutas
	for _, v := range configRouter.Services {
		r.PathPrefix(v.Prefix).HandlerFunc(BuildGeneralHandler(v))
	}

	r.HandleFunc("/auth/login", BuildAuthHandler(configRouter.Auth))
	r.Use(BuildAuthMiddleware(configRouter.Auth))
}

// Genera una ruta generica
func BuildGeneralHandler(service model.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//Le quitamos el prefix para que vaya al servicio
		pathWithoutPrefix := strings.Replace(r.URL.Path, service.Prefix, "", 1)

		//Mandamos la llamada al service
		res, err := utils.SendHtppRequest(service.ServiceUrl+pathWithoutPrefix, r, service.Headers)
		if err != nil {
			panic(err)
		}

		//Copiamos los headers de respuesta
		utils.CopyHeaders(res.Header, w.Header())

		//Seteamos el status code
		w.WriteHeader(res.StatusCode)

		//Copiamos el body de respuesta
		utils.CopyBody(res, w)
	}
}

func BuildAuthHandler(auth model.Auth) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//Mandamos la llamada al service
		res, err := utils.SendHtppRequest(auth.LoginUrl, r, map[string]string{"apiKey": auth.UserHubApiKey})
		if err != nil {
			panic(err)
		}

		//Copiamos los headers de respuesta
		utils.CopyHeaders(res.Header, w.Header())

		//Seteamos el status code
		w.WriteHeader(res.StatusCode)

		//Copiamos el body de respuesta
		utils.CopyBody(res, w)
	}
}

func BuildAuthMiddleware(auth model.Auth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !strings.HasPrefix(r.URL.Path, "/auth") {
				//Mandamos la llamada al service
				res, err := utils.SendHtppRequest(auth.ValidateTokenUrl, r, nil)
				if err != nil {
					panic(err)
				}
				fmt.Println(res)

			}
			next.ServeHTTP(w, r)
		})
	}
}
