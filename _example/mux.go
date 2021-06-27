package example

import (
	"net/http"

	"github.com/gorilla/mux"
)

func muxHelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

// MuxEngine is mux router.
func MuxEngine() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", muxHelloHandler)

	return r
}
