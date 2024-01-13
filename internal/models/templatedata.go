package models

import "github.com/amartin3659/VacationHomeRental/internal/forms"

// TemplateData holds any kind of data sent from handlers to templates
type TemplateData struct {
	StringMap       map[string]string
	IntMap          map[int]int
	FloatMap        map[float64]float64
	Data            map[string]interface{}
	CSRFToken       string
	Success         string
	Warning         string
	Error           string
	Form            *forms.Form
	IsAuthenticated int
}
