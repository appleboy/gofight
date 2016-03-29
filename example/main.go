package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func helloHandler(c *gin.Context) {
	c.String(http.StatusOK, "Hello World")
}

func main() {
	r := gin.New()

	r.GET("/", helloHandler)

	r.Run(":8080")
}
