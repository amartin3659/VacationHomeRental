package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	env "github.com/amartin3659/VacationHomeRental/cmd"
	"github.com/amartin3659/VacationHomeRental/internal/config"
	"github.com/amartin3659/VacationHomeRental/internal/driver"
	"github.com/amartin3659/VacationHomeRental/internal/handlers"
	"github.com/amartin3659/VacationHomeRental/internal/helpers"
	"github.com/amartin3659/VacationHomeRental/internal/models"
	"github.com/amartin3659/VacationHomeRental/internal/render"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
  db, err := run()
  if err != nil {
    log.Fatal(err)
  }

  defer db.SQL.Close()

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

func run() (*driver.DB, error) {

  // Data to be available in the session
  gob.Register(models.Reservation{})
  gob.Register(models.User{})
  gob.Register(models.Bungalow{})
  gob.Register(models.BungalowRestriction{})
  gob.Register(models.Restriction{})

	app.InProduction = false

  infoLog = log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime)
  app.InfoLog = infoLog
  
  errorLog = log.New(os.Stdout, "[ERROR]\t", log.Ldate|log.Ltime|log.Lshortfile)
  app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

  // connecting to database
  log.Println("Conecting to database...")
  env.SetPass()
  connStr := fmt.Sprintf("host=localhost port=5432 dbname=mygowebapp user=%s password=%s", env.GetUser(), env.GetPass())
  db, err := driver.ConnectSQL(connStr)
  if err != nil {
    log.Fatal("No connection to database! Terminating ...")
    return nil, err
  }
  log.Println("Successfully connected to database.")

	// create a template cache
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatalln("Error creating template cache", err)
    return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	render.NewRenderer(&app)

  helpers.NewHelpers(&app)
  return db, nil
}

