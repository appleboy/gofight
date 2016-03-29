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

	r.Run()
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
