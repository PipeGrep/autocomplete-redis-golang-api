package main

import (
	"encoding/base64"
	"errors"
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	"net/http"
	"strings"
)

type CustomAuthBasicMiddleware struct {
	Realm          string
	CheckedMethods []string
	Authenticate   func(string, string) bool
}

func (mw *CustomAuthBasicMiddleware) MiddlewareFunc(handler rest.HandlerFunc) rest.HandlerFunc {
	if mw.Realm == "" {
		log.Fatal("Realm is required")
	}

	if mw.Authenticate == nil {
		log.Fatal("Authenticate is required")
	}

	if mw.CheckedMethods == nil {
		mw.CheckedMethods = []string{}
	}

	return func(writer rest.ResponseWriter, request *rest.Request) {
		for i := 0; i < len(mw.CheckedMethods); i++ {
			if request.Method == mw.CheckedMethods[i] {

				authHeader := request.Header.Get("Authorization")
				if authHeader == "" {
					mw.unauthorized(writer)
					return
				}

				providedUsername, providedPassword, err := mw.decodeBasicAuthHeader(authHeader)
				if err != nil {
					rest.Error(writer, "Invalid authentication", http.StatusBadRequest)
					return
				}

				if !mw.Authenticate(providedUsername, providedPassword) {
					mw.unauthorized(writer)
					return
				}
				break
			}
		}

		handler(writer, request)
	}

}

func (mw *CustomAuthBasicMiddleware) unauthorized(writer rest.ResponseWriter) {
	writer.Header().Set("WWW-Authenticate", "Basic realm="+mw.Realm)
	rest.Error(writer, "Not Authorized", http.StatusUnauthorized)
}

func (mw *CustomAuthBasicMiddleware) decodeBasicAuthHeader(header string) (user string, password string, err error) {
	parts := strings.SplitN(header, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Basic") {
		return "", "", errors.New("Invalid authentication")
	}

	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", "", errors.New("Invalid base64")
	}

	creds := strings.SplitN(string(decoded), ":", 2)
	if len(creds) != 2 {
		return "", "", errors.New("Invalid authentication")
	}

	return creds[0], creds[1], nil
}

type HeadersMiddleware struct {
}

// MiddlewareFunc makes StatusMiddleware implement the Middleware interface.
func (mw *HeadersMiddleware) MiddlewareFunc(fn rest.HandlerFunc) rest.HandlerFunc {
	return func(res rest.ResponseWriter, req *rest.Request) {
		res.Header().Set("X-Powered-By", "http://grep.so")
		fn(res, req)
	}
}
