package main

import (
	"io"
	"net/http"
	"net/http/httputil"
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
		proxy := httputil.ReverseProxy{
			Director: buildDirectorFunc(v),
			ModifyResponse: func(r *http.Response) error {
				utils.EnableCORS(r)
				return nil
			},
		}
		r.PathPrefix(v.Prefix).Handler(&proxy)
	}

	r.Path("/auth/login").Handler(buildAuthHandler(configRouter.Auth)).Methods(http.MethodPost)
	r.Use(BuildAuthMiddleware(configRouter.Auth, configRouter.Services))
}

func buildDirectorFunc(service model.Service) func(r *http.Request) {
	return func(r *http.Request) {
		r.URL.Scheme = "http"
		r.URL.Host = service.ServiceUrl
		r.URL.Path = strings.Replace(r.URL.Path, service.Prefix, "", 1)
	}
}

func buildAuthHandler(auth model.Auth) *httputil.ReverseProxy {
	return &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL.Scheme = "http"
			r.URL.Host = auth.LoginUrl
			r.Header.Set("apiKey", auth.UserHubApiKey)
		},
		ModifyResponse: func(r *http.Response) error {
			utils.EnableCORS(r)
			return nil
		},
	}
}

//Genera el middleware para validar el token
func BuildAuthMiddleware(auth model.Auth, services []model.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// if r.Method == "OPTIONS" {
			// 	next.ServeHTTP(w, r)
			// 	return
			// }

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
			rq, err := http.NewRequest("POST", "http://"+auth.ValidateTokenUrl, nil)
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
