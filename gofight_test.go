package gofight

import (
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"testing"
	"time"

	"github.com/appleboy/gofight/v2/framework"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

var goVersion = runtime.Version()

func TestHttpURL(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/example", func(c *gin.Context) { c.String(http.StatusOK, "it worked") })

	go func() {
		assert.NoError(t, router.Run())
	}()
	// have to wait for the goroutine to start and run the server
	// otherwise the main thread will complete
	time.Sleep(5 * time.Millisecond)

	TestRequest(t, "http://localhost:8080/example")
}

func TestHttpsURL(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/example", func(c *gin.Context) { c.String(http.StatusOK, "it worked") })

	go func() {
		assert.NoError(t, router.RunTLS(":8088", "certificate/localhost.cert", "certificate/localhost.key"))
	}()
	// have to wait for the goroutine to start and run the server
	// otherwise the main thread will complete
	time.Sleep(5 * time.Millisecond)

	TestRequest(t, "https://localhost:8088/example")
}

func TestGinHelloWorld(t *testing.T) {
	r := New()

	r.GET("/hello").
		SetDebug(true).
		Run(framework.GinEngine(), func(r HTTPResponse, rq HTTPRequest) {
			data := []byte(r.Body.String())
			value := gjson.GetBytes(data, "hello")
			assert.Equal(t, "world", value.String())
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestGinHeader(t *testing.T) {
	r := New()

	r.GET("/text").
		SetHeader(H{
			"Content-Type": "text/plain",
			"Go-Version":   goVersion,
		}).
		Run(framework.GinEngine(), func(r HTTPResponse, rq HTTPRequest) {
			assert.Equal(t, goVersion, rq.Header.Get("Go-Version"))
			assert.Equal(t, "Gofight-client/"+Version, rq.Header.Get(UserAgent))
			assert.Equal(t, "text/plain", rq.Header.Get(ContentType))
			assert.Equal(t, "Hello World", r.Body.String())
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestGinCookie(t *testing.T) {
	r := New()

	r.GET("/text").
		SetCookie(H{
			"foo": "bar",
		}).
		Run(framework.GinEngine(), func(r HTTPResponse, rq HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
			assert.Equal(t, "foo=bar", rq.Header.Get("cookie"))
		})
}

func TestGinQuery(t *testing.T) {
	r := New()

	r.GET("/query?text=world&foo=bar").
		Run(framework.GinEngine(), func(r HTTPResponse, rq HTTPRequest) {
			data := []byte(r.Body.String())
			value := gjson.GetBytes(data, "hello")
			foo := gjson.GetBytes(data, "foo")

			assert.Equal(t, "world", value.String())
			assert.Equal(t, "bar", foo.String())
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestGinPostFormData(t *testing.T) {
	r := New()

	r.POST("/form").
		SetForm(H{
			"a": "1",
			"b": "2",
		}).
		Run(framework.GinEngine(), func(r HTTPResponse, rq HTTPRequest) {
			data := []byte(r.Body.String())
			a := gjson.GetBytes(data, "a")
			b := gjson.GetBytes(data, "b")

			assert.Equal(t, "1", a.String())
			assert.Equal(t, "2", b.String())
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestGinPostJSONData(t *testing.T) {
	r := New()

	r.POST("/json").
		SetJSON(D{
			"a": 1,
			"b": 2,
		}).
		Run(framework.GinEngine(), func(r HTTPResponse, rq HTTPRequest) {
			data := []byte(r.Body.String())
			a := gjson.GetBytes(data, "a")
			b := gjson.GetBytes(data, "b")

			assert.Equal(t, int64(1), a.Int())
			assert.Equal(t, int64(2), b.Int())
			assert.Equal(t, http.StatusOK, r.Code)
			assert.Equal(t, "application/json; charset=utf-8", r.HeaderMap.Get("Content-Type"))
		})
}

func TestGinPut(t *testing.T) {
	r := New()

	r.PUT("/update").
		SetBody("c=1&d=2").
		Run(framework.GinEngine(), func(r HTTPResponse, rq HTTPRequest) {
			data := []byte(r.Body.String())
			c := gjson.GetBytes(data, "c")
			d := gjson.GetBytes(data, "d")

			assert.Equal(t, "1", c.String())
			assert.Equal(t, "2", d.String())
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestGinDelete(t *testing.T) {
	r := New()

	r.DELETE("/delete").
		Run(framework.GinEngine(), func(r HTTPResponse, rq HTTPRequest) {
			data := []byte(r.Body.String())
			hello := gjson.GetBytes(data, "hello")

			assert.Equal(t, "world", hello.String())
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestGinPatch(t *testing.T) {
	r := New()

	r.PATCH("/patch").
		SetJSON(D{
			"a": 1,
			"b": 2,
		}).
		Run(framework.GinEngine(), func(r HTTPResponse, rq HTTPRequest) {
			data := []byte(r.Body.String())
			a := gjson.GetBytes(data, "a")
			b := gjson.GetBytes(data, "b")

			assert.Equal(t, int64(1), a.Int())
			assert.Equal(t, int64(2), b.Int())
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestGinHead(t *testing.T) {
	r := New()

	r.HEAD("/head").
		SetJSON(D{
			"a": 1,
			"b": 2,
		}).
		Run(framework.GinEngine(), func(r HTTPResponse, rq HTTPRequest) {
			data := []byte(r.Body.String())
			a := gjson.GetBytes(data, "a")
			b := gjson.GetBytes(data, "b")

			assert.Equal(t, int64(1), a.Int())
			assert.Equal(t, int64(2), b.Int())
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestGinOptions(t *testing.T) {
	r := New()

	r.OPTIONS("/options").
		SetJSON(D{
			"a": 1,
			"b": 2,
		}).
		Run(framework.GinEngine(), func(r HTTPResponse, rq HTTPRequest) {
			data := []byte(r.Body.String())
			a := gjson.GetBytes(data, "a")
			b := gjson.GetBytes(data, "b")

			assert.Equal(t, int64(1), a.Int())
			assert.Equal(t, int64(2), b.Int())
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestEchoHelloWorld(t *testing.T) {
	r := New()

	r.GET("/hello").
		SetDebug(true).
		Run(framework.EchoEngine(), func(r HTTPResponse, rq HTTPRequest) {
			data := []byte(r.Body.String())
			hello := gjson.GetBytes(data, "hello")

			assert.Equal(t, "world", hello.String())
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestEchoHeader(t *testing.T) {
	r := New()

	r.GET("/text").
		SetHeader(H{
			"Content-Type": "text/plain",
			"Go-Version":   goVersion,
		}).
		Run(framework.EchoEngine(), func(r HTTPResponse, rq HTTPRequest) {

			assert.Equal(t, goVersion, rq.Header.Get("Go-Version"))
			assert.Equal(t, r.Body.String(), "Hello World")
			assert.Equal(t, r.Code, http.StatusOK)
		})
}

func TestEchoQuery(t *testing.T) {
	r := New()

	r.GET("/query?text=world&foo=bar").
		Run(framework.EchoEngine(), func(r HTTPResponse, rq HTTPRequest) {
			data := []byte(r.Body.String())
			value := gjson.GetBytes(data, "hello")
			foo := gjson.GetBytes(data, "foo")

			assert.Equal(t, "world", value.String())
			assert.Equal(t, "bar", foo.String())
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestEchoPostFormData(t *testing.T) {
	r := New()

	r.POST("/form").
		SetBody("a=1&b=2").
		Run(framework.EchoEngine(), func(r HTTPResponse, rq HTTPRequest) {
			data := []byte(r.Body.String())
			a := gjson.GetBytes(data, "a")
			b := gjson.GetBytes(data, "b")

			assert.Equal(t, "1", a.String())
			assert.Equal(t, "2", b.String())
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestEchoPostJSONData(t *testing.T) {
	r := New()

	r.POST("/json").
		SetJSON(D{
			"a": 1,
			"b": 2,
		}).
		Run(framework.EchoEngine(), func(r HTTPResponse, rq HTTPRequest) {
			data := []byte(r.Body.String())
			a := gjson.GetBytes(data, "a")
			b := gjson.GetBytes(data, "b")

			assert.Equal(t, int64(1), a.Int())
			assert.Equal(t, int64(2), b.Int())
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestEchoPut(t *testing.T) {
	r := New()

	r.PUT("/update").
		SetBody("c=1&d=2").
		Run(framework.EchoEngine(), func(r HTTPResponse, rq HTTPRequest) {
			data := []byte(r.Body.String())
			c := gjson.GetBytes(data, "c")
			d := gjson.GetBytes(data, "d")

			assert.Equal(t, "1", c.String())
			assert.Equal(t, "2", d.String())
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestEchoDelete(t *testing.T) {
	r := New()

	r.DELETE("/delete").
		Run(framework.EchoEngine(), func(r HTTPResponse, rq HTTPRequest) {
			data := []byte(r.Body.String())
			hello := gjson.GetBytes(data, "hello")

			assert.Equal(t, "world", hello.String())
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestEchoPatch(t *testing.T) {
	r := New()

	r.PATCH("/patch").
		Run(framework.EchoEngine(), func(r HTTPResponse, rq HTTPRequest) {
			data := []byte(r.Body.String())
			hello := gjson.GetBytes(data, "hello")

			assert.Equal(t, "world", hello.String())
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestEchoHead(t *testing.T) {
	r := New()

	r.HEAD("/head").
		Run(framework.EchoEngine(), func(r HTTPResponse, rq HTTPRequest) {
			data := []byte(r.Body.String())
			hello := gjson.GetBytes(data, "hello")

			assert.Equal(t, "world", hello.String())
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestEchoOptions(t *testing.T) {
	r := New()

	r.OPTIONS("/options").
		Run(framework.EchoEngine(), func(r HTTPResponse, rq HTTPRequest) {
			data := []byte(r.Body.String())
			hello := gjson.GetBytes(data, "hello")

			assert.Equal(t, "world", hello.String())
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestSetQueryString(t *testing.T) {
	r := New()

	r.GET("/hello").
		SetQuery(H{
			"a": "1",
			"b": "2",
		})

	assert.Equal(t, "/hello?a=1&b=2", r.Path)

	r.GET("/hello?foo=bar").
		SetQuery(H{
			"a": "1",
			"b": "2",
		})

	assert.Equal(t, "/hello?foo=bar&a=1&b=2", r.Path)
}

type User struct {
	// Username user name
	Username string `json:"account"`
	// Password account password
	Password string `json:"password"`
}

func TestSetJSONInterface(t *testing.T) {
	r := New()

	r.POST("/user").
		SetJSONInterface(User{
			Username: "foo",
			Password: "bar",
		}).
		Run(framework.GinEngine(), func(r HTTPResponse, rq HTTPRequest) {
			data := []byte(r.Body.String())

			username := gjson.GetBytes(data, "username")
			password := gjson.GetBytes(data, "password")

			assert.Equal(t, "foo", username.String())
			assert.Equal(t, "bar", password.String())
			assert.Equal(t, http.StatusOK, r.Code)
			assert.Equal(t, "application/json; charset=utf-8", r.HeaderMap.Get("Content-Type"))
		})
}

func TestUploadFile(t *testing.T) {
	r := New()

	r.POST("/upload").
		SetDebug(true).
		SetFileFromPath([]UploadFile{
			{
				Path: "./testdata/hello.txt",
				Name: "hello",
			},
			{
				Path: "./testdata/world.txt",
				Name: "world",
			},
		}, H{
			"foo": "bar",
			"bar": "foo",
		}).
		Run(framework.GinEngine(), func(r HTTPResponse, rq HTTPRequest) {
			data := []byte(r.Body.String())

			hello := gjson.GetBytes(data, "hello")
			world := gjson.GetBytes(data, "world")
			foo := gjson.GetBytes(data, "foo")
			bar := gjson.GetBytes(data, "bar")
			ip := gjson.GetBytes(data, "ip")
			helloSize := gjson.GetBytes(data, "helloSize")
			worldSize := gjson.GetBytes(data, "worldSize")

			assert.Equal(t, "world\n", helloSize.String())
			assert.Equal(t, "hello\n", worldSize.String())
			assert.Equal(t, "hello.txt", hello.String())
			assert.Equal(t, "world.txt", world.String())
			assert.Equal(t, "bar", foo.String())
			assert.Equal(t, "foo", bar.String())
			assert.Equal(t, "", ip.String())
			assert.Equal(t, http.StatusOK, r.Code)
			assert.Equal(t, "application/json; charset=utf-8", r.HeaderMap.Get("Content-Type"))
		})
}

func TestUploadFileByContent(t *testing.T) {
	r := New()

	helloContent, err := ioutil.ReadFile("./testdata/hello.txt")
	if err != nil {
		log.Fatal(err)
	}

	worldContent, err := ioutil.ReadFile("./testdata/world.txt")
	if err != nil {
		log.Fatal(err)
	}

	r.POST("/upload").
		SetDebug(true).
		SetFileFromPath([]UploadFile{
			{
				Path:    "hello.txt",
				Name:    "hello",
				Content: helloContent,
			},
			{
				Path:    "world.txt",
				Name:    "world",
				Content: worldContent,
			},
		}, H{
			"foo": "bar",
			"bar": "foo",
		}).
		Run(framework.GinEngine(), func(r HTTPResponse, rq HTTPRequest) {
			data := []byte(r.Body.String())

			hello := gjson.GetBytes(data, "hello")
			world := gjson.GetBytes(data, "world")
			foo := gjson.GetBytes(data, "foo")
			bar := gjson.GetBytes(data, "bar")
			ip := gjson.GetBytes(data, "ip")
			helloSize := gjson.GetBytes(data, "helloSize")
			worldSize := gjson.GetBytes(data, "worldSize")

			assert.Equal(t, "world\n", helloSize.String())
			assert.Equal(t, "hello\n", worldSize.String())
			assert.Equal(t, "hello.txt", hello.String())
			assert.Equal(t, "world.txt", world.String())
			assert.Equal(t, "bar", foo.String())
			assert.Equal(t, "foo", bar.String())
			assert.Equal(t, "", ip.String())
			assert.Equal(t, http.StatusOK, r.Code)
			assert.Equal(t, "application/json; charset=utf-8", r.HeaderMap.Get("Content-Type"))
		})
}
