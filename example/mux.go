package example

import (
	"github.com/gorilla/mux"
	"net/http"
)

func MuxHelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

func MuxEngine() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", MuxHelloHandler)

	return r
}
