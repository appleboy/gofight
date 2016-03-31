package example

import (
	"github.com/appleboy/gofight"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGinHelloWorld(t *testing.T) {
	r := gofight.New()

	r.GET("/").
		SetDebug(true).
		Run(GinEngine(), func(r gofight.HttpResponse, rq gofight.HttpRequest) {
			assert.Equal(t, "Hello World", r.Body.String())
			assert.Equal(t, http.StatusOK, r.Code)
		})
}
