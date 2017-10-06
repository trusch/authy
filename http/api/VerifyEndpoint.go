package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/trusch/authy/jwt"
)

// VerifyEndpoint is an http.Handler which serves CRUD operations for schemas
type VerifyEndpoint struct {
	publicKey interface{}
}

// ServeHTTP is the function needed to implement http.Handler
func (endpoint *VerifyEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Infof("verify request: %v %v", r.Method, r.URL)
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	claims, err := jwt.ValidateToken(string(bs), endpoint.publicKey)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(claims)
}

// NewVerifyEndpoint constructs a new handler instances
func NewVerifyEndpoint(publicKey interface{}) http.Handler {
	endpoint := &VerifyEndpoint{publicKey}
	return endpoint
}
