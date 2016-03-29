package mocha

import (
	"github.com/buger/jsonparser"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Binding from JSON
type Login struct {
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

func postHandler(c *gin.Context) {
	a := c.PostForm("a")
	b := c.PostForm("b")

	c.JSON(http.StatusOK, gin.H{
		"a": a,
		"b": b,
	})
}

func testPostJSONHandler(c *gin.Context) {
	var json Login
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

func TestHelloWorld(t *testing.T) {
	r := &RequestConfig{
		Handler: helloHandler,
		Callback: func(r *httptest.ResponseRecorder) {
			data := []byte(r.Body.String())

			value, _ := jsonparser.GetString(data, "hello")

			assert.Equal(t, value, "world")
			assert.Equal(t, r.Code, http.StatusOK)
		},
	}

	r.SetDebug(true).Run()
}

func TestHeader(t *testing.T) {
	r := &RequestConfig{
		Handler: textHandler,
		Headers: map[string]string{
			"Content-Type": "text/plain",
		},
		Callback: func(r *httptest.ResponseRecorder) {

			assert.Equal(t, r.Body.String(), "Hello World")
			assert.Equal(t, r.Code, http.StatusOK)
		},
	}

	r.Run()
}

func TestQuery(t *testing.T) {
	r := &RequestConfig{
		Path:    "/hello?text=world&foo=bar",
		Handler: queryHandler,
		Callback: func(r *httptest.ResponseRecorder) {
			data := []byte(r.Body.String())

			hello, _ := jsonparser.GetString(data, "hello")
			foo, _ := jsonparser.GetString(data, "foo")

			assert.Equal(t, hello, "world")
			assert.Equal(t, foo, "bar")
			assert.Equal(t, r.Code, http.StatusOK)
		},
	}

	r.Run()
}

func TestPost(t *testing.T) {
	r := &RequestConfig{
		Method:  "POST",
		Body:    "a=1&b=2",
		Handler: postHandler,
		Callback: func(r *httptest.ResponseRecorder) {
			data := []byte(r.Body.String())

			a, _ := jsonparser.GetString(data, "a")
			b, _ := jsonparser.GetString(data, "b")

			assert.Equal(t, a, "1")
			assert.Equal(t, b, "2")
			assert.Equal(t, r.Code, http.StatusOK)
		},
	}

	r.Run()
}

func TestPostJsonData(t *testing.T) {
	r := &RequestConfig{
		Method:  "POST",
		Body:    `{"a":1,"b":2}`,
		Handler: testPostJSONHandler,
		Callback: func(r *httptest.ResponseRecorder) {
			data := []byte(r.Body.String())

			a, _ := jsonparser.GetInt(data, "a")
			b, _ := jsonparser.GetInt(data, "b")

			assert.Equal(t, int(a), 1)
			assert.Equal(t, int(b), 2)
			assert.Equal(t, r.Code, http.StatusOK)
		},
	}

	r.SetDebug(true).Run()
}

func TestPut(t *testing.T) {
	r := &RequestConfig{
		Method:  "PUT",
		Body:    "c=1&d=2",
		Handler: putHandler,
		Callback: func(r *httptest.ResponseRecorder) {
			data := []byte(r.Body.String())

			c, _ := jsonparser.GetString(data, "c")
			d, _ := jsonparser.GetString(data, "d")

			assert.Equal(t, c, "1")
			assert.Equal(t, d, "2")
			assert.Equal(t, r.Code, http.StatusOK)
		},
	}

	r.Run()
}
