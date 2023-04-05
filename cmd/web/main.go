package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/toothsy/bookings-app/internal/config"
	"github.com/toothsy/bookings-app/internal/driver"
	"github.com/toothsy/bookings-app/internal/handlers"
	"github.com/toothsy/bookings-app/internal/helpers"
	"github.com/toothsy/bookings-app/internal/models"
	"github.com/toothsy/bookings-app/internal/render"
	dbrepo "github.com/toothsy/bookings-app/internal/repository/DBrepo"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// main is the main function
func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	fmt.Printf("Staring application on http://localhost%s", portNumber)
	fmt.Println("")
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func run() (*driver.DB, error) {
	// what am I going to put in the session
	gob.Register(models.Reservation{})

	// change this to true when in production
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// set up the session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session
	// connecting to DB
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=postgres password=123")

	if err != nil {
		log.Fatal("cannot connec to DB")
	}

	dbrepo.NewPostgresConnection(db.SQL, &app)
	log.Println("connecting to DB")
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)

	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
