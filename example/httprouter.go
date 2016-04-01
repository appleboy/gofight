package example

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"fmt"
)

func HttpRouterHelloHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Hello World")
}

func HttpRouterEngine() http.Handler {
	r := httprouter.New()

	r.GET("/", HttpRouterHelloHandler)

	return r
}
