package example

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ginHelloHandler(c *gin.Context) {
	c.String(http.StatusOK, "Hello World")
}

func ginUserHandler(c *gin.Context) {
	name := c.Param("name")
	c.String(http.StatusOK, fmt.Sprintf("Hello, %s", name))
}

// GinEngine is gin router.
func GinEngine() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.GET("/", ginHelloHandler)
	r.GET("/user/:name", ginUserHandler)

	return r
}
