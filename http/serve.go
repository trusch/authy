package http

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/trusch/authy/http/api"
	"github.com/trusch/authy/jwt"
	"github.com/trusch/streamstore/uriparser"
)

// Serve constructs the router and starts serving http requests
func Serve(address, storageURI string, options uriparser.Options, privkeyFile, pubkeyFile string) error {
	router := mux.NewRouter().StrictSlash(false)
	store, err := uriparser.NewFromURI(storageURI, options)
	if err != nil {
		return err
	}
	log.Debugf("created store %v", storageURI)

	privateKey, err := jwt.LoadPrivateKey(privkeyFile)
	if err != nil {
		return err
	}
	log.Debugf("loaded private key %v", privkeyFile)

	publicKey, err := jwt.LoadPublicKey(pubkeyFile)
	if err != nil {
		return err
	}
	log.Debugf("loaded public key %v", pubkeyFile)

	userRouter := router.PathPrefix("/api/v1/user").Subrouter()
	userRouter.NewRoute().Handler(http.StripPrefix("/api/v1/user", api.NewUserEndpoint(store)))

	schemaRouter := router.PathPrefix("/api/v1/schema").Subrouter()
	schemaRouter.NewRoute().Handler(http.StripPrefix("/api/v1/schema", api.NewSchemaEndpoint(store)))

	router.Path("/api/v1/login").Handler(api.NewLoginEndpoint(store, privateKey))
	router.Path("/api/v1/verify").Handler(api.NewVerifyEndpoint(publicKey))

	log.Debugf("start listening on %v", address)
	return http.ListenAndServe(address, router)
}
