package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/trusch/crud"
	"github.com/trusch/streamstore"
	"golang.org/x/crypto/bcrypt"
)

// UserEndpoint is an http.Handler which serves CRUD operations for users
// So this is CRUD
// + if payload is json hash the "password" field with bcrypt on all write operations (POST, PUT, PATCH)
// + if payload is json respect the "id" field as object key in POST requests
type UserEndpoint struct {
	router       *mux.Router
	store        streamstore.Storage
	prefix       string
	crudEndpoint http.Handler
}

// ServeHTTP is the function needed to implement http.Handler
func (endpoint *UserEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Infof("user request: %v %v", r.Method, r.URL)
	endpoint.router.ServeHTTP(w, r)
	return
}

// NewUserEndpoint constructs a new handler instances
func NewUserEndpoint(store streamstore.Storage) http.Handler {
	endpoint := &UserEndpoint{mux.NewRouter(), store, "user", crud.NewEndpoint("user", store)}
	endpoint.router.PathPrefix("/").Methods("POST", "PUT", "PATCH").HandlerFunc(endpoint.handleSaveUser)
	endpoint.router.PathPrefix("/").Handler(endpoint.crudEndpoint)
	return endpoint
}

func (endpoint *UserEndpoint) handleSaveUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	obj := make(map[string]interface{})
	if err := decoder.Decode(&obj); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	if pass, ok := obj["password"].(string); ok {
		log.Debug("save user request, hashing password")
		hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		obj["password"] = base64.StdEncoding.EncodeToString(hash)
	}
	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)
	if err := encoder.Encode(obj); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	r.Body = ioutil.NopCloser(buf)
	if id, ok := obj["id"].(string); ok && r.Method == "POST" {
		if endpoint.store.Has("user::" + id) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("user already exists"))
			return
		}
		log.Debug("save user request, rewrite to POST and respect 'id' field")
		r.Method = "PUT"
		r.URL.Path = "/" + id
	}
	endpoint.crudEndpoint.ServeHTTP(w, r)
}
