package example

import (
	"github.com/astaxie/beego"
)

// UserController for beego router
type UserController struct {
	beego.Controller
}

// SayHelloWorld for say hello
func (c *UserController) SayHelloWorld() {
	c.Ctx.ResponseWriter.Status = 200
	c.Ctx.ResponseWriter.Write([]byte("Hello, World"))
}
