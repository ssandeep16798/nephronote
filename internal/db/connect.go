package db

import (
	"database/sql"
	"fmt"
)

const (
	dbname   = "users"
	host     = "localhost"
	port     = 5432
	password = "sandeep"
	user     = "sandeep"
)

func getDBConnection() (*sql.DB, error) {
	host := "localhost"   // replace with your DB host
	port := "5432"        // replace with your DB port
	user := "sandeep"     // replace with your DB user
	password := "sandeep" // replace with your DB password
	dbname := "users"     // replace with your DB name

	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %v", err)
	}

	// Test the database connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %v", err)
	}

	return db, nil
}
