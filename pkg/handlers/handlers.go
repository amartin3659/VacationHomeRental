package handlers

import (
	"net/http"

	"github.com/amartin3659/VacationHomeRental/pkg/models"
	"github.com/amartin3659/VacationHomeRental/pkg/config"
	"github.com/amartin3659/VacationHomeRental/pkg/render"
)

type Repository struct {
  App *config.AppConfig
}

// Repo the repository used by the handlers
var Repo *Repository

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
  return &Repository{
    App: a,
  }
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
  Repo = r
}

// Home is the handler for home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
  
  remoteIP := r.RemoteAddr
  m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

  render.RenderTemplate(w, "home-page.html", &models.TemplateData{})
}

// About is the handler for the about page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
  
  sidekickMap := make(map[string]string)
  sidekickMap["morty"] = "Ooh, wee!"

  remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
  sidekickMap["remote_ip"] = remoteIP

  render.RenderTemplate(w, "about-page.html", &models.TemplateData{
    StringMap: sidekickMap,
  })
}
