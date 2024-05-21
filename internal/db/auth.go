package db

import (
	"database/sql"
	"fmt"
	"nephronote/internal/models"

	"golang.org/x/crypto/bcrypt"
)

// AuthenticateUser verifies the user's credentials and returns an authenticated user struct
func AuthenticateUser(username, password string) (*models.AuthenticatedUser, error) {
	dbConn, err := getDBConnection()
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %v", err)
	}
	defer dbConn.Close()

	var storedHashedPassword string
	var userID int
	var email string

	query := `SELECT id, email, password FROM users WHERE username=$1`
	err = dbConn.QueryRow(query, username).Scan(&userID, &email, &storedHashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invalid username or password")
		}
		return nil, fmt.Errorf("error fetching user data: %v", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid username or password")
	}

	user := &models.AuthenticatedUser{
		ID:       userID,
		UserName: username,
		Email:    email,
	}

	return user, nil
}
