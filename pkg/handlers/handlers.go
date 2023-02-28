package handlers

import (
	"github/toothsy/bookings/pkg/config"
	"github/toothsy/bookings/pkg/models"
	"github/toothsy/bookings/pkg/renderers"
	"net/http"
)

type Repository struct {
	App *config.AppConfig
}

var Repo *Repository

// NewRepository return  a new Repository object with app config initilized

func NewRepository(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandler initialises the local repo var so that I can use it in other packages
func NewHandler(r *Repository) {
	Repo = r
}

// HomeHandler this function handles all traffic to "/"
func (m *Repository) HomeHandler(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)
	renderers.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

// AboutHandler this function handles all traffic to "/about"
func (m *Repository) AboutHandler(w http.ResponseWriter, r *http.Request) {
	StrMap := map[string]string{}
	StrMap["test"] = "hello there"
	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	StrMap["remote_ip"] = remoteIP

	// m.App.
	renderers.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{
		StringMap: StrMap,
	})

}
