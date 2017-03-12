package example

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ginHelloHandler(c *gin.Context) {
	c.String(http.StatusOK, "Hello World")
}

// GinEngine is gin router.
func GinEngine() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.GET("/", ginHelloHandler)

	return r
}
