package handlers

// auto increment tips
//https://stackoverflow.com/questions/5342440/reset-auto-increment-counter-in-postgres
import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
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
var timeFormat string = "2006-01-02"

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
	sessionReservationValue, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("Can't get reservation from session")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	sd := sessionReservationValue.StartDate.Format(timeFormat)
	ed := sessionReservationValue.EndDate.Format(timeFormat)
	var rt models.RoomType = models.RoomType(sessionReservationValue.RoomId)

	data := make(map[string]interface{})
	data["reservation"] = sessionReservationValue
	sm := make(map[string]string)
	sm["start_date"] = sd
	sm["end_date"] = ed
	sm["room_type"] = rt.String()
	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: sm,
	})
}

// PostReservation handles the posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	sessionReservationValue, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("Can't get reservation from session")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	// form.IsEmail("email")
	sessionReservationValue.FirstName = r.Form.Get("first_name")
	sessionReservationValue.LastName = r.Form.Get("last_name")
	sessionReservationValue.Email = r.Form.Get("email")
	sessionReservationValue.Phone = r.Form.Get("phone")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = sessionReservationValue
		render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	newResId, err := m.DB.InsertReservation(sessionReservationValue)
	fmt.Println("sessionReservationValue.Room.RoomName:", sessionReservationValue.Room.RoomName)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	restr := models.RoomRestriction{
		RestrictionId: 1,
		ReservationId: newResId,
		StartDate:     sessionReservationValue.StartDate,
		EndDate:       sessionReservationValue.EndDate,
		RoomId:        sessionReservationValue.RoomId,
	}

	err = m.DB.InsertRoomReservation(restr)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	m.App.Session.Put(r.Context(), "reservation", sessionReservationValue)
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
	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if len(rooms) == 0 {
		m.App.Session.Put(r.Context(), "error", "No Availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	fmt.Println("rooms is ", rooms)
	data["rooms"] = rooms
	sessionReservationValue := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}
	m.App.Session.Put(r.Context(), "reservation", sessionReservationValue)
	render.Template(w, r, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJSON handles request for availability and sends JSON response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")
	layout := "02-01-2006"
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
	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	isAvailable, err := m.DB.SearchRoomReservationByRoomID(startDate, endDate, roomID)
	fmt.Println("sd is ", startDate, "end")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	resp := jsonResponse{
		OK:      isAvailable,
		Message: "",
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
	sessionReservationValue, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("Can't get reservation from session")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	sd := sessionReservationValue.StartDate.Format("2006-01-02")
	ed := sessionReservationValue.EndDate.Format("2006-01-02")

	data := make(map[string]interface{})
	data["reservation"] = sessionReservationValue
	sm := make(map[string]string)
	sm["start_date"] = sd
	sm["end_date"] = ed
	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: sm,
	})
	m.App.Session.Remove(r.Context(), "reservation")
}

// ChooseRooms allows the user to choose the available rooms
func (m *Repository) ChooseRooms(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	sessionReservationValue, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("Can't get reservation from session")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	sessionReservationValue.RoomId = roomID
	sessionReservationValue.Room.RoomName = models.RoomType(roomID).String()

	m.App.Session.Put(r.Context(), "reservation", sessionReservationValue)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)

}
