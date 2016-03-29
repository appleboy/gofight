package main

import (
	"github.com/appleboy/gin-mocha"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHelloWorld(t *testing.T) {
	r := &ginMocha.RequestConfig{
		Handler: helloHandler,
		Callback: func(r *httptest.ResponseRecorder) {

			assert.Equal(t, r.Body.String(), "Hello World")
			assert.Equal(t, r.Code, http.StatusOK)
		},
	}

	r.Run()
}
