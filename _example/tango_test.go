package example

import (
	"fmt"
	"net/http"
	"testing"

	"gitea.com/lunny/tango"
	"github.com/appleboy/gofight/v2"
	"github.com/stretchr/testify/assert"
)

func TestTangoHelloWorld(t *testing.T) {
	tan := tango.New(tango.Return())
	tan.Get("/", new(HelloAction))

	r := gofight.New()
	r.GET("/").
		SetDebug(true).
		Run(tan, func(rp gofight.HTTPResponse, rq gofight.HTTPRequest) {
			fmt.Println(rp.Code)
			assert.Equal(t, "Hello, World", rp.Body.String())
			assert.Equal(t, http.StatusOK, rp.Code)
		})
}
