package api_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"

	log "github.com/sirupsen/logrus"
)

func TestApi(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Api Suite")
}

func post(handler http.Handler, path, body string) (code int, data string) {
	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", path, bytes.NewBuffer([]byte(body)))
	handler.ServeHTTP(recorder, req)
	code = recorder.Code
	data = recorder.Body.String()
	return
}

func get(handler http.Handler, path string) (code int, data string) {
	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", path, nil)
	handler.ServeHTTP(recorder, req)
	code = recorder.Code
	data = recorder.Body.String()
	return
}

func put(handler http.Handler, path, body string) (code int, data string) {
	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", path, bytes.NewBuffer([]byte(body)))
	handler.ServeHTTP(recorder, req)
	code = recorder.Code
	data = recorder.Body.String()
	return
}

func patch(handler http.Handler, path, body string) (code int, data string) {
	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", path, bytes.NewBuffer([]byte(body)))
	handler.ServeHTTP(recorder, req)
	code = recorder.Code
	data = recorder.Body.String()
	return
}

func del(handler http.Handler, path string) (code int, data string) {
	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", path, nil)
	handler.ServeHTTP(recorder, req)
	code = recorder.Code
	data = recorder.Body.String()
	return
}
