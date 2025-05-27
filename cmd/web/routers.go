package main

import(
	"net/http"
)

// 修改返回值为 http.Handler
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/codebox/view", app.snippetView)
	mux.HandleFunc("/codebox/create", app.snippetCreate)
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	return app.recoverPanic(app.logRequest(secureHeaders(mux)))
}
