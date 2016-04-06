package example

import (
	"github.com/gorilla/mux"
	"net/http"
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
