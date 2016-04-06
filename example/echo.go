package example

import (
	"github.com/labstack/echo"
	"net/http"
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
	e.Get("/", echoHelloHandler())

	return e
}
