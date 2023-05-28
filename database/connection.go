package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB

// Connect initializes the database connection
func Connect(dbPath string) error {
	var err error

	// Open the database connection
	Db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	// Ping the database to verify the connection
	err = Db.Ping()
	if err != nil {
		return err
	}

	log.Println("Connected to the SQLite database")

	return nil
}

// Close disconnects from the database
func Close() error {
	if Db != nil {
		err := Db.Close()
		if err != nil {
			return err
		}
		log.Println("Disconnected from the SQLite database")
	}

	return nil
}

