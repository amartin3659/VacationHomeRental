package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/amartin3659/VacationHomeRental/internal/config"
	"github.com/amartin3659/VacationHomeRental/internal/models"
	"github.com/amartin3659/VacationHomeRental/internal/render"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justinas/nosurf"
)

var app config.AppConfig
var session *scs.SessionManager
var pathToTemplates = "./../../templates"

func TestMain(m *testing.M) {

  // Data to be available in the session
  gob.Register(models.Reservation{})

	app.InProduction = false

  infoLog := log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime)
  app.InfoLog = infoLog
  
  errorLog := log.New(os.Stdout, "[ERROR]\t", log.Ldate|log.Ltime|log.Lshortfile)
  app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// create a template cache
	tc, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatalln("Error creating template cache", err)
	}

	app.TemplateCache = tc
	app.UseCache = true 

	repo := NewTestRepo(&app)
	NewHandlers(repo)

	render.NewRenderer(&app)

  os.Exit(m.Run())
}

func getRoutes() http.Handler {

  mux := chi.NewRouter()

  mux.Use(middleware.Recoverer)
  // mux.Use(NoSurf)
  mux.Use(SessionLoad)

  mux.Get("/", Repo.Home)
  mux.Get("/about", Repo.About)
  mux.Get("/contact", Repo.Contact)
  mux.Get("/eremite", Repo.Eremite)
  mux.Get("/couple", Repo.Couple)
  mux.Get("/family", Repo.Family)
  mux.Get("/reservation", Repo.Reservation)
  mux.Post("/reservation", Repo.PostReservation)
  mux.Post("/reservation-json", Repo.ReservationJSON)
  mux.Get("/make-reservation", Repo.MakeReservation)
  mux.Post("/make-reservation", Repo.PostMakeReservation)
  mux.Get("/reservation-overview", Repo.ReservationOverview)

  fileServer := http.FileServer(http.Dir("./static/"))
  mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

  return mux
}

// NoSurf serves as a CSRF protection middleware
func NoSurf(next http.Handler) http.Handler {
  csrfHandler := nosurf.New(next)

  csrfHandler.SetBaseCookie(http.Cookie{
    HttpOnly: true,
    Path: "/",
    Secure: app.InProduction,
    SameSite: http.SameSiteLaxMode,
  })

  return csrfHandler
}

// SessionLoad loads, saves session data for each request
func SessionLoad(next http.Handler) http.Handler {
  return session.LoadAndSave(next)
}

func CreateTestTemplateCache() (map[string]*template.Template, error) {
  cache := map[string]*template.Template{}

  // get all files *-page.html from /templates folder
  pages, err := filepath.Glob(fmt.Sprintf("%s/*-page.html", pathToTemplates))
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
