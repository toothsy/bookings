package dbrepo

import (
	"database/sql"

	"github.com/toothsy/bookings-app/internal/config"
	"github.com/toothsy/bookings-app/internal/repository"
)

type postrgesDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresConnection(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postrgesDBRepo{
		App: a,
		DB:  conn,
	}
}
