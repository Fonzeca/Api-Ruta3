package main

import (
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/Fonzeca/Api-Ruta3/src/model"
	"github.com/Fonzeca/Api-Ruta3/src/utils"
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
	r.Use(BuildAuthMiddleware(configRouter.Auth, configRouter.Services))
}

// Genera una ruta generica
func BuildGeneralHandler(service model.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		rqService, err := utils.DeepCopyRequest(r)

		//Le quitamos el prefix para que vaya al servicio
		rqService.URL, _ = url.ParseRequestURI(service.ServiceUrl + strings.Replace(rqService.RequestURI, service.Prefix, "", 1))
		rqService.RequestURI = ""

		//Mandamos la llamada al service
		res, err := http.DefaultClient.Do(rqService)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()

		//Copiamos los headers de respuesta
		utils.CopyHeaders(res.Header, w.Header())

		//Copiamos el status code
		w.WriteHeader(res.StatusCode)

		//Copiamos el body de respuesta
		io.Copy(w, res.Body)
	}
}

//Genera la ruta de login
func BuildAuthHandler(auth model.Auth) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		rqService := r.Clone(r.Context())

		//Le quitamos el prefix para que vaya al servicio
		rqService.URL, _ = url.ParseRequestURI(auth.LoginUrl)
		rqService.RequestURI = "" //Da error si esto no esta vacio.
		rqService.Header.Set("apiKey", auth.UserHubApiKey)

		//Mandamos la llamada al service
		res, err := http.DefaultClient.Do(rqService)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()

		//Copiamos los headers de respuesta
		utils.CopyHeaders(res.Header, w.Header())

		//Copiamos el status code
		w.WriteHeader(res.StatusCode)

		//Copiamos el body de respuesta
		io.Copy(w, res.Body)

	}
}

//Genera el middleware para validar el token
func BuildAuthMiddleware(auth model.Auth, services []model.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "OPTIONS" {
				next.ServeHTTP(w, r)
				return
			}

			if strings.HasPrefix(r.URL.Path, "/auth") {
				next.ServeHTTP(w, r)
				return
			}

			for _, se := range services {
				if strings.HasPrefix(r.URL.Path, se.Prefix) {
					pathToVerificate := strings.TrimPrefix(r.URL.Path, se.Prefix)

					for _, publicUrl := range se.PublicUrls {
						matched, _ := regexp.MatchString("^"+publicUrl, pathToVerificate)
						if matched {
							next.ServeHTTP(w, r)
							return
						}
					}

				}
			}

			//Armamos el request para verificar el token
			rq, err := http.NewRequest("POST", auth.ValidateTokenUrl, nil)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				panic(err)
			}
			rq.Header.Set("apiKey", auth.UserHubApiKey)
			rq.Header.Set("Authorization", r.Header.Get("Authorization"))

			//Hacemos la llamada
			res, err := http.DefaultClient.Do(rq)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				panic(err)
			}

			if res.StatusCode == 200 {
				next.ServeHTTP(w, r)
				return
			}

			//Copiamos los headers de respuesta
			utils.CopyHeaders(res.Header, w.Header())

			//Si no es 200, devolvemos que no tiene acceso
			w.WriteHeader(http.StatusUnauthorized)

			//Copiamos el body de respuesta
			io.Copy(w, res.Body)
			return
		})
	}
}
