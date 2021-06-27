package gofight

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func basicHelloHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello World")
}

func basicEngine() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", basicHelloHandler)

	return mux
}

func TestBasicHelloWorld(t *testing.T) {
	r := New()
	version := "0.0.1"

	r.GET("/").
		// turn on the debug mode.
		SetDebug(true).
		SetHeader(H{
			"X-Version": version,
		}).
		Run(basicEngine(), func(r HTTPResponse, rq HTTPRequest) {
			assert.Equal(t, version, rq.Header.Get("X-Version"))
			assert.Equal(t, "Hello World", r.Body.String())
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func basicHTTPHelloHandler() {
	http.HandleFunc("/hello", basicHelloHandler)
}

func TestBasicHttpHelloWorld(t *testing.T) {
	basicHTTPHelloHandler()

	r := New()

	r.GET("/hello").
		// trun on the debug mode.
		SetDebug(true).
		Run(http.DefaultServeMux, func(r HTTPResponse, rq HTTPRequest) {
			assert.Equal(t, "Hello World", r.Body.String())
			assert.Equal(t, http.StatusOK, r.Code)
		})
}
