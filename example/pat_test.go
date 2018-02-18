package example

import (
	"github.com/appleboy/gofight"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestPatHelloWorld(t *testing.T) {
	r := gofight.New()

	r.GET("/").
		SetDebug(true).
		Run(PatEngine(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "Hello World", r.Body.String())
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestPatHelloHandler(t *testing.T) {
	r := gofight.New()

	r.GET("/user/appleboy").
		SetDebug(true).
		Run(GinEngine(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "Hello, appleboy", r.Body.String())
			assert.Equal(t, http.StatusOK, r.Code)
		})
}
