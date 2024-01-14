package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/amartin3659/VacationHomeRental/internal/config"
	"github.com/amartin3659/VacationHomeRental/internal/models"
	"github.com/justinas/nosurf"
)

var functions = template.FuncMap{
	"humanReadableDate": HumanReadableDate,
	"formatDate":        FormatDate,
	"iterate":           Iterate,
	"add":               Add,
}

// HumanReadableDate returns a time value in the YYYY-MM-DD format
func HumanReadableDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// FormatDate returns a time value in the YYYY-MM-DD format
func FormatDate(t time.Time, f string) string {
	return t.Format(f)
}

// Iterate creates and returns a slice of ints, starting at 1, going to count
func Iterate(count int) []int {
	var i int
	var items []int
	for i = 0; i < count; i++ {
		items = append(items, i)
	}

	return items
}

func Add(a, b int) int {
	return a + b
}

// AddDefaultData contains data which will be added to data sent to templates
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Success = app.Session.PopString(r.Context(), "success")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	if app.Session.Exists(r.Context(), "user_id") {
		td.IsAuthenticated = 1
	}
	return td
}

var app *config.AppConfig
var pathToTemplates = "./templates"

// NewRenderer sets the config for the template package
func NewRenderer(a *config.AppConfig) {
	app = a
}

// Template serves as a wrapper and renders
// a layout and a template from folder /templates to a desired writer
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
	var tc map[string]*template.Template
	if app.UseCache {
		// get the template from app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	// get the right template from cache
	t, ok := tc[tmpl]
	if !ok {
		return errors.New("template not in cache for some reason")
	}

	// store result in a buffer and double check if it is a valid value
	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	err := t.Execute(buf, td)
	if err != nil {
		log.Println(err)
	}

	// render that template
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}

	return nil
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// get all files *-page.html from /templates folder
	pages, err := filepath.Glob(fmt.Sprintf("%s/*-page.html", pathToTemplates))
	if err != nil {
		return cache, err
	}

	// range through the slice of *-page.html
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return cache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*-layout.html", pathToTemplates))
		if err != nil {
			return cache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*-layout.html", pathToTemplates))
			if err != nil {
				return cache, err
			}
		}
		cache[name] = ts
	}

	return cache, nil
}
