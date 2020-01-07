package example

import (
	"gitea.com/lunny/tango"
)

// HelloAction for tango router
type HelloAction struct {
	tango.Ctx
}

// Get will be executed when GET requests
func (c *HelloAction) Get() string {
	return "Hello, World"
}
