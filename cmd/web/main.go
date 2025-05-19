package main

import (
	"net/http"
)



func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/codebox/view", snippetView)
	mux.HandleFunc("/codebox/create", snippetCreate)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	server.ListenAndServe()
}
