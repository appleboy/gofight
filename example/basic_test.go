package example

import (
	"github.com/appleboy/mocha"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestBasicHelloWorld(t *testing.T) {
	r := mocha.New()
	version := "0.0.1"

	r.GET("/").
		SetDebug(true).
		SetHeader(mocha.H{
			"X-Version": version,
		}).
		Run(BasicEngine(), func(r mocha.HttpResponse, rq mocha.HttpRequest) {

			assert.Equal(t, version, rq.Header.Get("X-Version"))
			assert.Equal(t, "Hello World", r.Body.String())
			assert.Equal(t, http.StatusOK, r.Code)
		})
}
