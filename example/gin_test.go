package example

import (
	"github.com/appleboy/mocha"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGinHelloWorld(t *testing.T) {
	r := mocha.New()

	r.GET("/").
		SetDebug(true).
		RunGin(GinEngine(), func(r *httptest.ResponseRecorder) {
			assert.Equal(t, "Hello World", r.Body.String())
			assert.Equal(t, http.StatusOK, r.Code)
		})
}
