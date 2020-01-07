package example

import (
	"net/http"

	"github.com/gorilla/pat"
)

func patHelloHandler(wr http.ResponseWriter, req *http.Request) {
	wr.WriteHeader(http.StatusOK)
	wr.Write([]byte("Hello World"))
}

func patUserHandler(wr http.ResponseWriter, req *http.Request) {
	name := req.URL.Query().Get(":name")
	wr.WriteHeader(http.StatusOK)
	wr.Write([]byte("Hello, " + name))
}

// PatEngine is pat router.
func PatEngine() *pat.Router {
	r := pat.New()

	r.Get("/user/{name}", patUserHandler)
	r.Get("/", patHelloHandler)

	return r
}
