package example

import (
	"github.com/appleboy/mocha"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestEchoHelloWorld(t *testing.T) {
	r := mocha.New()

	r.GET("/").
		SetDebug(true).
		RunEcho(EchoEngine(), func(r mocha.EchoHttpResponse, rq mocha.EchoHttpRequest) {
			assert.Equal(t, "Hello World", r.Body.String())
			assert.Equal(t, http.StatusOK, r.Status())
		})
}
