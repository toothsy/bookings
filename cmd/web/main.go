package main

// I used fresh to keep running it
import (
	"encoding/gob"
	"fmt"
	"github/toothsy/bookings/internal/config"
	"github/toothsy/bookings/internal/handlers"
	"github/toothsy/bookings/internal/models"
	"github/toothsy/bookings/internal/renderers"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":6969"

var app config.AppConfig

var session *scs.SessionManager

func main() {
	app.UseSecure = false
	app.UseCache = false
	gob.Register(&models.Reservation{})
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.UseSecure

	app.Session = session
	templateCache, err := renderers.CreateTemplateCache()
	if err != nil {
		log.Fatal("coudlnt create cache ", err)
	}

	app.TemplateCache = templateCache
	renderers.SetConfig(&app)
	repoReference := handlers.NewRepository(&app)
	handlers.NewHandler(repoReference)

	fmt.Printf("the server is up and running at http://localhost%s\n", portNumber)
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal("could not start the server ", err)
}
