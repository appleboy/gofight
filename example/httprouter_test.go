package example

import (
	"net/http"
	"testing"

	"github.com/appleboy/gofight/v2"

	"github.com/stretchr/testify/assert"
)

func TestHttpRouterHelloWorld(t *testing.T) {
	r := gofight.New()

	r.GET("/").
		Run(HTTPRouterEngine(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "Hello World", r.Body.String())
			assert.Equal(t, http.StatusOK, r.Code)
		})
}
