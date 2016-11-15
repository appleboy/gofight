package example

import (
	"net/http"

	"github.com/labstack/echo"
)

func echoHelloHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World")
	}
}

// EchoEngine is echo router.
func EchoEngine() *echo.Echo {
	// Echo instance
	e := echo.New()

	// Routes
	e.GET("/", echoHelloHandler())

	return e
}
