package api

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/trusch/authy/jwt"
	"github.com/trusch/streamstore"
	"github.com/xeipuuv/gojsonschema"
	"golang.org/x/crypto/bcrypt"
)

// LoginEndpoint is an http.Handler which serves CRUD operations for schemas
type LoginEndpoint struct {
	store      streamstore.Storage
	privateKey interface{}
}

// LoginRequest is the struture of a Login Request (haha)
type LoginRequest struct {
	ID       string     `json:"id"`
	Password string     `json:"password"`
	Claims   jwt.Claims `json:"claims"`
}

// ServeHTTP is the function needed to implement http.Handler
func (endpoint *LoginEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Infof("login request: %v %v", r.Method, r.URL)
	decoder := json.NewDecoder(r.Body)
	req := &LoginRequest{}
	if err := decoder.Decode(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	userReader, err := endpoint.store.GetReader("user::" + req.ID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	defer userReader.Close()
	user := &struct {
		ID           string `json:"id"`
		Password     string `json:"password"`
		AuthSchemaID string `json:"authschema"`
	}{}
	decoder = json.NewDecoder(userReader)
	if err = decoder.Decode(user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	hash, err := base64.StdEncoding.DecodeString(user.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if err = bcrypt.CompareHashAndPassword(hash, []byte(req.Password)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return
	}
	schemaReader, err := endpoint.store.GetReader("schema::" + user.AuthSchemaID)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("no 'authschema' supplied"))
		return
	}
	defer schemaReader.Close()
	bs, err := ioutil.ReadAll(schemaReader)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	loader := gojsonschema.NewBytesLoader(bs)
	claimLoader := gojsonschema.NewGoLoader(req.Claims)
	if res, e := gojsonschema.Validate(loader, claimLoader); e != nil || !res.Valid() {
		log.Print(e)
		w.WriteHeader(http.StatusUnauthorized)
		if e != nil {
			w.Write([]byte(e.Error()))
			return
		}
		for _, e := range res.Errors() {
			w.Write([]byte(e.String()))
		}
		return
	}
	token, err := jwt.CreateToken(req.Claims, endpoint.privateKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(token))
}

// NewLoginEndpoint constructs a new handler instances
func NewLoginEndpoint(store streamstore.Storage, privateKey interface{}) http.Handler {
	endpoint := &LoginEndpoint{store, privateKey}
	return endpoint
}

func (endpoint *LoginEndpoint) loadUser(username string) (map[string]interface{}, error) {
	reader, err := endpoint.store.GetReader("user::" + username)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	obj := make(map[string]interface{})
	decoder := json.NewDecoder(reader)
	return obj, decoder.Decode(&obj)
}
