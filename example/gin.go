package example

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GinHelloHandler(c *gin.Context) {
	c.String(http.StatusOK, "Hello World")
}

func GinEngine() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.GET("/", GinHelloHandler)

	return r
}
