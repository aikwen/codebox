package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
	// import models package
	"github.com/aikwen/codebox/internal/models"
	"github.com/go-playground/form/v4"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"

	_ "github.com/go-sql-driver/mysql"
)


type application struct{
	errorLog *log.Logger
	infoLog *log.Logger
	snippets *models.SnippetModel
	templateCache map[string]*template.Template
	formDecoder *form.Decoder
	sessionManager *scs.SessionManager
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


	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	sessionManager.Cookie.Secure = true
	// application
	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
		snippets: &models.SnippetModel{DB: db},
		templateCache: templateCache,
		formDecoder: formDecoder,
		sessionManager: sessionManager,
	}

	infoLog.Printf("starting server on %s", *addr)

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	server := http.Server{
		Addr:    *addr,
		Handler: app.routes(),
		ErrorLog: errorLog,
		TLSConfig: tlsConfig,
		IdleTimeout: time.Minute,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	err = server.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}
