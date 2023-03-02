package handlers

import (
	"encoding/json"
	"fmt"
	"github/toothsy/bookings/internal/config"
	"github/toothsy/bookings/internal/models"
	"github/toothsy/bookings/internal/renderers"
	"log"
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
	renderers.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})
}

// AboutHandler this function handles all traffic to "/about"
func (m *Repository) AboutHandler(w http.ResponseWriter, r *http.Request) {
	StrMap := map[string]string{}
	StrMap["test"] = "hello there"
	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	StrMap["remote_ip"] = remoteIP

	// m.App.
	renderers.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{
		StringMap: StrMap,
	})

}

// GodHandler this function handles all traffic to "/god-rooms"
func (m *Repository) GodHandler(w http.ResponseWriter, r *http.Request) {
	// m.App.
	renderers.RenderTemplate(w, r, "god.page.tmpl", &models.TemplateData{})

}

// EmpHandler this function handles all traffic to "/emp-rooms"
func (m *Repository) EmpHandler(w http.ResponseWriter, r *http.Request) {
	renderers.RenderTemplate(w, r, "emp.page.tmpl", &models.TemplateData{})

}

// KingHandler this function handles all traffic to "/king-rooms"
func (m *Repository) KingHandler(w http.ResponseWriter, r *http.Request) {
	renderers.RenderTemplate(w, r, "king.page.tmpl", &models.TemplateData{})

}

// SaintHandler this function handles all traffic to "/saint-rooms"
func (m *Repository) SaintHandler(w http.ResponseWriter, r *http.Request) {
	renderers.RenderTemplate(w, r, "saint.page.tmpl", &models.TemplateData{})

}

// Reservation this function handles all traffic to "/make-reservation"
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	renderers.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{})

}

// Availability this function handles all traffic to "/search-availability"
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	renderers.RenderTemplate(w, r, "search-availability.page.tmpl", &models.TemplateData{})

}

// PostAvailability this function handles all post traffic to "/search-availability"
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	w.Write([]byte(fmt.Sprintf("<h1>start and end date are %s and %s</h1>", start, end)))

}

type jsonResponse struct {
	Ok  bool   `json:"ok"`
	Msg string `json:"msg"`
}

// AvailabilityJSON this function responds with JSON based on availability "/search-availability"
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		Ok:  true,
		Msg: "hello form routes",
	}
	out, err := json.MarshalIndent(resp, "", "\t")
	if err != nil {
		log.Fatal("coudnt marshal JSON")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)

}

// Contact this function handles all traffic to "/contact"
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	renderers.RenderTemplate(w, r, "contact.page.tmpl", &models.TemplateData{})

}
