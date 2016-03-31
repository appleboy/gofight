package example

import (
	"github.com/appleboy/mocha"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGinHelloWorld(t *testing.T) {
	r := mocha.New()

	r.GET("/").
		SetDebug(true).
		Run(GinEngine(), func(r mocha.HttpResponse, rq mocha.HttpRequest) {
			assert.Equal(t, "Hello World", r.Body.String())
			assert.Equal(t, http.StatusOK, r.Code)
		})
}
