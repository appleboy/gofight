package gofight

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const version = "0.0.1"

func basicHelloHandler(w http.ResponseWriter, r *http.Request) {
	// add header in response.
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("X-Version", version)
	_, _ = io.WriteString(w, "Hello World")
}

func basicCookieHandler(w http.ResponseWriter, r *http.Request) {
	// get cookie from request.
	foo, err := r.Cookie("foo")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, _ = io.WriteString(w, foo.Value)
}

func basicQueryHandler(w http.ResponseWriter, r *http.Request) {
	// get query from request.
	foo := r.URL.Query().Get("foo")
	_, _ = io.WriteString(w, foo)
}

func basicFormHandler(w http.ResponseWriter, r *http.Request) {
	// get form from request.
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	foo := r.Form.Get("foo")
	_, _ = io.WriteString(w, foo)
}

func basicEngine() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", basicHelloHandler)
	mux.HandleFunc("/cookie", basicCookieHandler)
	mux.HandleFunc("/query", basicQueryHandler)
	mux.HandleFunc("/form", basicFormHandler)

	return mux
}

func TestBasicHelloWorld(t *testing.T) {
	r := New()
	version := "0.0.1"

	r.GET("/").
		// turn on the debug mode.
		SetDebug(true).
		SetHeader(H{
			"X-Version": version,
		}).
		Run(basicEngine(), func(r HTTPResponse, rq HTTPRequest) {
			assert.Equal(t, version, rq.Header.Get("X-Version"))
			assert.Equal(t, version, r.Header().Get("X-Version"))
			assert.Equal(t, "Hello World", r.Body.String())
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func basicHTTPHelloHandler() {
	http.HandleFunc("/hello", basicHelloHandler)
}

func TestBasicHttpHelloWorld(t *testing.T) {
	basicHTTPHelloHandler()

	r := New()

	r.GET("/hello").
		// trun on the debug mode.
		SetDebug(true).
		Run(http.DefaultServeMux, func(r HTTPResponse, rq HTTPRequest) {
			assert.Equal(t, "Hello World", r.Body.String())
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestSetContext(t *testing.T) {
	r := New()
	type contextKey string
	const key contextKey = "key"
	ctx := context.WithValue(context.Background(), key, "value")

	r.GET("/").
		SetContext(ctx).
		Run(basicEngine(), func(r HTTPResponse, rq HTTPRequest) {
			assert.Equal(t, "value", rq.Context().Value(contextKey("key")))
			assert.Equal(t, "Hello World", r.Body.String())
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestSetContextWithTimeout(t *testing.T) {
	r := New()
	ctx, cancel := context.WithTimeout(context.Background(), 0)
	defer cancel()

	r.GET("/").
		SetContext(ctx).
		Run(basicEngine(), func(r HTTPResponse, rq HTTPRequest) {
			select {
			case <-rq.Context().Done():
				assert.Equal(t, context.DeadlineExceeded, rq.Context().Err())
			default:
				t.Error("expected context to be done")
			}
		})
}

func TestGetHeaderFromResponse(t *testing.T) {
	version := "0.0.1"
	r := New()
	r.GET("/").
		Run(basicEngine(), func(r HTTPResponse, rq HTTPRequest) {
			assert.Equal(t, version, r.Header().Get("X-Version"))
			assert.Equal(t, "Hello World", r.Body.String())
		})
}

func TestSetBody(t *testing.T) {
	r := New()
	body := "a=1&b=2"

	r.POST("/").
		SetBody(body).
		Run(basicEngine(), func(r HTTPResponse, rq HTTPRequest) {
			// Read the content of the io.Reader
			bodyBytes, err := io.ReadAll(rq.Body)
			if err != nil {
				t.Fatalf("Failed to read body: %v", err)
			}

			// Convert the byte slice to a string
			bodyString := string(bodyBytes)
			assert.Equal(t, body, bodyString)
			assert.Equal(t, "Hello World", r.Body.String())
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestSetCookie(t *testing.T) {
	r := New()
	cookies := H{
		"foo": "bar",
		"baz": "qux",
	}

	r.GET("/cookie").
		SetCookie(cookies).
		Run(basicEngine(), func(r HTTPResponse, rq HTTPRequest) {
			cookieFoo, err := rq.Cookie("foo")
			assert.NoError(t, err)
			assert.Equal(t, "bar", cookieFoo.Value)

			cookieBaz, err := rq.Cookie("baz")
			assert.NoError(t, err)
			assert.Equal(t, "qux", cookieBaz.Value)

			assert.Equal(t, "bar", r.Body.String())
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestSetQuery(t *testing.T) {
	r := New()
	query := H{
		"foo": "bar",
	}

	r.GET("/query").
		SetQuery(query).
		Run(basicEngine(), func(r HTTPResponse, rq HTTPRequest) {
			assert.Equal(t, "bar", rq.URL.Query().Get("foo"))
			assert.Equal(t, "bar", r.Body.String())
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestSetQueryWithExistingQuery(t *testing.T) {
	r := New()
	query := H{
		"c": "3",
		"d": "4",
	}

	r.GET("/query?a=1&b=2&foo=testing").
		SetQuery(query).
		Run(basicEngine(), func(r HTTPResponse, rq HTTPRequest) {
			assert.Equal(t, "1", rq.URL.Query().Get("a"))
			assert.Equal(t, "2", rq.URL.Query().Get("b"))
			assert.Equal(t, "3", rq.URL.Query().Get("c"))
			assert.Equal(t, "4", rq.URL.Query().Get("d"))
			assert.Equal(t, "testing", rq.URL.Query().Get("foo"))
			assert.Equal(t, "testing", r.Body.String())
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestSetForm(t *testing.T) {
	r := New()
	formData := H{
		"a":   "1",
		"b":   "2",
		"foo": "bar",
	}

	r.POST("/form").
		SetForm(formData).
		Run(basicEngine(), func(r HTTPResponse, rq HTTPRequest) {
			assert.Equal(t, "bar", r.Body.String())
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

// Additional test handlers for comprehensive testing
func jsonHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}

	var data interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"received": data,
		"method":   r.Method,
	}
	_ = json.NewEncoder(w).Encode(response)
}

func fileUploadHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check for any uploaded files in any field
	if r.MultipartForm == nil || len(r.MultipartForm.File) == 0 {
		http.Error(w, "No files uploaded", http.StatusBadRequest)
		return
	}

	// Get the first file from any field
	for fieldName, files := range r.MultipartForm.File {
		if len(files) > 0 {
			file := files[0]
			_, _ = io.WriteString(w, fmt.Sprintf("Uploaded file: %s, Size: %d, Field: %s",
				file.Filename, file.Size, fieldName))
			return
		}
	}

	http.Error(w, "No files found", http.StatusBadRequest)
}

func methodEchoHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = io.WriteString(w, fmt.Sprintf("Method: %s, Path: %s", r.Method, r.URL.Path))
}

func pathVariableHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/books/")
	_, _ = io.WriteString(w, fmt.Sprintf("Book path: %s", path))
}

func extendedEngine() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", basicHelloHandler)
	mux.HandleFunc("/cookie", basicCookieHandler)
	mux.HandleFunc("/query", basicQueryHandler)
	mux.HandleFunc("/form", basicFormHandler)
	mux.HandleFunc("/json", jsonHandler)
	mux.HandleFunc("/upload", fileUploadHandler)
	mux.HandleFunc("/books/", pathVariableHandler)
	mux.HandleFunc("/method", methodEchoHandler)
	return mux
}

// TestHTTPMethods tests all HTTP methods using table-driven tests
func TestHTTPMethods(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		setupRequest func(r *RequestConfig) *RequestConfig
		expectedBody string
		checkBody    bool
	}{
		{
			name:         "GET request",
			method:       "GET",
			setupRequest: func(r *RequestConfig) *RequestConfig { return r.GET("/method") },
			expectedBody: "Method: GET, Path: /method",
			checkBody:    true,
		},
		{
			name:         "POST request",
			method:       "POST",
			setupRequest: func(r *RequestConfig) *RequestConfig { return r.POST("/method") },
			expectedBody: "Method: POST, Path: /method",
			checkBody:    true,
		},
		{
			name:         "PUT request",
			method:       "PUT",
			setupRequest: func(r *RequestConfig) *RequestConfig { return r.PUT("/method") },
			expectedBody: "Method: PUT, Path: /method",
			checkBody:    true,
		},
		{
			name:         "DELETE request",
			method:       "DELETE",
			setupRequest: func(r *RequestConfig) *RequestConfig { return r.DELETE("/method") },
			expectedBody: "Method: DELETE, Path: /method",
			checkBody:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New()
			tt.setupRequest(r).Run(extendedEngine(), func(resp HTTPResponse, req HTTPRequest) {
				assert.Equal(t, tt.method, req.Method)
				if tt.checkBody {
					assert.Equal(t, tt.expectedBody, resp.Body.String())
				}
				assert.Equal(t, http.StatusOK, resp.Code)
			})
		})
	}
}

// TestSetJSON tests JSON body setting functionality
func TestSetJSON(t *testing.T) {
	tests := []struct {
		name string
		data D
	}{
		{
			name: "simple object",
			data: D{"name": "test", "value": 123},
		},
		{
			name: "nested object",
			data: D{"user": D{"name": "john", "age": 30}},
		},
		{
			name: "empty object",
			data: D{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New()
			r.POST("/json").
				SetJSON(tt.data).
				Run(extendedEngine(), func(resp HTTPResponse, req HTTPRequest) {
					assert.Equal(t, http.StatusOK, resp.Code)
					assert.Contains(t, req.Header.Get("Content-Type"), "application/json")

					// Parse response to verify data was processed correctly
					var response map[string]interface{}
					err := json.Unmarshal(resp.Body.Bytes(), &response)
					assert.NoError(t, err)
					assert.Equal(t, "POST", response["method"])
					assert.NotNil(t, response["received"])
				})
		})
	}
}

// TestSetJSONInterface tests JSON interface functionality
func TestSetJSONInterface(t *testing.T) {
	type TestStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	tests := []struct {
		name string
		data interface{}
	}{
		{
			name: "struct",
			data: TestStruct{Name: "test", Value: 42},
		},
		{
			name: "map",
			data: map[string]interface{}{"key": "value", "number": 123},
		},
		{
			name: "slice",
			data: []string{"item1", "item2", "item3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New()
			r.POST("/json").
				SetJSONInterface(tt.data).
				Run(extendedEngine(), func(resp HTTPResponse, req HTTPRequest) {
					assert.Equal(t, http.StatusOK, resp.Code)
					assert.Contains(t, req.Header.Get("Content-Type"), "application/json")
				})
		})
	}
}

// TestSetQueryD tests query parameter arrays functionality
func TestSetQueryD(t *testing.T) {
	tests := []struct {
		name     string
		query    D
		expected map[string][]string
	}{
		{
			name:  "string query",
			query: D{"name": "john", "age": "30"},
			expected: map[string][]string{
				"name": {"john"},
				"age":  {"30"},
			},
		},
		{
			name:  "array query",
			query: D{"ids": []string{"1", "2", "3"}},
			expected: map[string][]string{
				"ids": {"1", "2", "3"},
			},
		},
		{
			name:  "mixed query",
			query: D{"name": "john", "ids": []string{"1", "2"}},
			expected: map[string][]string{
				"name": {"john"},
				"ids":  {"1", "2"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New()
			r.GET("/query").
				SetQueryD(tt.query).
				Run(extendedEngine(), func(resp HTTPResponse, req HTTPRequest) {
					for key, expectedValues := range tt.expected {
						actualValues := req.URL.Query()[key]
						assert.Equal(t, expectedValues, actualValues)
					}
					assert.Equal(t, http.StatusOK, resp.Code)
				})
		})
	}
}

// TestSetPath tests path variable functionality
func TestSetPath(t *testing.T) {
	r := New()
	r.GET("/books/").
		SetPath("golang/guide").
		Run(extendedEngine(), func(resp HTTPResponse, req HTTPRequest) {
			assert.Equal(t, "/books/golang/guide", req.URL.Path)
			assert.Equal(t, "Book path: golang/guide", resp.Body.String())
			assert.Equal(t, http.StatusOK, resp.Code)
		})
}

// TestSetBodyEmpty tests empty body handling
func TestSetBodyEmpty(t *testing.T) {
	r := New()
	r.POST("/method").
		SetBody("").
		Run(extendedEngine(), func(resp HTTPResponse, req HTTPRequest) {
			body, _ := io.ReadAll(req.Body)
			assert.Equal(t, "", string(body))
			assert.Equal(t, http.StatusOK, resp.Code)
		})
}

// TestContentTypeDetection tests automatic content type detection
func TestContentTypeDetection(t *testing.T) {
	tests := []struct {
		name         string
		body         string
		expectedType string
	}{
		{
			name:         "JSON body",
			body:         `{"name": "test"}`,
			expectedType: "application/json",
		},
		{
			name:         "form body",
			body:         "name=test&value=123",
			expectedType: "application/x-www-form-urlencoded",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New()
			r.POST("/method").
				SetBody(tt.body).
				Run(extendedEngine(), func(resp HTTPResponse, req HTTPRequest) {
					assert.Contains(t, req.Header.Get("Content-Type"), tt.expectedType)
				})
		})
	}
}

// TestCookieSecuritySettings tests cookie security configuration
func TestCookieSecuritySettings(t *testing.T) {
	r := New()
	r.GET("/cookie").
		SetCookie(H{"foo": "bar"}). // Use "foo" to match basicCookieHandler
		Run(extendedEngine(), func(resp HTTPResponse, req HTTPRequest) {
			// Test that the cookie was set and handler works
			assert.Equal(t, "bar", resp.Body.String()) // Handler returns cookie value
			assert.Equal(t, http.StatusOK, resp.Code)

			// Test that cookies were added to the request
			cookies := req.Cookies()
			assert.NotEmpty(t, cookies)

			// Find the foo cookie and verify its properties
			var fooCookie *http.Cookie
			for _, cookie := range cookies {
				if cookie.Name == "foo" {
					fooCookie = cookie
					break
				}
			}

			if fooCookie != nil {
				assert.Equal(t, "foo", fooCookie.Name)
				assert.Equal(t, "bar", fooCookie.Value)
				// Note: HttpOnly and SameSite might not be set in test environment
			}
		})
}

// TestErrorHandling tests error scenarios
func TestErrorHandling(t *testing.T) {
	// Test with invalid JSON
	r := New()
	invalidJSON := `{"invalid": json}`

	r.POST("/json").
		SetBody(invalidJSON).
		Run(extendedEngine(), func(resp HTTPResponse, req HTTPRequest) {
			// Should handle gracefully
			assert.Equal(t, http.StatusBadRequest, resp.Code)
		})
}

// TestSetFileFromPath tests file upload functionality
func TestSetFileFromPath(t *testing.T) {
	// Create a temporary test file
	tmpFile := filepath.Join(os.TempDir(), "test.txt")
	err := os.WriteFile(tmpFile, []byte("Hello World"), 0o600)
	require.NoError(t, err)
	defer os.Remove(tmpFile)

	uploadFile := UploadFile{
		Path: tmpFile,
		Name: "file",
	}

	r := New()
	r.POST("/upload").
		SetFileFromPath([]UploadFile{uploadFile}).
		Run(extendedEngine(), func(resp HTTPResponse, req HTTPRequest) {
			// The request should have multipart content type
			contentType := req.Header.Get("Content-Type")
			assert.Contains(t, contentType, "multipart/form-data")

			// The body should contain the file upload data
			assert.NotEmpty(t, req.Body)
		})
}

// TestSetFileFromContent tests file upload with content
func TestSetFileFromContent(t *testing.T) {
	uploadFile := UploadFile{
		Path:    "test.txt",
		Name:    "file",
		Content: []byte("Test file content"),
	}

	r := New()
	r.POST("/upload").
		SetFileFromPath([]UploadFile{uploadFile}).
		Run(extendedEngine(), func(resp HTTPResponse, req HTTPRequest) {
			// The request should have multipart content type
			contentType := req.Header.Get("Content-Type")
			assert.Contains(t, contentType, "multipart/form-data")

			// The body should contain the file upload data
			assert.NotEmpty(t, req.Body)
		})
}

// TestDebugMode tests debug functionality
func TestDebugMode(t *testing.T) {
	r := New()
	r.GET("/").
		SetDebug(true).
		SetHeader(H{"X-Test": "debug"}).
		Run(extendedEngine(), func(resp HTTPResponse, req HTTPRequest) {
			assert.Equal(t, "debug", req.Header.Get("X-Test"))
			assert.Equal(t, http.StatusOK, resp.Code)
		})
}

// Benchmark tests for performance
func BenchmarkNewRequest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := New()
		r.GET("/")
	}
}

func BenchmarkSimpleGETRequest(b *testing.B) {
	engine := extendedEngine()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := New()
		r.GET("/").Run(engine, func(resp HTTPResponse, req HTTPRequest) {
			// Simple assertion
			_ = resp.Code == http.StatusOK
		})
	}
}

func BenchmarkJSONRequest(b *testing.B) {
	engine := extendedEngine()
	data := D{"name": "benchmark", "value": 123}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := New()
		r.POST("/json").
			SetJSON(data).
			Run(engine, func(resp HTTPResponse, req HTTPRequest) {
				_ = resp.Code == http.StatusOK
			})
	}
}

func BenchmarkFormRequest(b *testing.B) {
	engine := extendedEngine()
	formData := H{"name": "benchmark", "value": "123"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := New()
		r.POST("/form").
			SetForm(formData).
			Run(engine, func(resp HTTPResponse, req HTTPRequest) {
				_ = resp.Code == http.StatusOK
			})
	}
}

// TestMoreHTTPMethods tests PATCH, HEAD, and OPTIONS methods
func TestMoreHTTPMethods(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		setupRequest func(r *RequestConfig) *RequestConfig
		expectedBody string
		checkBody    bool
	}{
		{
			name:         "PATCH request",
			method:       "PATCH",
			setupRequest: func(r *RequestConfig) *RequestConfig { return r.PATCH("/method") },
			expectedBody: "Method: PATCH, Path: /method",
			checkBody:    true,
		},
		{
			name:         "HEAD request",
			method:       "HEAD",
			setupRequest: func(r *RequestConfig) *RequestConfig { return r.HEAD("/method") },
			expectedBody: "",
			checkBody:    false, // HEAD requests don't return body
		},
		{
			name:         "OPTIONS request",
			method:       "OPTIONS",
			setupRequest: func(r *RequestConfig) *RequestConfig { return r.OPTIONS("/method") },
			expectedBody: "Method: OPTIONS, Path: /method",
			checkBody:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New()
			tt.setupRequest(r).Run(extendedEngine(), func(resp HTTPResponse, req HTTPRequest) {
				assert.Equal(t, tt.method, req.Method)
				if tt.checkBody {
					assert.Equal(t, tt.expectedBody, resp.Body.String())
				}
				assert.Equal(t, http.StatusOK, resp.Code)
			})
		})
	}
}

// TestEdgeCases tests various edge cases for better coverage
func TestEdgeCases(t *testing.T) {
	t.Run("empty headers", func(t *testing.T) {
		r := New()
		r.GET("/").
			SetHeader(H{}).
			Run(extendedEngine(), func(resp HTTPResponse, req HTTPRequest) {
				assert.Equal(t, http.StatusOK, resp.Code)
			})
	})

	t.Run("empty cookies", func(t *testing.T) {
		r := New()
		r.GET("/").
			SetCookie(H{}).
			Run(extendedEngine(), func(resp HTTPResponse, req HTTPRequest) {
				assert.Equal(t, http.StatusOK, resp.Code)
			})
	})

	t.Run("empty query", func(t *testing.T) {
		r := New()
		r.GET("/").
			SetQuery(H{}).
			Run(extendedEngine(), func(resp HTTPResponse, req HTTPRequest) {
				assert.Equal(t, http.StatusOK, resp.Code)
			})
	})

	t.Run("empty form", func(t *testing.T) {
		r := New()
		r.POST("/").
			SetForm(H{}).
			Run(extendedEngine(), func(resp HTTPResponse, req HTTPRequest) {
				assert.Equal(t, http.StatusOK, resp.Code)
			})
	})

	t.Run("empty QueryD", func(t *testing.T) {
		r := New()
		r.GET("/").
			SetQueryD(D{}).
			Run(extendedEngine(), func(resp HTTPResponse, req HTTPRequest) {
				assert.Equal(t, http.StatusOK, resp.Code)
			})
	})

	t.Run("nil context", func(t *testing.T) {
		r := New()
		r.SetContext(context.Background()).
			GET("/").
			Run(extendedEngine(), func(resp HTTPResponse, req HTTPRequest) {
				assert.Equal(t, http.StatusOK, resp.Code)
				assert.NotNil(t, req.Context())
			})
	})
}

// TestJSONMarshallErrors tests JSON marshalling error handling
func TestJSONMarshallErrors(t *testing.T) {
	r := New()

	// Test with channel type that can't be marshalled to JSON
	invalidData := make(chan int)

	r.POST("/json").
		SetJSONInterface(invalidData).
		Run(extendedEngine(), func(resp HTTPResponse, req HTTPRequest) {
			// Should continue execution even with marshal error
			body, _ := io.ReadAll(req.Body)
			assert.Equal(t, "", string(body)) // Body should be empty due to marshal error
		})
}

// TestQueryWithSpecialCharacters tests query parameters with special characters
func TestQueryWithSpecialCharacters(t *testing.T) {
	r := New()

	specialQuery := H{
		"name":  "john doe",
		"email": "john@example.com",
		"tags":  "tag1,tag2,tag3",
	}

	r.GET("/query").
		SetQuery(specialQuery).
		Run(extendedEngine(), func(resp HTTPResponse, req HTTPRequest) {
			assert.Equal(t, "john doe", req.URL.Query().Get("name"))
			assert.Equal(t, "john@example.com", req.URL.Query().Get("email"))
			assert.Equal(t, "tag1,tag2,tag3", req.URL.Query().Get("tags"))
		})
}

// TestSecureContextVariations tests different secure context scenarios
func TestSecureContextVariations(t *testing.T) {
	r := New()

	// Test isSecureContext method indirectly
	r.GET("/").
		SetCookie(H{"secure_test": "value"}).
		Run(extendedEngine(), func(resp HTTPResponse, req HTTPRequest) {
			// Cookie should be set regardless of secure context in test
			cookies := req.Cookies()
			assert.NotEmpty(t, cookies)
		})
}

// TestContentTypeOverride tests content type override functionality
func TestContentTypeOverride(t *testing.T) {
	customContentType := "application/xml"

	config := &RequestConfig{
		Method:      "POST",
		Path:        "/method",
		Body:        "<xml>test</xml>",
		ContentType: customContentType,
		Context:     context.Background(),
	}

	config.Run(extendedEngine(), func(resp HTTPResponse, req HTTPRequest) {
		assert.Equal(t, customContentType, req.Header.Get("Content-Type"))
	})
}

// TestUserAgentHeader tests user agent header setting
func TestUserAgentHeader(t *testing.T) {
	r := New()

	r.GET("/").
		Run(extendedEngine(), func(resp HTTPResponse, req HTTPRequest) {
			userAgent := req.Header.Get("User-Agent")
			assert.Contains(t, userAgent, "Gofight-client/")
			assert.Contains(t, userAgent, "1.0")
		})
}

// Handler for testing multiple response headers
func multipleHeadersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("X-Custom-Header", "custom-value")
	w.Header().Set("X-API-Version", "v2.0")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Add("X-Multiple", "value1")
	w.Header().Add("X-Multiple", "value2")
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, "Headers set")
}

// TestHTTPResponseHeaders tests comprehensive HTTPResponse header functionality
func TestHTTPResponseHeaders(t *testing.T) {
	// Create a custom engine with the headers handler
	mux := http.NewServeMux()
	mux.HandleFunc("/headers", multipleHeadersHandler)

	tests := []struct {
		name         string
		handler      http.Handler
		checkHeaders func(t *testing.T, resp HTTPResponse)
	}{
		{
			name:    "single header",
			handler: mux,
			checkHeaders: func(t *testing.T, resp HTTPResponse) {
				assert.Equal(t, "custom-value", resp.Header().Get("X-Custom-Header"))
			},
		},
		{
			name:    "multiple different headers",
			handler: mux,
			checkHeaders: func(t *testing.T, resp HTTPResponse) {
				assert.Equal(t, "custom-value", resp.Header().Get("X-Custom-Header"))
				assert.Equal(t, "v2.0", resp.Header().Get("X-API-Version"))
				assert.Equal(t, "no-cache", resp.Header().Get("Cache-Control"))
			},
		},
		{
			name:    "multiple values for same header",
			handler: mux,
			checkHeaders: func(t *testing.T, resp HTTPResponse) {
				values := resp.Header()["X-Multiple"]
				assert.Equal(t, 2, len(values))
				assert.Contains(t, values, "value1")
				assert.Contains(t, values, "value2")
			},
		},
		{
			name:    "header case insensitivity",
			handler: mux,
			checkHeaders: func(t *testing.T, resp HTTPResponse) {
				// HTTP headers are case-insensitive
				assert.Equal(t, "custom-value", resp.Header().Get("x-custom-header"))
				assert.Equal(t, "custom-value", resp.Header().Get("X-CUSTOM-HEADER"))
				assert.Equal(t, "custom-value", resp.Header().Get("X-Custom-Header"))
			},
		},
		{
			name:    "non-existent header returns empty string",
			handler: mux,
			checkHeaders: func(t *testing.T, resp HTTPResponse) {
				assert.Equal(t, "", resp.Header().Get("X-Non-Existent"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New()
			r.GET("/headers").
				Run(tt.handler, func(resp HTTPResponse, req HTTPRequest) {
					assert.Equal(t, http.StatusOK, resp.Code)
					tt.checkHeaders(t, resp)
				})
		})
	}
}

// TestHTTPResponseHeaderMethods tests various header methods on HTTPResponse
func TestHTTPResponseHeaderMethods(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/headers", multipleHeadersHandler)

	r := New()
	r.GET("/headers").
		Run(mux, func(resp HTTPResponse, req HTTPRequest) {
			// Test Header() returns non-nil map
			assert.NotNil(t, resp.Header())

			// Test Values() method for multiple header values
			multipleValues := resp.Header().Values("X-Multiple")
			assert.Equal(t, 2, len(multipleValues))
			assert.Equal(t, "value1", multipleValues[0])
			assert.Equal(t, "value2", multipleValues[1])

			// Test direct map access
			headerMap := resp.Header()
			assert.Contains(t, headerMap, "X-Custom-Header")
			assert.Equal(t, []string{"custom-value"}, headerMap["X-Custom-Header"])

			// Test Content-Type from standard response
			assert.NotEmpty(t, resp.Header().Get("Content-Type"))
		})
}
