package main

import (
	"net/http"

	"github.com/justinas/alice"
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

	dynamic := alice.New(app.sessionManager.LoadAndSave)
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/codebox/view/:id", dynamic.ThenFunc(app.snippetView))
	router.Handler(http.MethodGet, "/codebox/create", dynamic.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/codebox/create", dynamic.ThenFunc(app.snippetCreatePost))
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(router)
}
