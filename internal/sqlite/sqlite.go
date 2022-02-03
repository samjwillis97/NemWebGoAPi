package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3" // Required for SQLite
	log "github.com/sirupsen/logrus"
)

// New, returns a new instance of sql db
func New(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		log.Fatalln("Cannot connect to DB: ", err)
	} else {
		log.Infof("Successfully connected to DB at %s", filepath)
	}

	return db
}
