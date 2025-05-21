package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)


type application struct{
	errorLog *log.Logger
	infoLog *log.Logger
}

func main() {

	// flag command-line
	addr := flag.String("addr", ":8080", "HTTP network address")
	flag.Parse()

	// logging
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// application
	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
	}

	infoLog.Printf("starting server on %s", *addr)
	server := http.Server{
		Addr:    *addr,
		Handler: app.routes(),
		ErrorLog: errorLog,
	}
	err := server.ListenAndServe()
	errorLog.Fatal(err)
}
