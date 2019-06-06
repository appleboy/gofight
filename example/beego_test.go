package example

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/appleboy/gofight/v2"

	"github.com/astaxie/beego"
	"github.com/stretchr/testify/assert"
)

func TestSayHelloWorld(t *testing.T) {
	uri := "/say"

	// LoadAppConfig allow developer to apply a config file
	// beego.LoadAppConfig("ini", "../conf/app.conf")
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
