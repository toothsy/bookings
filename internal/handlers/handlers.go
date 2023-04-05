package handlers

// auto increment tips
//https://stackoverflow.com/questions/5342440/reset-auto-increment-counter-in-postgres
import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/toothsy/bookings-app/internal/config"
	"github.com/toothsy/bookings-app/internal/driver"
	"github.com/toothsy/bookings-app/internal/forms"
	"github.com/toothsy/bookings-app/internal/helpers"
	"github.com/toothsy/bookings-app/internal/models"
	"github.com/toothsy/bookings-app/internal/render"
	"github.com/toothsy/bookings-app/internal/repository"
	dbrepo "github.com/toothsy/bookings-app/internal/repository/DBrepo"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresConnection(db.SQL, a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the handler for the home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About is the handler for the about page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

// Reservation renders the make a reservation page and displays form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	emptyReservation.Phone = "123-123-1234"
	emptyReservation.Email = "adf@dot.com"
	data["reservation"] = emptyReservation

	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostReservation handles the posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")
	timeFormat := "2006-01-02"
	startDate, err := time.Parse(timeFormat, sd)

	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(timeFormat, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return

	}
	roomId, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return

	}
	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomId:    roomId,
	}

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	// form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	resId, err := m.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	restr := models.RoomRestriction{
		RestrictionId: 1,
		ReservationId: resId,
		StartDate:     startDate,
		EndDate:       endDate,
		RoomId:        roomId,
	}

	err = m.DB.InsertRoomReservation(restr)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// GodHandler this function handles all traffic to "/god-rooms"
func (m *Repository) GodHandler(w http.ResponseWriter, r *http.Request) {
	// m.App.
	render.Template(w, r, "god.page.tmpl", &models.TemplateData{})

}

// EmpHandler this function handles all traffic to "/emp-rooms"
func (m *Repository) EmpHandler(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "emp.page.tmpl", &models.TemplateData{})

}

// KingHandler this function handles all traffic to "/king-rooms"
func (m *Repository) KingHandler(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "king.page.tmpl", &models.TemplateData{})

}

// SaintHandler this function handles all traffic to "/saint-rooms"
func (m *Repository) SaintHandler(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "saint.page.tmpl", &models.TemplateData{})

}

// Availability renders the search availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostAvailability handles post
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")
	timeFormat := "2006-01-02"
	startDate, err := time.Parse(timeFormat, start)

	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(timeFormat, end)
	if err != nil {
		helpers.ServerError(w, err)
		return

	}
	m.DB.SearchAvailability(startDate, endDate)
	w.Write([]byte(fmt.Sprintf("start date is %s and end is %s", start, end)))
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJSON handles request for availability and sends JSON response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK:      true,
		Message: "Available!",
	}

	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Contact renders the contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}

// ReservationSummary displays the res summary page
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("Can't get reservation from session")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
