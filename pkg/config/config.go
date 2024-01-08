package config

import (
	"html/template"

	"github.com/alexedwards/scs/v2"
)

// AppConfig is a struct holding this application's configuration
type AppConfig struct {
	TemplateCache map[string]*template.Template
	UseCache      bool
	InProduction  bool
	Session       *scs.SessionManager
}
