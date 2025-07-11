package http

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func newRouter(html *html, ezAPI *api) http.Handler {
	router := httprouter.New()
	router.GET("/", html.Index)
	router.GET("/versions", ezAPI.versions())
	router.GET("/sum", ezAPI.sumScore())
	router.ServeFiles("/static/*filepath", http.Dir("assets/static"))
	return recovery(noDirListing(router))
}
