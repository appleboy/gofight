package example

import (
	"io"
	"net/http"
)

func basicHelloHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello World")
}

// BasicEngine is basic router.
func BasicEngine() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", basicHelloHandler)

	return mux
}
