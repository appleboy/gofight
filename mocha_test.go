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

func TestHelloWorld(t *testing.T) {
	r := RequestConfig{
		Method:  "GET",
		Path:    "/hello",
		Body:    `{}`,
		Handler: helloHandler,
		Callback: func(r *httptest.ResponseRecorder) {
			data := []byte(r.Body.String())

			value, _ := jsonparser.GetString(data, "hello")

			assert.Equal(t, value, "world")
			assert.Equal(t, r.Code, http.StatusOK)
		},
	}

	RunRequest(r)
}
