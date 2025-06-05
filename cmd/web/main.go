package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	// import models package
	"github.com/aikwen/codebox/internal/models"
	"github.com/go-playground/form/v4"

	_ "github.com/go-sql-driver/mysql"
)


type application struct{
	errorLog *log.Logger
	infoLog *log.Logger
	snippets *models.SnippetModel
	templateCache map[string]*template.Template
	formDecoder *form.Decoder
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil{
		return nil, err
	}
	return db, nil
}



func main() {

	// flag command-line
	// addr
	addr := flag.String("addr", ":8080", "HTTP network address")
	// mysql dsn
	dsn := flag.String("dsn", "codebox:codebox@/codebox?parseTime=true","mysql data source name")
	flag.Parse()

	// logging
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// mysql
	db, err := openDB(*dsn)
	if err != nil{
		errorLog.Fatal(err)
	}

	defer db.Close()

	// template cache
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	formDecoder := form.NewDecoder()
	// application
	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
		snippets: &models.SnippetModel{DB: db},
		templateCache: templateCache,
		formDecoder: formDecoder,
	}

	infoLog.Printf("starting server on %s", *addr)
	server := http.Server{
		Addr:    *addr,
		Handler: app.routes(),
		ErrorLog: errorLog,
	}
	err = server.ListenAndServe()
	errorLog.Fatal(err)
}
