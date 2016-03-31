package mocha

import (
	"github.com/appleboy/mocha/framework"
	"github.com/buger/jsonparser"
	"github.com/stretchr/testify/assert"
	"net/http"
	"runtime"
	"testing"
)

var go_version = runtime.Version()

func TestGinHelloWorld(t *testing.T) {
	r := New()

	r.GET("/hello").
		SetDebug(true).
		Run(framework.GinEngine(), func(r HttpResponse, rq HttpRequest) {
			data := []byte(r.Body.String())

			value, _ := jsonparser.GetString(data, "hello")

			assert.Equal(t, "world", value)
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestGinHeader(t *testing.T) {
	r := New()

	go_version := runtime.Version()

	r.GET("/text").
		SetHeader(H{
			"Content-Type": "text/plain",
			"Go-Version":   go_version,
		}).
		Run(framework.GinEngine(), func(r HttpResponse, rq HttpRequest) {

			assert.Equal(t, go_version, rq.Header.Get("Go-Version"))
			assert.Equal(t, "Hello World", r.Body.String())
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestGinQuery(t *testing.T) {
	r := New()

	r.GET("/query?text=world&foo=bar").
		Run(framework.GinEngine(), func(r HttpResponse, rq HttpRequest) {

			data := []byte(r.Body.String())

			hello, _ := jsonparser.GetString(data, "hello")
			foo, _ := jsonparser.GetString(data, "foo")

			assert.Equal(t, "world", hello)
			assert.Equal(t, "bar", foo)
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestGinPostFormData(t *testing.T) {
	r := New()

	r.POST("/form").
		SetFORM(H{
			"a": "1",
			"b": "2",
		}).
		Run(framework.GinEngine(), func(r HttpResponse, rq HttpRequest) {
			data := []byte(r.Body.String())

			a, _ := jsonparser.GetString(data, "a")
			b, _ := jsonparser.GetString(data, "b")

			assert.Equal(t, "1", a)
			assert.Equal(t, "2", b)
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestGinPostJSONData(t *testing.T) {
	r := New()

	r.POST("/json").
		SetJSON(D{
			"a": 1,
			"b": 2,
		}).
		Run(framework.GinEngine(), func(r HttpResponse, rq HttpRequest) {
			data := []byte(r.Body.String())

			a, _ := jsonparser.GetInt(data, "a")
			b, _ := jsonparser.GetInt(data, "b")

			assert.Equal(t, 1, int(a))
			assert.Equal(t, 2, int(b))
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestGinPut(t *testing.T) {
	r := New()

	r.PUT("/update").
		SetBody("c=1&d=2").
		Run(framework.GinEngine(), func(r HttpResponse, rq HttpRequest) {
			data := []byte(r.Body.String())

			c, _ := jsonparser.GetString(data, "c")
			d, _ := jsonparser.GetString(data, "d")

			assert.Equal(t, "1", c)
			assert.Equal(t, "2", d)
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestGinDelete(t *testing.T) {
	r := New()

	r.DELETE("/delete").
		Run(framework.GinEngine(), func(r HttpResponse, rq HttpRequest) {
			data := []byte(r.Body.String())

			hello, _ := jsonparser.GetString(data, "hello")

			assert.Equal(t, "world", hello)
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestEchoHelloWorld(t *testing.T) {
	r := New()

	r.GET("/hello").
		SetDebug(true).
		RunEcho(framework.EchoEngine(), func(r EchoHttpResponse, rq EchoHttpRequest) {
			data := []byte(r.Body.String())

			value, _ := jsonparser.GetString(data, "hello")

			assert.Equal(t, "world", value)
			assert.Equal(t, http.StatusOK, r.Status())
		})
}

func TestEchoHeader(t *testing.T) {
	r := New()

	r.GET("/text").
		SetHeader(H{
			"Content-Type": "text/plain",
			"Go-Version":   go_version,
		}).
		RunEcho(framework.EchoEngine(), func(r EchoHttpResponse, rq EchoHttpRequest) {

			assert.Equal(t, go_version, rq.Header().Get("Go-Version"))
			assert.Equal(t, r.Body.String(), "Hello World")
			assert.Equal(t, r.Status(), http.StatusOK)
		})
}

func TestEchoQuery(t *testing.T) {
	r := New()

	r.GET("/query?text=world&foo=bar").
		RunEcho(framework.EchoEngine(), func(r EchoHttpResponse, rq EchoHttpRequest) {
			data := []byte(r.Body.String())

			hello, _ := jsonparser.GetString(data, "hello")
			foo, _ := jsonparser.GetString(data, "foo")

			assert.Equal(t, "world", hello)
			assert.Equal(t, "bar", foo)
			assert.Equal(t, http.StatusOK, r.Status())
		})
}

func TestEchoPostFormData(t *testing.T) {
	r := New()

	r.POST("/form").
		SetBody("a=1&b=2").
		RunEcho(framework.EchoEngine(), func(r EchoHttpResponse, rq EchoHttpRequest) {
			data := []byte(r.Body.String())

			a, _ := jsonparser.GetString(data, "a")
			b, _ := jsonparser.GetString(data, "b")

			assert.Equal(t, "1", a)
			assert.Equal(t, "2", b)
			assert.Equal(t, http.StatusOK, r.Status())
		})
}

func TestEchoPostJSONData(t *testing.T) {
	r := New()

	r.POST("/json").
		SetJSON(D{
			"a": 1,
			"b": 2,
		}).
		RunEcho(framework.EchoEngine(), func(r EchoHttpResponse, rq EchoHttpRequest) {
			data := []byte(r.Body.String())

			a, _ := jsonparser.GetInt(data, "a")
			b, _ := jsonparser.GetInt(data, "b")

			assert.Equal(t, 1, int(a))
			assert.Equal(t, 2, int(b))
			assert.Equal(t, http.StatusOK, r.Status())
		})
}

func TestEchoPut(t *testing.T) {
	r := New()

	r.PUT("/update").
		SetBody("c=1&d=2").
		RunEcho(framework.EchoEngine(), func(r EchoHttpResponse, rq EchoHttpRequest) {
			data := []byte(r.Body.String())

			c, _ := jsonparser.GetString(data, "c")
			d, _ := jsonparser.GetString(data, "d")

			assert.Equal(t, "1", c)
			assert.Equal(t, "2", d)
			assert.Equal(t, http.StatusOK, r.Status())
		})
}

func TestEchoDelete(t *testing.T) {
	r := New()

	r.DELETE("/delete").
		RunEcho(framework.EchoEngine(), func(r EchoHttpResponse, rq EchoHttpRequest) {
			data := []byte(r.Body.String())

			hello, _ := jsonparser.GetString(data, "hello")

			assert.Equal(t, "world", hello)
			assert.Equal(t, http.StatusOK, r.Status())
		})
}
