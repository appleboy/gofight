package mocha

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine"
	"github.com/labstack/echo/test"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
)

// Media types
const (
	ContentType     = "Content-Type"
	ApplicationJSON = "application/json"
	ApplicationForm = "application/x-www-form-urlencoded"
)

// Gin http request and response
type HttpResponse *httptest.ResponseRecorder
type HttpRequest *http.Request

// Echo http request and response
type EchoHttpResponse *test.ResponseRecorder
type EchoHttpRequest engine.Request

// response handling func type
type ResponseFunc func(HttpResponse, HttpRequest)

// echo response handling func type
type EchoResponseFunc func(EchoHttpResponse, EchoHttpRequest)

type RequestConfig struct {
	Method  string
	Path    string
	Body    string
	Headers map[string]string
	Debug   bool
}

func New() *RequestConfig {

	return &RequestConfig{}
}

func (rc *RequestConfig) SetDebug(enable bool) *RequestConfig {
	rc.Debug = enable

	return rc
}

func (rc *RequestConfig) GET(path string) *RequestConfig {
	rc.Path = path
	rc.Method = "GET"

	return rc
}

func (rc *RequestConfig) POST(path string) *RequestConfig {
	rc.Path = path
	rc.Method = "POST"

	return rc
}

func (rc *RequestConfig) PUT(path string) *RequestConfig {
	rc.Path = path
	rc.Method = "PUT"

	return rc
}

func (rc *RequestConfig) DELETE(path string) *RequestConfig {
	rc.Path = path
	rc.Method = "DELETE"

	return rc
}

func (rc *RequestConfig) SetHeader(headers map[string]string) *RequestConfig {
	if len(headers) > 0 {
		rc.Headers = headers
	}

	return rc
}

func (rc *RequestConfig) SetBody(body string) *RequestConfig {
	if len(body) > 0 {
		rc.Body = body
	}

	return rc
}

func (rc *RequestConfig) InitGinTest() (*http.Request, *httptest.ResponseRecorder) {
	qs := ""
	if strings.Contains(rc.Path, "?") {
		ss := strings.Split(rc.Path, "?")
		rc.Path = ss[0]
		qs = ss[1]
	}

	body := bytes.NewBufferString(rc.Body)

	req, _ := http.NewRequest(rc.Method, rc.Path, body)

	if len(qs) > 0 {
		req.URL.RawQuery = qs
	}

	if rc.Method == "POST" || rc.Method == "PUT" {
		if strings.HasPrefix(rc.Body, "{") {
			req.Header.Set(ContentType, ApplicationJSON)
		} else {
			req.Header.Set(ContentType, ApplicationForm)
		}
	}

	if len(rc.Headers) > 0 {
		for k, v := range rc.Headers {
			req.Header.Set(k, v)
		}
	}

	if rc.Debug {
		log.Printf("Request Method: %s", rc.Method)
		log.Printf("Request Path: %s", rc.Path)
		log.Printf("Request Body: %s", rc.Body)
		log.Printf("Request Headers: %s", rc.Headers)
		log.Printf("Request Header: %s", req.Header)
	}

	w := httptest.NewRecorder()

	return req, w
}

func (rc *RequestConfig) RunGin(r *gin.Engine, response ResponseFunc) {

	req, w := rc.InitGinTest()
	r.ServeHTTP(w, req)

	response(w, req)
}

func (rc *RequestConfig) InitEchoTest() (engine.Request, *test.ResponseRecorder) {

	rq := test.NewRequest(rc.Method, rc.Path, strings.NewReader(rc.Body))
	rec := test.NewResponseRecorder()

	if rc.Method == "POST" || rc.Method == "PUT" {
		if strings.HasPrefix(rc.Body, "{") {
			rq.Header().Add(ContentType, ApplicationJSON)
		} else {
			rq.Header().Add(ContentType, ApplicationForm)
		}
	}

	for k, v := range rc.Headers {
		rq.Header().Add(k, v)
	}

	return rq, rec
}

func (rc *RequestConfig) RunEcho(e *echo.Echo, response EchoResponseFunc) {

	rq, rec := rc.InitEchoTest()
	e.ServeHTTP(rq, rec)

	response(rec, rq)
}
