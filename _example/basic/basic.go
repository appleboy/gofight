package example

import (
	"io"
	"net/http"
)

func basicHelloHandler(w http.ResponseWriter, r *http.Request) {
	// add header in response.
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("X-Version", "0.0.1")
	_, _ = io.WriteString(w, "Hello World")
}

// BasicEngine is basic router.
func BasicEngine() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", basicHelloHandler)

	return mux
}
