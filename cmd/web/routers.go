package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)


func (app *application) routes() http.Handler {
	// Initialize the router.
	router := httprouter.New()

	// 404
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath",
		http.StripPrefix("/static", fileServer))
	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/codebox/view/:id", app.snippetView)
	router.HandlerFunc(http.MethodGet, "/codebox/create", app.snippetCreate)
	router.HandlerFunc(http.MethodPost, "/codebox/create", app.snippetCreatePost)
	return app.recoverPanic(app.logRequest(secureHeaders(router)))
}
