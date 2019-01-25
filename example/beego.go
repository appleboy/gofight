package example

import (
	"github.com/astaxie/beego"
)

type UserController struct {
	beego.Controller
}

func (c *UserController) SayHelloWorld() {
	c.Ctx.ResponseWriter.Status = 200
	c.Ctx.ResponseWriter.Write([]byte("Hello, World"))
}
