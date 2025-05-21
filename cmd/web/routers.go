package main

import(
	"net/http"
)


func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/codebox/view", app.snippetView)
	mux.HandleFunc("/codebox/create", app.snippetCreate)
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	return mux
}
