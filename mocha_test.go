package mocha

import (
	"github.com/appleboy/mocha/framework"
	"github.com/buger/jsonparser"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGinHelloWorld(t *testing.T) {
	r := New()

	r.GET("/hello").
		SetDebug(true).
		RunGin(framework.GinEngine(), func(r *httptest.ResponseRecorder) {
			data := []byte(r.Body.String())

			value, _ := jsonparser.GetString(data, "hello")

			assert.Equal(t, value, "world")
			assert.Equal(t, r.Code, http.StatusOK)
		})
}

func TestGinHeader(t *testing.T) {
	r := New()

	r.GET("/text").
		SetHeader(map[string]string{
			"Content-Type": "text/plain",
		}).
		RunGin(framework.GinEngine(), func(r *httptest.ResponseRecorder) {

			assert.Equal(t, r.Body.String(), "Hello World")
			assert.Equal(t, r.Code, http.StatusOK)
		})
}

func TestGinQuery(t *testing.T) {
	r := New()

	r.GET("/query?text=world&foo=bar").
		RunGin(framework.GinEngine(), func(r *httptest.ResponseRecorder) {

			data := []byte(r.Body.String())

			hello, _ := jsonparser.GetString(data, "hello")
			foo, _ := jsonparser.GetString(data, "foo")

			assert.Equal(t, hello, "world")
			assert.Equal(t, foo, "bar")
			assert.Equal(t, r.Code, http.StatusOK)
		})
}

func TestGinPostFormData(t *testing.T) {
	r := New()

	r.POST("/form").
		SetBody("a=1&b=2").
		RunGin(framework.GinEngine(), func(r *httptest.ResponseRecorder) {
			data := []byte(r.Body.String())

			a, _ := jsonparser.GetString(data, "a")
			b, _ := jsonparser.GetString(data, "b")

			assert.Equal(t, a, "1")
			assert.Equal(t, b, "2")
			assert.Equal(t, r.Code, http.StatusOK)
		})
}

func TestGinPostJSONData(t *testing.T) {
	r := New()

	r.POST("/json").
		SetBody(`{"a":1,"b":2}`).
		RunGin(framework.GinEngine(), func(r *httptest.ResponseRecorder) {
			data := []byte(r.Body.String())

			a, _ := jsonparser.GetInt(data, "a")
			b, _ := jsonparser.GetInt(data, "b")

			assert.Equal(t, int(a), 1)
			assert.Equal(t, int(b), 2)
			assert.Equal(t, r.Code, http.StatusOK)
		})
}

func TestGinPut(t *testing.T) {
	r := New()

	r.PUT("/update").
		SetBody("c=1&d=2").
		RunGin(framework.GinEngine(), func(r *httptest.ResponseRecorder) {
			data := []byte(r.Body.String())

			c, _ := jsonparser.GetString(data, "c")
			d, _ := jsonparser.GetString(data, "d")

			assert.Equal(t, c, "1")
			assert.Equal(t, d, "2")
			assert.Equal(t, r.Code, http.StatusOK)
		})
}

func TestGinDelete(t *testing.T) {
	r := New()

	r.DELETE("/delete").
		RunGin(framework.GinEngine(), func(r *httptest.ResponseRecorder) {
			data := []byte(r.Body.String())

			hello, _ := jsonparser.GetString(data, "hello")

			assert.Equal(t, hello, "world")
			assert.Equal(t, r.Code, http.StatusOK)
		})
}
