package example

import (
	"net/http"
	"testing"

	"github.com/appleboy/gofight/v2"

	"github.com/stretchr/testify/assert"
)

func TestGinHelloWorld(t *testing.T) {
	r := gofight.New()

	r.GET("/").
		SetDebug(true).
		Run(GinEngine(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "Hello World", r.Body.String())
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestGinHelloHandler(t *testing.T) {
	r := gofight.New()

	r.GET("/user/appleboy").
		SetDebug(true).
		Run(GinEngine(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "Hello, appleboy", r.Body.String())
			assert.Equal(t, http.StatusOK, r.Code)
		})
}
