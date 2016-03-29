# gin-mocha

[![Build Status](https://travis-ci.org/appleboy/gin-mocha.svg?branch=master)](https://travis-ci.org/appleboy/gin-mocha) [![Coverage Status](https://coveralls.io/repos/github/appleboy/gin-mocha/badge.svg?branch=master)](https://coveralls.io/github/appleboy/gin-mocha?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/appleboy/gin-mocha)](https://goreportcard.com/report/github.com/appleboy/gin-mocha)

API Handler Testing for Gin framework written in Golang.

## Installation

```
$ go get -u github.com/appleboy/gin-mocha
```

## Usage

main.go

```go
package main

import (
  "github.com/gin-gonic/gin"
  "net/http"
)

func helloHandler(c *gin.Context) {
  c.String(http.StatusOK, "Hello World")
}

func main() {
  r := gin.New()

  r.GET("/", helloHandler)

  r.Run(":8080")
}
```

main_test.go

```go
package main

import (
  "github.com/appleboy/gin-mocha"
  "github.com/stretchr/testify/assert"
  "net/http"
  "net/http/httptest"
  "testing"
)

func TestHelloWorld(t *testing.T) {
  r := &ginMocha.RequestConfig{
    Handler: helloHandler,
    Callback: func(r *httptest.ResponseRecorder) {

      assert.Equal(t, r.Body.String(), "Hello World")
      assert.Equal(t, r.Code, http.StatusOK)
    },
  }

  r.Run()
}
```
