package main

import (
	"html/template"
	"path/filepath"
	"time"
	"io/fs"

	"github.com/aikwen/codebox/ui"
	"github.com/aikwen/codebox/internal/models"
)

// Define a templateData type to act as the holding structure for
// any dynamic data that we want to pass to our HTML templates.
// At the moment it only contains one field, but we'll add more
// to it as the build progresses.
type templateData struct {
	Snippet *models.Snippet
	Snippets []*models.Snippet
	CurrentYear int
	Form any
	Flash string
	IsAuthenticated bool
	CSRFToken string // Add a CSRFToken field.
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate":humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := fs.Glob(ui.Files,"html/pages/*.html")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		// Extract the file name (like 'home.html') from the full filepath
		name := filepath.Base(page)
		patterns := []string{
			"html/base.html",
			"html/partials/*.html",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}

