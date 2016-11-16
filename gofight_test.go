package gofight

import (
	"net/http"
	"runtime"
	"testing"
	"time"

	"github.com/appleboy/gofight/framework"
	"github.com/buger/jsonparser"
	"github.com/stretchr/testify/assert"
	"gopkg.in/gin-gonic/gin.v1"
)

var goVersion = runtime.Version()

func TestHttpURL(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/example", func(c *gin.Context) { c.String(http.StatusOK, "it worked") })

	go func() {
		assert.NoError(t, router.Run())
	}()
	// have to wait for the goroutine to start and run the server
	// otherwise the main thread will complete
	time.Sleep(5 * time.Millisecond)

	TestRequest(t, "http://localhost:8080/example")
}

func TestHttpsURL(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/example", func(c *gin.Context) { c.String(http.StatusOK, "it worked") })

	go func() {
		assert.NoError(t, router.RunTLS(":8088", "certificate/localhost.cert", "certificate/localhost.key"))
	}()
	// have to wait for the goroutine to start and run the server
	// otherwise the main thread will complete
	time.Sleep(5 * time.Millisecond)

	TestRequest(t, "https://localhost:8088/example")
}

func TestGinHelloWorld(t *testing.T) {
	r := New()

	r.GET("/hello").
		SetDebug(true).
		Run(framework.GinEngine(), func(r HTTPResponse, rq HTTPRequest) {
			data := []byte(r.Body.String())

			value, _ := jsonparser.GetString(data, "hello")

			assert.Equal(t, "world", value)
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestGinHeader(t *testing.T) {
	r := New()

	r.GET("/text").
		SetHeader(H{
			"Content-Type": "text/plain",
			"Go-Version":   goVersion,
		}).
		Run(framework.GinEngine(), func(r HTTPResponse, rq HTTPRequest) {

			assert.Equal(t, goVersion, rq.Header.Get("Go-Version"))
			assert.Equal(t, "Gofight-client/"+Version, rq.Header.Get(UserAgent))
			assert.Equal(t, "text/plain", rq.Header.Get(ContentType))
			assert.Equal(t, "Hello World", r.Body.String())
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestGinCookie(t *testing.T) {
	r := New()

	r.GET("/text").
		SetCookie(H{
			"foo": "bar",
		}).
		Run(framework.GinEngine(), func(r HTTPResponse, rq HTTPRequest) {

			assert.Equal(t, http.StatusOK, r.Code)
			assert.Equal(t, "foo=bar", rq.Header.Get("cookie"))
		})
}

func TestGinQuery(t *testing.T) {
	r := New()

	r.GET("/query?text=world&foo=bar").
		Run(framework.GinEngine(), func(r HTTPResponse, rq HTTPRequest) {

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
		SetForm(H{
			"a": "1",
			"b": "2",
		}).
		Run(framework.GinEngine(), func(r HTTPResponse, rq HTTPRequest) {
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
		Run(framework.GinEngine(), func(r HTTPResponse, rq HTTPRequest) {
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
		Run(framework.GinEngine(), func(r HTTPResponse, rq HTTPRequest) {
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
		Run(framework.GinEngine(), func(r HTTPResponse, rq HTTPRequest) {
			data := []byte(r.Body.String())

			hello, _ := jsonparser.GetString(data, "hello")

			assert.Equal(t, "world", hello)
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestGinPatch(t *testing.T) {
	r := New()

	r.PATCH("/patch").
		SetJSON(D{
			"a": 1,
			"b": 2,
		}).
		Run(framework.GinEngine(), func(r HTTPResponse, rq HTTPRequest) {
			data := []byte(r.Body.String())

			a, _ := jsonparser.GetInt(data, "a")
			b, _ := jsonparser.GetInt(data, "b")

			assert.Equal(t, 1, int(a))
			assert.Equal(t, 2, int(b))
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestGinHead(t *testing.T) {
	r := New()

	r.HEAD("/head").
		SetJSON(D{
			"a": 1,
			"b": 2,
		}).
		Run(framework.GinEngine(), func(r HTTPResponse, rq HTTPRequest) {
			data := []byte(r.Body.String())

			a, _ := jsonparser.GetInt(data, "a")
			b, _ := jsonparser.GetInt(data, "b")

			assert.Equal(t, 1, int(a))
			assert.Equal(t, 2, int(b))
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestGinOptions(t *testing.T) {
	r := New()

	r.OPTIONS("/options").
		SetJSON(D{
			"a": 1,
			"b": 2,
		}).
		Run(framework.GinEngine(), func(r HTTPResponse, rq HTTPRequest) {
			data := []byte(r.Body.String())

			a, _ := jsonparser.GetInt(data, "a")
			b, _ := jsonparser.GetInt(data, "b")

			assert.Equal(t, 1, int(a))
			assert.Equal(t, 2, int(b))
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestEchoHelloWorld(t *testing.T) {
	r := New()

	r.GET("/hello").
		SetDebug(true).
		Run(framework.EchoEngine(), func(r HTTPResponse, rq HTTPRequest) {
			data := []byte(r.Body.String())

			value, _ := jsonparser.GetString(data, "hello")

			assert.Equal(t, "world", value)
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestEchoHeader(t *testing.T) {
	r := New()

	r.GET("/text").
		SetHeader(H{
			"Content-Type": "text/plain",
			"Go-Version":   goVersion,
		}).
		Run(framework.EchoEngine(), func(r HTTPResponse, rq HTTPRequest) {

			assert.Equal(t, goVersion, rq.Header.Get("Go-Version"))
			assert.Equal(t, r.Body.String(), "Hello World")
			assert.Equal(t, r.Code, http.StatusOK)
		})
}

func TestEchoQuery(t *testing.T) {
	r := New()

	r.GET("/query?text=world&foo=bar").
		Run(framework.EchoEngine(), func(r HTTPResponse, rq HTTPRequest) {
			data := []byte(r.Body.String())

			hello, _ := jsonparser.GetString(data, "hello")
			foo, _ := jsonparser.GetString(data, "foo")

			assert.Equal(t, "world", hello)
			assert.Equal(t, "bar", foo)
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestEchoPostFormData(t *testing.T) {
	r := New()

	r.POST("/form").
		SetBody("a=1&b=2").
		Run(framework.EchoEngine(), func(r HTTPResponse, rq HTTPRequest) {
			data := []byte(r.Body.String())

			a, _ := jsonparser.GetString(data, "a")
			b, _ := jsonparser.GetString(data, "b")

			assert.Equal(t, "1", a)
			assert.Equal(t, "2", b)
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestEchoPostJSONData(t *testing.T) {
	r := New()

	r.POST("/json").
		SetJSON(D{
			"a": 1,
			"b": 2,
		}).
		Run(framework.EchoEngine(), func(r HTTPResponse, rq HTTPRequest) {
			data := []byte(r.Body.String())

			a, _ := jsonparser.GetInt(data, "a")
			b, _ := jsonparser.GetInt(data, "b")

			assert.Equal(t, 1, int(a))
			assert.Equal(t, 2, int(b))
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestEchoPut(t *testing.T) {
	r := New()

	r.PUT("/update").
		SetBody("c=1&d=2").
		Run(framework.EchoEngine(), func(r HTTPResponse, rq HTTPRequest) {
			data := []byte(r.Body.String())

			c, _ := jsonparser.GetString(data, "c")
			d, _ := jsonparser.GetString(data, "d")

			assert.Equal(t, "1", c)
			assert.Equal(t, "2", d)
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestEchoDelete(t *testing.T) {
	r := New()

	r.DELETE("/delete").
		Run(framework.EchoEngine(), func(r HTTPResponse, rq HTTPRequest) {
			data := []byte(r.Body.String())

			hello, _ := jsonparser.GetString(data, "hello")

			assert.Equal(t, "world", hello)
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestEchoPatch(t *testing.T) {
	r := New()

	r.PATCH("/patch").
		Run(framework.EchoEngine(), func(r HTTPResponse, rq HTTPRequest) {
			data := []byte(r.Body.String())

			value, _ := jsonparser.GetString(data, "hello")

			assert.Equal(t, "world", value)
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestEchoHead(t *testing.T) {
	r := New()

	r.HEAD("/head").
		Run(framework.EchoEngine(), func(r HTTPResponse, rq HTTPRequest) {
			data := []byte(r.Body.String())

			value, _ := jsonparser.GetString(data, "hello")

			assert.Equal(t, "world", value)
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestEchoOptions(t *testing.T) {
	r := New()

	r.OPTIONS("/options").
		Run(framework.EchoEngine(), func(r HTTPResponse, rq HTTPRequest) {
			data := []byte(r.Body.String())

			value, _ := jsonparser.GetString(data, "hello")

			assert.Equal(t, "world", value)
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestSetQueryString(t *testing.T) {
	r := New()

	r.GET("/hello").
		SetQuery(H{
			"a": "1",
			"b": "2",
		})

	assert.Equal(t, "/hello?a=1&b=2", r.Path)

	r.GET("/hello?foo=bar").
		SetQuery(H{
			"a": "1",
			"b": "2",
		})

	assert.Equal(t, "/hello?foo=bar&a=1&b=2", r.Path)
}
