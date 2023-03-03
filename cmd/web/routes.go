package main

import (
	"github/toothsy/bookings/internal/config"
	"github/toothsy/bookings/internal/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// tips for PAT routing
// routes defines the routes the app must take
// notice that the routes must be in desencding order or matching...
// if "/" is at the top, only that page will be served at all pages
// chi routing
// it doesnt care about the order, it'll serve nonethless
func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(NoSurfCSRFTokenCheck)
	mux.Use(SessionLoad)
	mux.Get("/", handlers.Repo.HomeHandler)
	mux.Get("/about", handlers.Repo.AboutHandler)
	mux.Get("/contact", handlers.Repo.Contact)
	mux.Get("/god-rooms", handlers.Repo.GodHandler)
	mux.Get("/emp-rooms", handlers.Repo.EmpHandler)
	mux.Get("/king-rooms", handlers.Repo.KingHandler)
	mux.Get("/saint-rooms", handlers.Repo.SaintHandler)
	mux.Get("/make-reservation", handlers.Repo.Reservation)
	mux.Post("/make-reservation", handlers.Repo.PostReservation)

	mux.Get("/reservation-summary", handlers.Repo.ReservationSummary)

	mux.Get("/search-availability", handlers.Repo.Availability)
	mux.Post("/search-availability", handlers.Repo.PostAvailability)
	mux.Post("/search-availability-json", handlers.Repo.AvailabilityJSON)

	fileserver := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileserver))
	// https://stackoverflow.com/questions/44662456/why-my-golang-template-is-not-picking-the-external-javascript-file
	// to properly understand what line 34 does, look it up
	return mux
}
