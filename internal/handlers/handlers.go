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

// NewTestRepo creates a new repository for basic tests
func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestingRepo(a),
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
		m.App.Session.Put(r.Context(), "error", "Can't get data from form")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	endDate, err := time.Parse(layout, end)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't get data from form")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	bungalows, err := m.DB.SearchAvailabilityByDatesForAllBungalows(startDate, endDate)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't get data from database")
		http.Redirect(w, r, "/", http.StatusSeeOther)
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
		EndDate:   endDate,
	}

	m.App.Session.Put(r.Context(), "reservation", res)

	render.Template(w, r, "choose-bungalow-page.html", &models.TemplateData{
		Data: data,
	})
}

type jsonResponse struct {
	OK         bool   `json:"ok"`
	Message    string `json:"message"`
	BungalowID string `json:"bungalow_id"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
}

// ReservationJSON is the handler for reservation-json and returns JSON
func (m *Repository) ReservationJSON(w http.ResponseWriter, r *http.Request) {

	// parse request body
	err := r.ParseForm()
	if err != nil {
		// can't parse form, so return appropriate json
		resp := jsonResponse{
			OK:      false,
			Message: "Internal server error",
		}

		output, _ := json.MarshalIndent(resp, "", "    ")

		w.Header().Set("Content-Type", "application/json")
		w.Write(output)
		return
	}

	bungalowID, err := strconv.Atoi(r.Form.Get("bungalow_id"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't get data from form")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	layout := "2006-01-02"

	startDate, err := time.Parse(layout, sd)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't get data from form")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't get data from form")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	available, err := m.DB.SearchAvailabilityByDatesByBungalowID(startDate, endDate, bungalowID)
	if err != nil {
		helpers.ServerError(w, err)
		resp := jsonResponse{
			OK:      false,
			Message: "Error querying database",
		}

		output, err := json.MarshalIndent(resp, "", "    ")
		if err != nil {
			m.App.Session.Put(r.Context(), "error", "Can't provide valid json")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(output)

		return
	}

	resp := jsonResponse{
		OK:         available,
		Message:    "",
		StartDate:  sd,
		EndDate:    ed,
		BungalowID: strconv.Itoa(bungalowID),
	}

	output, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		helpers.ServerError(w, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}

// MakeReservation is the handler for the make-reservation page
func (m *Repository) MakeReservation(w http.ResponseWriter, r *http.Request) {

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		// helpers.ServerError(w, errors.New("Cannot get reservation back from session"))
		// return
		m.App.Session.Put(r.Context(), "error", "Cannot get reservation back from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	bungalow, err := m.DB.GetBungalowByID(res.BungalowID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Cannot find bungalow!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.Bungalow.BungalowName = bungalow.BungalowName

	m.App.Session.Put(r.Context(), "reservation", res)

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(w, r, "make-reservation-page.html", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

// PostMakeReservation is the POST request handler for the reservation form
func (m *Repository) PostMakeReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't parse form")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	//	sd := r.Form.Get("start_date")
	//	ed := r.Form.Get("end_date")
	//
	//	layout := "2006-01-02"
	//
	//	startDate, err := time.Parse(layout, sd)
	//	if err != nil {
	//		m.App.Session.Put(r.Context(), "error", "Can't parse arrival date")
	//		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	//		return
	//	}
	//
	//	endDate, err := time.Parse(layout, ed)
	//	if err != nil {
	//		m.App.Session.Put(r.Context(), "error", "Can't parse departure date")
	//		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	//		return
	//	}
	//
	//	bungalowID, err := strconv.Atoi(r.Form.Get("bungalow_id"))
	//	if err != nil {
	//		m.App.Session.Put(r.Context(), "error", "Can't parse bungalow id")
	//		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	//		return
	//	}

	// getting the reservation data so far from session
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "Cannot get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	reservation := models.Reservation{
		FullName:   r.Form.Get("full_name"),
		Email:      r.Form.Get("email"),
		Phone:      r.Form.Get("phone"),
		StartDate:  res.StartDate,
		EndDate:    res.EndDate,
		BungalowID: res.BungalowID,
		Bungalow:   models.Bungalow{BungalowName: res.Bungalow.BungalowName},
	}

	// validate form data
	form := forms.New(r.PostForm)

	form.Required("full_name", "email")
	form.MinLength("full_name", 2)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		// if new rendering of page needed store already collected
		// (and maybe in session store) Dates as string in stringMap
		// to use when template is rendered again
		sd := res.StartDate.Format("2006-01-02")
		ed := res.EndDate.Format("2006-01-02")

		stringMap := make(map[string]string)
		stringMap["start_date"] = sd
		stringMap["end_date"] = ed

		// write new reservation in session as far collected (so far)
		m.App.Session.Put(r.Context(), "reservation", reservation)

		render.Template(w, r, "make-reservation-page.html", &models.TemplateData{
			Form:      form,
			Data:      data,
			StringMap: stringMap,
		})

		return
	}

	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't write reservation to database")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	restriction := models.BungalowRestriction{
		StartDate:     res.StartDate,
		EndDate:       res.EndDate,
		BungalowID:    res.BungalowID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}

	err = m.DB.InsertBungalowRestriction(restriction)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't reserve bungalow in database")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-overview", http.StatusSeeOther)
}

// ReservationOveriew displays the reservation summary page
func (m *Repository) ReservationOverview(w http.ResponseWriter, r *http.Request) {

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("Could not get item from session")
		m.App.Session.Put(r.Context(), "error", "No reservation data in this session available.")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(), "reservation")

	bungalow, err := m.DB.GetBungalowByID(res.BungalowID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Cannot find bungalow")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	res.Bungalow.BungalowName = bungalow.BungalowName

	data := make(map[string]interface{})
	data["reservation"] = res

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	render.Template(w, r, "reservation-overview-page.html", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
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

func (m *Repository) BookBungalow(w http.ResponseWriter, r *http.Request) {
	bungalowID, _ := strconv.Atoi(r.URL.Query().Get("id"))
	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")

	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	var res models.Reservation

	bungalow, err := m.DB.GetBungalowByID(bungalowID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Cannot find bungalow!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	res.Bungalow.BungalowName = bungalow.BungalowName
	res.BungalowID = bungalowID
	res.StartDate = startDate
	res.EndDate = endDate

	m.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}
