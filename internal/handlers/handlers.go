package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/amartin3659/VacationHomeRental/internal/config"
	"github.com/amartin3659/VacationHomeRental/internal/driver"
	"github.com/amartin3659/VacationHomeRental/internal/forms"
	"github.com/amartin3659/VacationHomeRental/internal/helpers"
	"github.com/amartin3659/VacationHomeRental/internal/models"
	"github.com/amartin3659/VacationHomeRental/internal/render"
	"github.com/amartin3659/VacationHomeRental/internal/repository"
	"github.com/amartin3659/VacationHomeRental/internal/repository/dbrepo"
)

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// Repo the repository used by the handlers
var Repo *Repository

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the handler for home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home-page.html", &models.TemplateData{})
}

// About is the handler for the about page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "about-page.html", &models.TemplateData{})
}

// Contact is the handler for the contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact-page.html", &models.TemplateData{})
}

// Eremite is the handler for the eremite page
func (m *Repository) Eremite(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "eremite-page.html", &models.TemplateData{})
}

// Couple is the handler for the couple page
func (m *Repository) Couple(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "couple-page.html", &models.TemplateData{})
}

// Family is the handler for the family page
func (m *Repository) Family(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "family-page.html", &models.TemplateData{})
}

// Reservation is the handler for the reservation page
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "check-availability-page.html", &models.TemplateData{})
}

// PostReservation is the handler for the reservation page and POST requests
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("startingDate")
	end := r.Form.Get("endingDate")

  layout := "2006-01-02"

  startDate, err := time.Parse(layout, start)
  if err != nil {
    helpers.ServerError(w, err)
    return
  }

  endDate, err := time.Parse(layout, end)
  if err != nil {
    helpers.ServerError(w, err)
    return
  }

  bungalows, err := m.DB.SearchAvailabilityByDatesForAllBungalows(startDate, endDate)
  if err != nil {
    helpers.ServerError(w, err)
    return
  }

  if len(bungalows) == 0 {
    m.App.Session.Put(r.Context(), "error", ":( No holiday home is available at that time.")
    http.Redirect(w, r, "/reservation", http.StatusSeeOther)
    return
  }

  data := make(map[string]interface{})
  data["bungalows"] = bungalows

  res := models.Reservation{
    StartDate: startDate,
    EndDate: endDate,
  }

  m.App.Session.Put(r.Context(), "reservation", res)
  
  render.Template(w, r, "choose-bungalow-page.html", &models.TemplateData{
    Data: data,
  })
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// ReservationJSON is the handler for reservation-json and returns JSON
func (m *Repository) ReservationJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK:      true,
		Message: "It's available",
	}

	output, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}

// MakeReservation is the handler for the make-reservation page
func (m *Repository) MakeReservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation

	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.Template(w, r, "make-reservation-page.html", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostMakeReservation is the POST request handler for the reservation form
func (m *Repository) PostMakeReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")

	layout := "2006-01-02"

	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	bungalowID, err := strconv.Atoi(r.Form.Get("bungalow_id"))
	if err != nil {
		helpers.ServerError(w, err)
	}

	reservation := models.Reservation{
		FullName:   r.Form.Get("full_name"),
		Email:      r.Form.Get("email"),
		Phone:      r.Form.Get("phone"),
		StartDate:  startDate,
		EndDate:    endDate,
		BungalowID: bungalowID,
	}

	form := forms.New(r.PostForm)

	form.Required("full_name", "email")
	form.MinLength("full_name", 2)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.Template(w, r, "make-reservation-page.html", &models.TemplateData{
			Form: form,
			Data: data,
		})

		return
	}

  newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

  restriction := models.BungalowRestriction{
    StartDate: startDate,
    EndDate: endDate,
    BungalowID: bungalowID,
    ReservationID: newReservationID,
    RestrictionID: 1,
  }

  err = m.DB.InsertBungalowRestriction(restriction)
  if err != nil {
    helpers.ServerError(w, err)
    return
  }

	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-overview", http.StatusSeeOther)
}

// ReservationOveriew displays the reservation summary page
func (m *Repository) ReservationOverview(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("Could not get item from session")
		m.App.Session.Put(r.Context(), "error", "No reservation data in this session available.")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Template(w, r, "reservation-overview-page.html", &models.TemplateData{
		Data: data,
	})
}

// ChooseBungalow displays list of available bungalows and lets the user choose a bungalow
func (m *Repository) ChooseBungalow(w http.ResponseWriter, r *http.Request) {
  
  exploded := strings.Split(r.RequestURI, "/")
  bungalowID, err := strconv.Atoi(exploded[2])
  if err != nil {
    m.App.Session.Put(r.Context(), "error", "Missing parameter from URL")
    http.Redirect(w, r, "/", http.StatusSeeOther)
    return
  }

  res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
  if !ok {
    m.App.Session.Put(r.Context(), "error", "Cannot get reservation back from session")
    http.Redirect(w, r, "/", http.StatusSeeOther)
    return
  }
  
  res.BungalowID = bungalowID

  m.App.Session.Put(r.Context(), "reservation", res)

  http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}
