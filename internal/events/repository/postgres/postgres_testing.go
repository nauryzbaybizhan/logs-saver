package events_postgres_repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"google.golang.org/appengine/log"
)

func MockDb() (*sqlx.DB, sqlmock.Sqlmock) {
	sqlDb, mock, err := sqlmock.New()
	if err != nil {
		log.Errorf(nil, "error while creating mock db: %s", err)
	}
	sqlxDb := sqlx.NewDb(sqlDb, "sqlmock")
	return sqlxDb, mock
}
