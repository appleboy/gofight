package framework

import (
	"github.com/labstack/echo"
	"log"
	"net/http"
)

type echoContent struct {
	Hello string `json:"hello"`
	Foo   string `json:"foo"`
	A     string `json:"a"`
	B     string `json:"b"`
	C     string `json:"c"`
	D     string `json:"d"`
}

// Binding from JSON
type echoJSONContent struct {
	A int `json:"a" binding:"required"`
	B int `json:"b" binding:"required"`
}

func echoHelloHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, &echoContent{
			Hello: "world",
		})
	}
}

func echoTextHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World")
	}
}

func echoQueryHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		text := c.QueryParam("text")
		foo := c.QueryParam("foo")

		return c.JSON(http.StatusOK, &echoContent{
			Hello: text,
			Foo:   foo,
		})
	}
}

func echoPostFormHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		a := c.FormValue("a")
		b := c.FormValue("b")

		return c.JSON(http.StatusOK, &echoContent{
			A: a,
			B: b,
		})
	}
}

func echoJSONHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		json := new(echoJSONContent)
		err := c.Bind(json)

		if err != nil {
			log.Println(err)
		}

		return c.JSON(http.StatusOK, json)
	}
}

func echoPutHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		foo := c.FormValue("c")
		bar := c.FormValue("d")

		return c.JSON(http.StatusOK, &echoContent{
			C: foo,
			D: bar,
		})
	}
}

func echoDeleteHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, &echoContent{
			Hello: "world",
		})
	}
}

// EchoEngine is echo router.
func EchoEngine() *echo.Echo {
	e := echo.New()

	e.Get("/hello", echoHelloHandler())
	e.Get("/text", echoTextHandler())
	e.Get("/query", echoQueryHandler())

	e.Post("/form", echoPostFormHandler())
	e.Post("/json", echoJSONHandler())
	e.Put("/update", echoPutHandler())
	e.Delete("/delete", echoDeleteHandler())

	e.Patch("/patch", echoHelloHandler())
	e.Options("/options", echoHelloHandler())
	e.Head("/head", echoHelloHandler())

	return e
}
