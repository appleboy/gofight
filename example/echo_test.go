package example

import (
	"github.com/appleboy/mocha"
	"github.com/labstack/echo/engine"
	"github.com/labstack/echo/test"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestEchoHelloWorld(t *testing.T) {
	r := mocha.New()

	r.GET("/").
		SetDebug(true).
		RunEcho(EchoEngine(), func(r *test.ResponseRecorder, rq engine.Request) {
			assert.Equal(t, "Hello World", r.Body.String())
			assert.Equal(t, http.StatusOK, r.Status())
		})
}
