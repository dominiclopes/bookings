package handlers

import (
	"net/http"

	"github.com/dominiclopes/bookings/pkg/config"
	"github.com/dominiclopes/bookings/pkg/models"
	"github.com/dominiclopes/bookings/pkg/render"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
}

func NewRepository(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	// Save the IP address of the request
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again."

	// Use the remote_ip stored in the session and store it in the stringMap
	remote_ip := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remote_ip

	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}
