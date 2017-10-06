package api

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/trusch/crud"
	"github.com/trusch/streamstore"
	"github.com/xeipuuv/gojsonschema"
)

// SchemaEndpoint is an http.Handler which serves CRUD operations for schemas
type SchemaEndpoint struct {
	router       *mux.Router
	store        streamstore.Storage
	prefix       string
	crudEndpoint http.Handler
}

// ServeHTTP is the function needed to implement http.Handler
func (endpoint *SchemaEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Infof("schema request: %v %v", r.Method, r.URL)
	endpoint.router.ServeHTTP(w, r)
	return
}

// NewSchemaEndpoint constructs a new handler instances
func NewSchemaEndpoint(store streamstore.Storage) http.Handler {
	endpoint := &SchemaEndpoint{mux.NewRouter(), store, "schema", crud.NewEndpoint("schema", store)}
	endpoint.router.PathPrefix("/").Methods("POST", "PUT").HandlerFunc(endpoint.handleSaveSchema)
	endpoint.router.PathPrefix("/").Handler(endpoint.crudEndpoint)
	return endpoint
}

func (endpoint *SchemaEndpoint) handleSaveSchema(w http.ResponseWriter, r *http.Request) {
	log.Debug("validate new schema...")
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bs))
	loader := gojsonschema.NewBytesLoader(bs)
	_, err = gojsonschema.NewSchema(loader)
	if err != nil {
		log.Errorf("try to post invalid json schema: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	endpoint.crudEndpoint.ServeHTTP(w, r)
}
