// Package gofight offers simple API http handler testing for Golang framework.
//
// Details about the gofight project are found in github page:
//
//	https://github.com/appleboy/gofight
//
// Installation:
//
//	$ go get -u github.com/appleboy/gofight
//
// Set Header: You can add custom header via SetHeader func.
//
//	SetHeader(gofight.H{
//	  "X-Version": version,
//	})
//
// Set Cookie: You can add custom cookie via SetCookie func.
//
//	SetCookie(gofight.H{
//	  "foo": "bar",
//	})
//
// Set query string: Using SetQuery to generate query string data.
//
//	SetQuery(gofight.H{
//	  "a": "1",
//	  "b": "2",
//	})
//
// POST FORM Data: Using SetForm to generate form data.
//
//	SetForm(gofight.H{
//	  "a": "1",
//	  "b": "2",
//	})
//
// POST JSON Data: Using SetJSON to generate json data.
//
//	SetJSON(gofight.H{
//	  "a": "1",
//	  "b": "2",
//	})
//
// POST RAW Data: Using SetBody to generate raw data.
//
//	SetBody("a=1&b=1")
//
// For more details, see the documentation and example.
package gofight

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// Media types
const (
	Version         = "1.0"
	UserAgent       = "User-Agent"
	ContentType     = "Content-Type"
	ApplicationJSON = "application/json"
	ApplicationForm = "application/x-www-form-urlencoded"
)

// HTTPResponse wraps the httptest.ResponseRecorder to provide additional
// functionality or to simplify the response handling in tests.
type HTTPResponse struct {
	*httptest.ResponseRecorder
}

// HTTPRequest is a wrapper around the standard http.Request.
// It embeds the http.Request struct, allowing you to use all the methods
// and fields of http.Request while also providing the ability to extend
// its functionality with additional methods or fields if needed.
type HTTPRequest struct {
	*http.Request
}

// ResponseFunc is a type alias for a function that takes an HTTPResponse and an HTTPRequest as parameters.
// It is used to define a callback function that can handle or process HTTP responses and requests.
type ResponseFunc func(HTTPResponse, HTTPRequest)

// H is HTTP Header Type
type H map[string]string

// D is HTTP Data Type
type D map[string]interface{}

// RequestConfig provide user input request structure
type RequestConfig struct {
	Method      string
	Path        string
	Body        string
	Headers     H
	Cookies     H
	Debug       bool
	ContentType string
	Context     context.Context
}

// UploadFile for upload file struct
type UploadFile struct {
	Path    string
	Name    string
	Content []byte
}

// New supply initial structure
func New() *RequestConfig {
	return &RequestConfig{
		Context: context.Background(),
	}
}

// SetDebug supply enable debug mode.
func (rc *RequestConfig) SetDebug(enable bool) *RequestConfig {
	rc.Debug = enable

	return rc
}

// SetContext sets the context for the RequestConfig.
// This allows the request to be aware of deadlines, cancellation signals, and other request-scoped values.
// It returns the updated RequestConfig instance.
//
// Parameters:
//
//	ctx - the context to be set for the RequestConfig
//
// Returns:
//
//	*RequestConfig - the updated RequestConfig instance with the new context
func (rc *RequestConfig) SetContext(ctx context.Context) *RequestConfig {
	rc.Context = ctx

	return rc
}

// GET is request method.
func (rc *RequestConfig) GET(path string) *RequestConfig {
	rc.Path = path
	rc.Method = "GET"

	return rc
}

// POST is request method.
func (rc *RequestConfig) POST(path string) *RequestConfig {
	rc.Path = path
	rc.Method = "POST"

	return rc
}

// PUT is request method.
func (rc *RequestConfig) PUT(path string) *RequestConfig {
	rc.Path = path
	rc.Method = "PUT"

	return rc
}

// DELETE is request method.
func (rc *RequestConfig) DELETE(path string) *RequestConfig {
	rc.Path = path
	rc.Method = "DELETE"

	return rc
}

