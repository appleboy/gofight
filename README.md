# Gofight

[![GoDoc](https://godoc.org/github.com/appleboy/gofight?status.svg)](https://godoc.org/github.com/appleboy/gofight) [![Build Status](https://travis-ci.org/appleboy/gofight.svg?branch=master)](https://travis-ci.org/appleboy/gofight) [![Go Report Card](https://goreportcard.com/badge/github.com/appleboy/gofight)](https://goreportcard.com/report/github.com/appleboy/gofight) [![codebeat badge](https://codebeat.co/badges/4d8b58ae-67ec-469e-bde6-be3dd336b30d)](https://codebeat.co/projects/github-com-appleboy-gofight) [![codecov](https://codecov.io/gh/appleboy/gofight/branch/master/graph/badge.svg)](https://codecov.io/gh/appleboy/gofight) [![Sourcegraph](https://sourcegraph.com/github.com/appleboy/gofight/-/badge.svg)](https://sourcegraph.com/github.com/appleboy/gofight?badge)

API Handler Testing for Golang Web framework.

## Support Framework

* [x] [Http Handler](https://golang.org/pkg/net/http/) Golang package http provides HTTP client and server implementations.
* [x] [Gin](https://github.com/gin-gonic/gin)
* [x] [Echo](https://github.com/labstack/echo) support [v3.0.0](https://github.com/labstack/echo/releases/tag/v3.0.0) up
* [x] [Mux](https://github.com/gorilla/mux)
* [x] [HttpRouter](https://github.com/julienschmidt/httprouter)
* [x] [Pat](https://github.com/gorilla/pat)

## Usage

Download this package.

```bash
$ go get github.com/appleboy/gofight
```

To import this package, add the following line to your code:

```go
import "github.com/appleboy/gofight"
```

## Usage

The following is basic testing example.

Main Program:

```go
package main

import (
  "io"
  "net/http"
)

func BasicHelloHandler(w http.ResponseWriter, r *http.Request) {
  io.WriteString(w, "Hello World")
}

func BasicEngine() http.Handler {
  mux := http.NewServeMux()
  mux.HandleFunc("/", BasicHelloHandler)

  return mux
}
```

Testing:

```go
package main

import (
  "github.com/appleboy/gofight"
  "github.com/stretchr/testify/assert"
  "net/http"
  "testing"
)

func TestBasicHelloWorld(t *testing.T) {
  r := gofight.New()

  r.GET("/").
    // turn on the debug mode.
    SetDebug(true).
    Run(BasicEngine(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {

      assert.Equal(t, "Hello World", r.Body.String())
      assert.Equal(t, http.StatusOK, r.Code)
    })
}
```

### Set Header

You can add custom header via `SetHeader` func.

```go
func TestBasicHelloWorld(t *testing.T) {
  r := gofight.New()
  version := "0.0.1"

  r.GET("/").
    // turn on the debug mode.
    SetDebug(true).
    SetHeader(gofight.H{
      "X-Version": version,
    }).
    Run(BasicEngine(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {

      assert.Equal(t, version, rq.Header.Get("X-Version"))
      assert.Equal(t, "Hello World", r.Body.String())
      assert.Equal(t, http.StatusOK, r.Code)
    })
}
```

### POST FORM Data

Using `SetForm` to generate form data.

```go
func TestPostFormData(t *testing.T) {
  r := gofight.New()

  r.POST("/form").
    SetForm(gofight.H{
      "a": "1",
      "b": "2",
    }).
    Run(BasicEngine(), func(r HTTPResponse, rq HTTPRequest) {
      data := []byte(r.Body.String())

      a, _ := jsonparser.GetString(data, "a")
      b, _ := jsonparser.GetString(data, "b")

      assert.Equal(t, "1", a)
      assert.Equal(t, "2", b)
      assert.Equal(t, http.StatusOK, r.Code)
    })
}
```

### POST JSON Data

Using `SetJSON` to generate JSON data.

```go
func TestPostJSONData(t *testing.T) {
  r := gofight.New()

  r.POST("/json").
    SetJSON(gofight.D{
      "a": 1,
      "b": 2,
    }).
    Run(BasicEngine, func(r HTTPResponse, rq HTTPRequest) {
      data := []byte(r.Body.String())

      a, _ := jsonparser.GetInt(data, "a")
      b, _ := jsonparser.GetInt(data, "b")

      assert.Equal(t, 1, int(a))
      assert.Equal(t, 2, int(b))
      assert.Equal(t, http.StatusOK, r.Code)
      assert.Equal(t, "application/json; charset=utf-8", r.HeaderMap.Get("Content-Type"))
    })
}
```

### POST RAW Data

Using `SetBody` to generate raw data.

```go
func TestPostRawData(t *testing.T) {
  r := gofight.New()

  r.POST("/raw").
    SetBody("a=1&b=1").
    Run(BasicEngine, func(r HTTPResponse, rq HTTPRequest) {
      data := []byte(r.Body.String())

      a, _ := jsonparser.GetString(data, "a")
      b, _ := jsonparser.GetString(data, "b")

      assert.Equal(t, "1", a)
      assert.Equal(t, "2", b)
      assert.Equal(t, http.StatusOK, r.Code)
    })
}
```

### Set Query String

Using `SetQuery` to generate raw data.

```go
func TestQueryString(t *testing.T) {
  r := gofight.New()

  r.GET("/hello").
    SetQuery(gofight.H{
      "a": "1",
      "b": "2",
    }).
    Run(BasicEngine, func(r HTTPResponse, rq HTTPRequest) {
      assert.Equal(t, http.StatusOK, r.Code)
    })
}
```

or append exist query parameter.

```go
func TestQueryString(t *testing.T) {
  r := gofight.New()

  r.GET("/hello?foo=bar").
    SetQuery(gofight.H{
      "a": "1",
      "b": "2",
    }).
    Run(BasicEngine, func(r HTTPResponse, rq HTTPRequest) {
      assert.Equal(t, http.StatusOK, r.Code)
    })
}
```

### Set Cookie String

Using `SetCookie` to generate raw data.

```go
func TestQueryString(t *testing.T) {
  r := gofight.New()

  r.GET("/hello").
    SetCookie(gofight.H{
      "foo": "bar",
    }).
    Run(BasicEngine, func(r HTTPResponse, rq HTTPRequest) {
      assert.Equal(t, http.StatusOK, r.Code)
      assert.Equal(t, "foo=bar", rq.Header.Get("cookie"))
    })
}
```

### Set JSON Struct

```go
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
```

## Example

* Basic HTTP Router: [basic.go](example/basic.go), [basic_test.go](example/basic_test.go)
* Gin Framework: [gin.go](example/gin.go), [gin_test.go](example/gin_test.go)
* Echo Framework: [echo.go](example/echo.go), [echo_test.go](example/echo_test.go)
* Mux Framework: [mux.go](example/mux.go), [mux_test.go](example/mux_test.go)
* HttpRouter Framework: [httprouter.go](example/httprouter.go), [httprouter_test.go](example/httprouter_test.go)

## Testing

```
$ make
```

## License

Copyright 2016 Bo-Yi Wu [@appleboy](https://twitter.com/appleboy).

Licensed under the MIT License.
