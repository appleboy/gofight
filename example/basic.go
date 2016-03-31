package example

import (
	"io"
	"net/http"
)

func BasicHelloHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello World")
}

func BasicEngine() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", BasicHelloHandler)

	return mux
}
