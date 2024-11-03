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

func basicCookieHandler(w http.ResponseWriter, r *http.Request) {
	// get cookie from request.
	foo, err := r.Cookie("foo")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, _ = io.WriteString(w, foo.Value)
}

func basicQueryHandler(w http.ResponseWriter, r *http.Request) {
	// get query from request.
	foo := r.URL.Query().Get("foo")
	_, _ = io.WriteString(w, foo)
}

func basicFormHandler(w http.ResponseWriter, r *http.Request) {
	// get form from request.
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	foo := r.Form.Get("foo")
	_, _ = io.WriteString(w, foo)
}

func basicEngine() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", basicHelloHandler)
	mux.HandleFunc("/cookie", basicCookieHandler)
	mux.HandleFunc("/query", basicQueryHandler)
	mux.HandleFunc("/form", basicFormHandler)

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

func TestSetCookie(t *testing.T) {
	r := New()
	cookies := H{
		"foo": "bar",
		"baz": "qux",
	}

	r.GET("/cookie").
		SetCookie(cookies).
		Run(basicEngine(), func(r HTTPResponse, rq HTTPRequest) {
			cookieFoo, err := rq.Cookie("foo")
			assert.NoError(t, err)
			assert.Equal(t, "bar", cookieFoo.Value)

			cookieBaz, err := rq.Cookie("baz")
			assert.NoError(t, err)
			assert.Equal(t, "qux", cookieBaz.Value)

			assert.Equal(t, "bar", r.Body.String())
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestSetQuery(t *testing.T) {
	r := New()
	query := H{
		"foo": "bar",
	}

	r.GET("/query").
		SetQuery(query).
		Run(basicEngine(), func(r HTTPResponse, rq HTTPRequest) {
			assert.Equal(t, "bar", rq.URL.Query().Get("foo"))
			assert.Equal(t, "bar", r.Body.String())
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestSetQueryWithExistingQuery(t *testing.T) {
	r := New()
	query := H{
		"c": "3",
		"d": "4",
	}

	r.GET("/query?a=1&b=2&foo=testing").
		SetQuery(query).
		Run(basicEngine(), func(r HTTPResponse, rq HTTPRequest) {
			assert.Equal(t, "1", rq.URL.Query().Get("a"))
			assert.Equal(t, "2", rq.URL.Query().Get("b"))
			assert.Equal(t, "3", rq.URL.Query().Get("c"))
			assert.Equal(t, "4", rq.URL.Query().Get("d"))
			assert.Equal(t, "testing", rq.URL.Query().Get("foo"))
			assert.Equal(t, "testing", r.Body.String())
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestSetForm(t *testing.T) {
	r := New()
	formData := H{
		"a":   "1",
		"b":   "2",
		"foo": "bar",
	}

	r.POST("/form").
		SetForm(formData).
		Run(basicEngine(), func(r HTTPResponse, rq HTTPRequest) {
			assert.Equal(t, "bar", r.Body.String())
			assert.Equal(t, http.StatusOK, r.Code)
		})
}
