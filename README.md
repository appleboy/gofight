# mocha

[![Build Status](https://travis-ci.org/appleboy/mocha.svg?branch=master)](https://travis-ci.org/appleboy/mocha) [![Coverage Status](https://coveralls.io/repos/github/appleboy/mocha/badge.svg?branch=master)](https://coveralls.io/github/appleboy/mocha?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/appleboy/mocha)](https://goreportcard.com/report/github.com/appleboy/mocha) [![codebeat badge](https://codebeat.co/badges/4d8b58ae-67ec-469e-bde6-be3dd336b30d)](https://codebeat.co/projects/github-com-appleboy-mocha)

API Handler Testing for Golang framework.

## Support Framework

* [x] [Gin](https://github.com/gin-gonic/gin)
* [x] [Echo](https://github.com/labstack/echo)

## Installation

```
$ go get -u github.com/appleboy/mocha
```

## Usage

### Gin Framework

[gin.go](example/gin.go)

```go
package example

import (
  "github.com/gin-gonic/gin"
  "net/http"
)

func helloHandler(c *gin.Context) {
  c.String(http.StatusOK, "Hello World")
}

func GinEngine() *gin.Engine {
  gin.SetMode(gin.TestMode)
  r := gin.New()

  r.GET("/", helloHandler)

  return r
}

```

[gin_test.go](example/gin.go)

```go
package example

import (
  "github.com/appleboy/mocha"
  "github.com/stretchr/testify/assert"
  "net/http"
  "testing"
)

func TestGinHelloWorld(t *testing.T) {
  r := mocha.New()

  r.GET("/").
    SetDebug(true).
    RunGin(GinEngine(), func(r mocha.HttpResponse, rq mocha.HttpRequest) {
      assert.Equal(t, "Hello World", r.Body.String())
      assert.Equal(t, http.StatusOK, r.Code)
    })
}
```

### Echo Framework

[echo.go](example/echo.go)

```go
package example

import (
  "github.com/labstack/echo"
  "net/http"
)

// Handler
func hello() echo.HandlerFunc {
  return func(c echo.Context) error {
    return c.String(http.StatusOK, "Hello World")
  }
}

func EchoEngine() *echo.Echo {
  // Echo instance
  e := echo.New()

  // Routes
  e.Get("/", hello())

  return e
}

```

[echo_test.go](example/echo.go)

```go
package example

import (
  "github.com/appleboy/mocha"
  "github.com/stretchr/testify/assert"
  "net/http"
  "testing"
)

func TestEchoHelloWorld(t *testing.T) {
  r := mocha.New()

  r.GET("/").
    SetDebug(true).
    RunEcho(EchoEngine(), func(r mocha.EchoHttpResponse, rq mocha.EchoHttpRequest) {
      assert.Equal(t, "Hello World", r.Body.String())
      assert.Equal(t, http.StatusOK, r.Status())
    })
}

```

