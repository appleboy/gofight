package framework

import (
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"
)

// Binding from JSON
type ginJSONContent struct {
	A int `json:"a" binding:"required"`
	B int `json:"b" binding:"required"`
}

func ginHelloHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "world",
	})
}

func ginTextHandler(c *gin.Context) {
	c.String(http.StatusOK, "Hello World")
}

func ginQueryHandler(c *gin.Context) {
	text := c.Query("text")
	foo := c.Query("foo")

	c.JSON(http.StatusOK, gin.H{
		"hello": text,
		"foo":   foo,
	})
}

func ginPostFormHandler(c *gin.Context) {
	a := c.PostForm("a")
	b := c.PostForm("b")

	c.JSON(http.StatusOK, gin.H{
		"a": a,
		"b": b,
	})
}

func ginJSONHandler(c *gin.Context) {
	var json ginJSONContent
	if c.BindJSON(&json) == nil {
		c.JSON(http.StatusOK, gin.H{
			"a": json.A,
			"b": json.B,
		})
	}
}

func ginPutHandler(c *gin.Context) {
	foo := c.PostForm("c")
	bar := c.PostForm("d")

	c.JSON(http.StatusOK, gin.H{
		"c": foo,
		"d": bar,
	})
}

func ginDeleteHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "world",
	})
}

// GinEngine is gin router.
func GinEngine() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.GET("/hello", ginHelloHandler)
	r.GET("/text", ginTextHandler)
	r.GET("/query", ginQueryHandler)

	r.POST("/form", ginPostFormHandler)
	r.POST("/json", ginJSONHandler)
	r.PUT("/update", ginPutHandler)
	r.DELETE("/delete", ginDeleteHandler)

	r.PATCH("/patch", ginJSONHandler)
	r.HEAD("/head", ginJSONHandler)
	r.OPTIONS("/options", ginJSONHandler)

	return r
}
