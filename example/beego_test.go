package example

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/astaxie/beego"
	"github.com/stretchr/testify/assert"

	"github.com/appleboy/gofight"
)

func TestSayHelloWorld(t *testing.T) {
	uri := "/say"

	c := beego.NewControllerRegister()
	c.Add(uri, &UserController{}, "get:SayHelloWorld")

	r := gofight.New()
	r.GET(uri).
		SetDebug(true).
		Run(c, func(rp gofight.HTTPResponse, rq gofight.HTTPRequest) {
			fmt.Println(rp.Code)
			assert.Equal(t, "Hello, World", rp.Body.String())
			assert.Equal(t, http.StatusOK, rp.Code)
		})
}
