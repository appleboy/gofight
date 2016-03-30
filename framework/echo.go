package framework

import (
	"github.com/labstack/echo"
	"net/http"
)

type echoContent struct {
	Hello string `json:"hello"`
	Foo   string `json:"foo"`
	A     string `json:"a"`
	B     string `json:"b"`
	C     string `json:"c"`
	D     string `json:"d"`
	E     string `json:"e"`
	F     string `json:"f"`
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

func echoPostJSONHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		var json echoContent
		err := c.Bind(json)

		if err != nil {

		}

		return c.JSON(http.StatusOK, json)
	}
}

func echoPutHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		foo := c.FormValue("c")
		bar := c.FormValue("d")

		return c.JSON(http.StatusOK, &echoContent{
			A: foo,
			B: bar,
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

func EchoEngine() *echo.Echo {
	e := echo.New()

	e.Get("/hello", echoHelloHandler())
	e.Get("/text", echoTextHandler())
	e.Get("/query", echoQueryHandler())

	e.Post("/form", echoPostFormHandler())
	e.Post("/json", echoPostJSONHandler())
	e.Put("/update", echoPutHandler())
	e.Delete("/delete", echoDeleteHandler())

	return e
}
