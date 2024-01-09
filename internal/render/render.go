package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/amartin3659/VacationHomeRental/internal/config"
	"github.com/amartin3659/VacationHomeRental/internal/models"
	"github.com/justinas/nosurf"
)

// AddDefaultData contains data which will be added to data sent to templates
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
  td.Flash = app.Session.PopString(r.Context(), "flash")
  td.Error = app.Session.PopString(r.Context(), "error")
  td.Warning = app.Session.PopString(r.Context(), "warning")
  td.CSRFToken = nosurf.Token(r)
  return td
}

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
  app = a
}

// RenderTemplate serves as a wrapper and renders
// a layout and a template from folder /templates to a desired writer
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {
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
    log.Fatalln("Template not in cache")
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
}

func CreateTemplateCache() (map[string]*template.Template, error) {
  cache := map[string]*template.Template{}

  // get all files *-page.html from /templates folder
  pages, err := filepath.Glob("./templates/*-page.html")
  if err != nil {
    return cache, err
  }

  // range through the slice of *-page.html
  for _, page := range pages {
    name := filepath.Base(page)
    ts, err := template.New(name).ParseFiles(page)
    if err != nil {
      return cache, err
    }

    matches, err := filepath.Glob("./templates/*-layout.html")
    if err != nil {
      return cache, err
    }

    if len(matches) > 0 {
      ts, err = ts.ParseGlob("./templates/*-layout.html")
      if err != nil {
        return cache, err
      }
    }
    cache[name] = ts
  }

  return cache, nil
}
