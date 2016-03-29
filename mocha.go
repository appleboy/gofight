package ginMocha

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
)

// request handling func type to replace gin.HandlerFunc
type RequestFunc func(*gin.Context)

// response handling func type
type ResponseFunc func(*httptest.ResponseRecorder)

type RequestConfig struct {
	Method      string
	Path        string
	Body        string
	Headers     map[string]string
	Middlewares []gin.HandlerFunc
	Handler     RequestFunc
	Callback    ResponseFunc
	Debug       bool
}

func (rc *RequestConfig) SetDebug(enable bool) *RequestConfig {
	rc.Debug = enable

	return rc
}

func (rc *RequestConfig) Run() {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	if rc.Method == "" {
		rc.Method = "GET"
	}

	if rc.Path == "" {
		rc.Path = "/"
	}

	if rc.Middlewares != nil && len(rc.Middlewares) > 0 {
		for _, mw := range rc.Middlewares {
			r.Use(mw)
		}
	}

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

	if len(rc.Headers) > 0 {
		for k, v := range rc.Headers {
			req.Header.Set(k, v)
		}
	} else if rc.Method == "POST" || rc.Method == "PUT" {
		if strings.HasPrefix(rc.Body, "{") {
			req.Header.Set("Content-Type", "application/json")
		} else {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
	}

	r.Handle(rc.Method, rc.Path, func(c *gin.Context) {
		//change argument if necessary here
		rc.Handler(c)
	})

	if rc.Debug {
		log.Printf("Request Method: %s", rc.Method)
		log.Printf("Request Path: %s", rc.Path)
		log.Printf("Request Body: %s", rc.Body)
		log.Printf("Request Headers: %s", rc.Headers)
		log.Printf("Request Header: %s", req.Header)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if rc.Callback != nil {
		rc.Callback(w)
	}
}
