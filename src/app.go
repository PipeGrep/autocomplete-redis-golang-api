package main

import (
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	"net/http"
	"strconv"
)

func main() {
	api := rest.NewApi()

	// middlewares
	api.Use(rest.DefaultProdStack...)
	api.Use(&HeadersMiddleware{})
	api.Use(&CustomAuthBasicMiddleware{
		Realm:          "This autocompletion API requires a basic authentication",
		CheckedMethods: []string{"POST", "PUT", "DELETE"},
		Authenticate: func(username string, password string) bool {
			if username == config.Http.Basic.Username && password == config.Http.Basic.Password {
				return true
			}
			return false
		},
	})
	api.Use(&rest.CorsMiddleware{
		RejectNonCorsRequests: false,
		OriginValidator: func(origin string, request *rest.Request) bool {
			return true
		},
		AllowedMethods: []string{"GET", "OPTIONS", "HEAD"},
		AllowedHeaders: []string{
			"Accept", "Content-Type", "X-Custom-Header", "Origin", "Cache-Control"},
		AccessControlAllowCredentials: true,
		AccessControlMaxAge:           3600,
	})


	router, err := GetRouter()
	if err != nil {
		log.Fatal(err)
	}

	api.SetApp(router)

	port := strconv.Itoa(config.Http.Port.(int))
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, api.MakeHandler()))
}