// PATCH is request method.
func (rc *RequestConfig) PATCH(path string) *RequestConfig {
	rc.Path = path
	rc.Method = "PATCH"

	return rc
}

// HEAD is request method.
func (rc *RequestConfig) HEAD(path string) *RequestConfig {
	rc.Path = path
	rc.Method = "HEAD"

	return rc
}

// OPTIONS is request method.
func (rc *RequestConfig) OPTIONS(path string) *RequestConfig {
	rc.Path = path
	rc.Method = "OPTIONS"

	return rc
}

// SetHeader supply http header what you defined.
func (rc *RequestConfig) SetHeader(headers H) *RequestConfig {
	if len(headers) > 0 {
		rc.Headers = headers
	}

	return rc
}

// SetJSON supply JSON body.
func (rc *RequestConfig) SetJSON(body D) *RequestConfig {
	if b, err := json.Marshal(body); err == nil {
		rc.Body = string(b)
	}

	return rc
}

// SetJSONInterface supply JSON body
func (rc *RequestConfig) SetJSONInterface(body interface{}) *RequestConfig {
	if b, err := json.Marshal(body); err == nil {
		rc.Body = string(b)
	}

	return rc
}

// SetForm sets the form data for the request configuration.
// It takes a map of string keys and values, converts it to url.Values,
// and encodes it as a URL-encoded form string, which is then assigned to the Body field.
//
// Parameters:
//
//	body (H): A map containing the form data to be set.
//
// Returns:
//
//	*RequestConfig: The updated request configuration.
func (rc *RequestConfig) SetForm(body H) *RequestConfig {
	f := make(url.Values)

	for k, v := range body {
		f.Set(k, v)
	}

	rc.Body = f.Encode()

	return rc
}

// SetFileFromPath upload new file.
func (rc *RequestConfig) SetFileFromPath(uploads []UploadFile, params ...H) *RequestConfig {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	for _, f := range uploads {
		reader := bytes.NewReader(f.Content)
		if reader.Size() == 0 {
			file, err := os.Open(f.Path)
			if err != nil {
				return rc
			}

			defer file.Close()
			part, err := writer.CreateFormFile(f.Name, filepath.Base(f.Path))
			if err != nil {
				return rc
			}
			_, err = io.Copy(part, file)
			if err != nil {
				return rc
			}
		} else {
			part, err := writer.CreateFormFile(f.Name, filepath.Base(f.Path))
			if err != nil {
				return rc
			}
			_, err = reader.WriteTo(part)
			if err != nil {
				return rc
			}
		}
	}

	if len(params) > 0 {
		for key, val := range params[0] {
			_ = writer.WriteField(key, val)
		}
	}

	err := writer.Close()
	if err != nil {
		return rc
	}

	rc.ContentType = writer.FormDataContentType()
	rc.Body = body.String()

	return rc
}

// SetPath supply new request path to deal with path variable request
// ex. /reqpath/:book/:apple , usage: r.POST("/reqpath/").SetPath("book1/apple2")...
func (rc *RequestConfig) SetPath(str string) *RequestConfig {
	rc.Path += str
	return rc
}

// SetQueryD supply query string, support query using string array input.
// ex. /reqpath/?Ids[]=E&Ids[]=M usage:
// IDArray:=[]string{"E","M"} r.GET("reqpath").SetQueryD(gofight.D{`Ids[]`: IDArray})
func (rc *RequestConfig) SetQueryD(query D) *RequestConfig {
	var buf strings.Builder
	buf.WriteString("?")
	for k, v := range query {
		switch v := v.(type) {
		case string:
			buf.WriteString(k + "=" + v)
			buf.WriteString("&")
		case []string:
			for _, info := range v {
				buf.WriteString(k + "=" + info)
				buf.WriteString("&")
			}
		}
	}
	rc.Path += buf.String()[:len(buf.String())-1]
	return rc
}

