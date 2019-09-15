package example

import (
	"gitea.com/lunny/tango"
)

// HelloAction for tango router
type HelloAction struct {
	tango.Ctx
}

func (c *HelloAction) Get() string {
	return "Hello, World"
}
