package api_test

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"

	. "github.com/trusch/authy/http/api"
	"github.com/trusch/streamstore"
	"github.com/trusch/streamstore/uriparser"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UserEndpoint", func() {
	var (
		recorder *httptest.ResponseRecorder
		handler  http.Handler
		store    streamstore.Storage
		err      error
	)

	BeforeEach(func() {
		recorder = httptest.NewRecorder()
		store, err = uriparser.NewFromURI("file:///tmp/test", nil)
		Expect(err).NotTo(HaveOccurred())
		handler = NewUserEndpoint(store)
	})

	AfterEach(func() {
		os.RemoveAll("/tmp/test")
	})

	It("should be possible to create and retrieve a user", func() {
		content := `{"id":"admin","username":"admin","password":"admin"}`
		code, responseData := post(handler, "/", content)
		Expect(code).To(Equal(http.StatusOK))
		Expect(responseData).NotTo(BeEmpty())
		code, responseData = get(handler, "/admin")
		Expect(code).To(Equal(http.StatusOK))
		Expect(responseData).NotTo(BeEmpty())
		log.Print(responseData)
	})

	It("should be possible to update a user", func() {
		content1 := `{"id":"admin","username":"admin","password":"admin"}`
		content2 := `{"username":"bla"}`

		code, responseData := post(handler, "/", content1)
		Expect(code).To(Equal(http.StatusOK))
		Expect(responseData).NotTo(BeEmpty())

		code, responseData = patch(handler, "/admin", content2)
		Expect(code).To(Equal(http.StatusOK))
		Expect(responseData).NotTo(BeEmpty())

		code, responseData = get(handler, "/admin")
		Expect(code).To(Equal(http.StatusOK))
		Expect(responseData).NotTo(BeEmpty())
		obj := make(map[string]interface{})
		Expect(json.Unmarshal([]byte(responseData), &obj)).To(Succeed())
		Expect(len(obj)).To(Equal(3))
		Expect(obj["id"]).To(Equal("admin"))
		Expect(obj["username"]).To(Equal("bla"))
		Expect(obj["password"]).NotTo(BeEmpty())
	})

})
