package http

import (
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/trusch/authy/http/api"
	"github.com/trusch/authy/jwt"
	"github.com/trusch/streamstore/uriparser"
)

// Serve constructs the router and starts serving http requests
func Serve(address, webroot, storageURI string, options uriparser.Options, keyfile string) error {
	router := mux.NewRouter().StrictSlash(false)
	store, err := uriparser.NewFromURI(storageURI, options)
	if err != nil {
		return err
	}
	key, err := jwt.LoadPrivateKey(keyfile)
	if err != nil {
		return err
	}

	userRouter := router.PathPrefix("/api/v1/user").Subrouter()
	userRouter.NewRoute().Handler(http.StripPrefix("/api/v1/user", api.NewUserEndpoint(store)))

	schemaRouter := router.PathPrefix("/api/v1/schema").Subrouter()
	schemaRouter.NewRoute().Handler(http.StripPrefix("/api/v1/schema", api.NewSchemaEndpoint(store)))

	router.Path("/api/v1/login").Handler(api.NewLoginEndpoint(store, key))
	router.Path("/index.html").Handler(rewriteToIndexHandler(filepath.Join(webroot, "static/index.html")))
	router.PathPrefix("/static/").Handler(http.FileServer(http.Dir(webroot)))
	router.NotFoundHandler = rewriteToIndexHandler(filepath.Join(webroot, "static/index.html"))

	return http.ListenAndServe(address, router)
}

type rewriteToIndexHandler string

func (h rewriteToIndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, string(h))
}
