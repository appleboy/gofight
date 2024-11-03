package gofight

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const version = "0.0.1"

func basicHelloHandler(w http.ResponseWriter, r *http.Request) {
	// add header in response.
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("X-Version", version)
	_, _ = io.WriteString(w, "Hello World")
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

func TestSetContext(t *testing.T) {
	r := New()
	type contextKey string
	const key contextKey = "key"
	ctx := context.WithValue(context.Background(), key, "value")

	r.GET("/").
		SetContext(ctx).
		Run(basicEngine(), func(r HTTPResponse, rq HTTPRequest) {
			assert.Equal(t, "value", rq.Context().Value(contextKey("key")))
			assert.Equal(t, "Hello World", r.Body.String())
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestSetContextWithTimeout(t *testing.T) {
	r := New()
	ctx, cancel := context.WithTimeout(context.Background(), 0)
	defer cancel()

	r.GET("/").
		SetContext(ctx).
		Run(basicEngine(), func(r HTTPResponse, rq HTTPRequest) {
			select {
			case <-rq.Context().Done():
				assert.Equal(t, context.DeadlineExceeded, rq.Context().Err())
			default:
				t.Error("expected context to be done")
			}
		})
}

func TestGetHeaderFromResponse(t *testing.T) {
	version := "0.0.1"
	r := New()
	r.GET("/").
		Run(basicEngine(), func(r HTTPResponse, rq HTTPRequest) {
			assert.Equal(t, version, r.Header().Get("X-Version"))
			assert.Equal(t, "Hello World", r.Body.String())
		})
}

func TestSetBody(t *testing.T) {
	r := New()
	body := "a=1&b=2"

	r.POST("/").
		SetBody(body).
		Run(basicEngine(), func(r HTTPResponse, rq HTTPRequest) {
			// Read the content of the io.Reader
			bodyBytes, err := io.ReadAll(rq.Body)
			if err != nil {
				t.Fatalf("Failed to read body: %v", err)
			}

			// Convert the byte slice to a string
			bodyString := string(bodyBytes)
			assert.Equal(t, body, bodyString)
			assert.Equal(t, "Hello World", r.Body.String())
			assert.Equal(t, http.StatusOK, r.Code)
		})
}
