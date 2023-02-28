package main

import (
	"github/toothsy/bookings/pkg/config"
	"github/toothsy/bookings/pkg/handlers"
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

	fileserver := http.FileServer(http.Dir("/static/"))
	mux.Handle("/static", http.StripPrefix("/static", fileserver))

	return mux
}
