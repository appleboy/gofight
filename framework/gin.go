package framework

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Binding from JSON
type input struct {
	A int `json:"a" binding:"required"`
	B int `json:"b" binding:"required"`
}

func helloHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "world",
	})
}

func textHandler(c *gin.Context) {
	c.String(http.StatusOK, "Hello World")
}

func queryHandler(c *gin.Context) {
	text := c.Query("text")
	foo := c.Query("foo")

	c.JSON(http.StatusOK, gin.H{
		"hello": text,
		"foo":   foo,
	})
}

func postFormHandler(c *gin.Context) {
	a := c.PostForm("a")
	b := c.PostForm("b")

	c.JSON(http.StatusOK, gin.H{
		"a": a,
		"b": b,
	})
}

func postJSONHandler(c *gin.Context) {
	var json input
	if c.BindJSON(&json) == nil {
		c.JSON(http.StatusOK, gin.H{
			"a": json.A,
			"b": json.B,
		})
	}
}

func putHandler(c *gin.Context) {
	foo := c.PostForm("c")
	bar := c.PostForm("d")

	c.JSON(http.StatusOK, gin.H{
		"c": foo,
		"d": bar,
	})
}

func deleteHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "world",
	})
}

func GinEngine() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.GET("/hello", helloHandler)
	r.GET("/text", textHandler)
	r.GET("/query", queryHandler)

	r.POST("/form", postFormHandler)
	r.POST("/json", postJSONHandler)
	r.PUT("/update", putHandler)
	r.DELETE("/delete", deleteHandler)

	return r
}