// SetQuery sets the query parameters for the request configuration.
// It takes a map of query parameters and their values, and appends them
// to the existing path of the request configuration. If the path already
// contains query parameters, the new parameters are appended with an '&'.
// Otherwise, they are appended with a '?'.
//
// Parameters:
//
//	query (H): A map containing the query parameters and their values.
//
// Returns:
//
//	*RequestConfig: The updated request configuration with the query parameters set.
func (rc *RequestConfig) SetQuery(query H) *RequestConfig {
	f := make(url.Values)

	for k, v := range query {
		f.Set(k, v)
	}

	if strings.Contains(rc.Path, "?") {
		rc.Path = rc.Path + "&" + f.Encode()
	} else {
		rc.Path = rc.Path + "?" + f.Encode()
	}

	return rc
}

// SetBody sets the body of the request if the provided body string is not empty.
// It returns the updated RequestConfig instance.
//
// Parameters:
//   - body: A string representing the body content to be set.
//
// Returns:
//   - *RequestConfig: The updated RequestConfig instance.
func (rc *RequestConfig) SetBody(body string) *RequestConfig {
	if len(body) > 0 {
		rc.Body = body
	}

	return rc
}

// SetCookie sets the cookies for the request configuration.
// It takes a map of cookies and assigns it to the Cookies field of the RequestConfig
// if the provided map is not empty.
//
// Parameters:
//   - cookies: A map of cookies to be set.
//
// Returns:
//   - A pointer to the updated RequestConfig.
func (rc *RequestConfig) SetCookie(cookies H) *RequestConfig {
	if len(cookies) > 0 {
		rc.Cookies = cookies
	}

	return rc
}

func (rc *RequestConfig) initTest() (*http.Request, *httptest.ResponseRecorder) {
	qs := ""
	if strings.Contains(rc.Path, "?") {
		ss := strings.Split(rc.Path, "?")
		rc.Path = ss[0]
		qs = ss[1]
	}

	body := bytes.NewBufferString(rc.Body)

	req, _ := http.NewRequestWithContext(rc.Context, rc.Method, rc.Path, body)
	req.RequestURI = req.URL.RequestURI()

	if len(qs) > 0 {
		req.URL.RawQuery = qs
	}

	// Auto add user agent
	req.Header.Set(UserAgent, "Gofight-client/"+Version)

	if rc.Method == "POST" || rc.Method == "PUT" || rc.Method == "PATCH" {
		if strings.HasPrefix(rc.Body, "{") {
			req.Header.Set(ContentType, ApplicationJSON)
		} else {
			req.Header.Set(ContentType, ApplicationForm)
		}
	}

	if rc.ContentType != "" {
		req.Header.Set(ContentType, rc.ContentType)
	}

	if len(rc.Headers) > 0 {
		for k, v := range rc.Headers {
			req.Header.Set(k, v)
		}
	}

	if len(rc.Cookies) > 0 {
		for k, v := range rc.Cookies {
			req.AddCookie(&http.Cookie{
				Name:     k,
				Value:    v,
				HttpOnly: true,
				Secure:   false,
			})
		}
	}

	if rc.Debug {
		log.Printf("Request QueryString: %s", qs)
		log.Printf("Request Method: %s", rc.Method)
		log.Printf("Request Path: %s", rc.Path)
		log.Printf("Request Body: %s", rc.Body)
		log.Printf("Request Headers: %s", rc.Headers)
		log.Printf("Request Cookies: %s", rc.Cookies)
		log.Printf("Request Header: %s", req.Header)
	}

	w := httptest.NewRecorder()

	return req, w
}

// Run executes the HTTP request using the provided http.Handler and processes
// the response using the given ResponseFunc. It initializes the test request
// and response writer, serves the HTTP request, and then passes the HTTP
// response and request to the response function.
//
// Parameters:
//   - r: The http.Handler that will handle the HTTP request.
//   - response: A function that processes the HTTP response and request.
func (rc *RequestConfig) Run(r http.Handler, response ResponseFunc) {
	req, w := rc.initTest()
	r.ServeHTTP(w, req)
	response(
		HTTPResponse{
			w,
		},
		HTTPRequest{
			req,
		},
	)
}
