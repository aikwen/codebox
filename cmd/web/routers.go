package main

import (
	"net/http"

	"github.com/aikwen/codebox/ui"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	// Initialize the router.
	router := httprouter.New()

	// 404
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.FS(ui.Files))

	router.Handler(http.MethodGet, "/static/*filepath",
		fileServer)
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/codebox/view/:id", dynamic.ThenFunc(app.snippetView))
	// user
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))
	// Protected (authenticated-only) application routes,
	protected := dynamic.Append(app.requireAuthentication)
	router.Handler(http.MethodGet, "/codebox/create", protected.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/codebox/create", protected.ThenFunc(app.snippetCreatePost))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(router)
}
