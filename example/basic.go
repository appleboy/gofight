package example

import (
	"io"
	"net/http"
)

func basicHelloHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello World")
}

func basicHTTPHelloHandler() {
	http.HandleFunc("/hello", basicHelloHandler)
}

// BasicEngine is basic router.
func BasicEngine() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", basicHelloHandler)

	return mux
}
