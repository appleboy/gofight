# mocha

[![Build Status](https://travis-ci.org/appleboy/mocha.svg?branch=master)](https://travis-ci.org/appleboy/mocha) [![Coverage Status](https://coveralls.io/repos/github/appleboy/mocha/badge.svg?branch=master)](https://coveralls.io/github/appleboy/mocha?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/appleboy/mocha)](https://goreportcard.com/report/github.com/appleboy/mocha) [![codebeat badge](https://codebeat.co/badges/4d8b58ae-67ec-469e-bde6-be3dd336b30d)](https://codebeat.co/projects/github-com-appleboy-mocha)

API Handler Testing for Golang framework.

## Support Framework

* [x] [Gin](https://github.com/gin-gonic/gin)
* [ ] [Echo](https://github.com/labstack/echo)

## Installation

```
$ go get -u github.com/appleboy/mocha
```

## Usage

For Gin framework: gin.go

```go
package main

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

main_test.go

```go
package main

import (
  "github.com/appleboy/mocha"
  "github.com/stretchr/testify/assert"
  "net/http"
  "net/http/httptest"
  "testing"
)

func TestHelloWorld(t *testing.T) {
  r := mocha.New()

  r.GET("/").
    SetDebug(true).
    RunGinEngine(GinEngine(), func(r *httptest.ResponseRecorder) {
      assert.Equal(t, r.Body.String(), "Hello World")
      assert.Equal(t, r.Code, http.StatusOK)
    })
}
```
