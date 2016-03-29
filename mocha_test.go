package mocha

import (
	"github.com/buger/jsonparser"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func helloHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "world",
	})
}

func textHandler(c *gin.Context) {
	c.String(http.StatusOK, "Hello World")
}

func queryHandler(c *gin.Context) {
	text := c.Query("text")
	foo := c.Query("foo")

	c.JSON(http.StatusOK, gin.H{
		"hello": text,
		"foo": foo,
	})
}

func TestHelloWorld(t *testing.T) {
	r := &RequestConfig{
		Handler: helloHandler,
		Callback: func(r *httptest.ResponseRecorder) {
			data := []byte(r.Body.String())

			value, _ := jsonparser.GetString(data, "hello")

			assert.Equal(t, value, "world")
			assert.Equal(t, r.Code, http.StatusOK)
		},
	}

	r.SetDebug(true).Run()
}

func TestHeader(t *testing.T) {
	r := &RequestConfig{
		Handler: textHandler,
		Headers: map[string]string {
			"Content-Type": "text/plain",
		},
		Callback: func(r *httptest.ResponseRecorder) {

			assert.Equal(t, r.Body.String(), "Hello World")
			assert.Equal(t, r.Code, http.StatusOK)
		},
	}

	r.Run()
}

func TestQuery(t *testing.T) {
	r := &RequestConfig{
		Path: "/hello?text=world&foo=bar",
		Handler: queryHandler,
		Callback: func(r *httptest.ResponseRecorder) {
			data := []byte(r.Body.String())

			hello, _ := jsonparser.GetString(data, "hello")
			foo, _ := jsonparser.GetString(data, "foo")

			assert.Equal(t, hello, "world")
			assert.Equal(t, foo, "bar")
			assert.Equal(t, r.Code, http.StatusOK)
		},
	}

	r.Run()
}
