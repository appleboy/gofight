package example

import (
	"github.com/appleboy/gofight/v2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestFiberEngine(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "TestHelloWorld",
			path: "/",
			want: "God Love the World ! 👴 john is 75 years old~",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gofight.New()
			r.GET(tt.path).SetQueryD(gofight.D{
				"name": "john",
				"age":  "75",
			}).
				SetDebug(true).
				RunX(FiberEngine(), func(res gofight.HTTPResponse, req gofight.HTTPRequest) {
					assert.Equal(t, tt.want, res.Body.String())
					assert.Equal(t, http.StatusOK, res.Code)
				})
		})
	}
}
