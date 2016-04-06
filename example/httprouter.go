package example

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func httpRouterHelloHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Hello World")
}

// HTTPRouterEngine is httprouter router.
func HTTPRouterEngine() http.Handler {
	r := httprouter.New()

	r.GET("/", httpRouterHelloHandler)

	return r
}
