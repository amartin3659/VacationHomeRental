package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/amartin3659/VacationHomeRental/internal/config"
	"github.com/amartin3659/VacationHomeRental/internal/handlers"
	"github.com/amartin3659/VacationHomeRental/internal/models"
	"github.com/amartin3659/VacationHomeRental/internal/render"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {
  err := run()
  if err != nil {
    log.Fatal(err)
  }
	fmt.Println("Starting server on port: ", portNumber)

	src := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = src.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}

func run() error {

  // Data to be available in the session
  gob.Register(models.Reservation{})

	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// create a template cache
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatalln("Error creating template cache", err)
    return err
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

  return nil
}

