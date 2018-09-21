package framework

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

type ginUserContent struct {
	// Username user name
	Username string `json:"account"`
	// Password account password
	Password string `json:"password"`
}

func ginUserHandler(c *gin.Context) {
	var json ginUserContent
	if c.ShouldBind(&json) == nil {
		c.JSON(http.StatusOK, gin.H{
			"username": json.Username,
			"password": json.Password,
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

func gintTimeoutHandler(c *gin.Context) {
	time.Sleep(10 * time.Second)
	c.JSON(http.StatusOK, gin.H{
		"hello": "world",
	})
}

func gintFileUploadHandler(c *gin.Context) {
	file, err := c.FormFile("test")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	foo := c.PostForm("foo")
	bar := c.PostForm("bar")
	c.JSON(http.StatusOK, gin.H{
		"hello":    "world",
		"filename": file.Filename,
		"foo":      foo,
		"bar":      bar,
	})
}

// GinEngine is gin router.
func GinEngine() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.GET("/hello", ginHelloHandler)
	r.GET("/text", ginTextHandler)
	r.GET("/query", ginQueryHandler)
	r.GET("/timeout", gintTimeoutHandler)

	r.POST("/form", ginPostFormHandler)
	r.POST("/json", ginJSONHandler)
	r.POST("/user", ginUserHandler)
	r.PUT("/update", ginPutHandler)
	r.DELETE("/delete", ginDeleteHandler)

	r.PATCH("/patch", ginJSONHandler)
	r.HEAD("/head", ginJSONHandler)
	r.OPTIONS("/options", ginJSONHandler)
	r.POST("/upload", gintFileUploadHandler)

	return r
}
