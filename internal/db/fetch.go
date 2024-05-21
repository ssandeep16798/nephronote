package db

import (
	"database/sql"
	"fmt"
	"nephronote/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func Fetchdata(username string, password string) (models.UserRegistrationForm, error) {
	var registerReq models.UserRegistrationForm

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return registerReq, fmt.Errorf("error opening database connection: %v", err)
	}
	defer db.Close()

	var storedHashedPassword string
	fetchPasswordStmt := `SELECT "password" FROM users WHERE "username"=$1`
	err = db.QueryRow(fetchPasswordStmt, username).Scan(&storedHashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return registerReq, fmt.Errorf("username not found")
		}
		return registerReq, fmt.Errorf("error fetching user password: %v", err)
	}

	// Check if the stored password is a valid bcrypt hash
	if len(storedHashedPassword) == 0 || storedHashedPassword[0] != '$' {
		return registerReq, fmt.Errorf("stored password is not a valid bcrypt hash")
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(password))
	if err != nil {
		return registerReq, fmt.Errorf("invalid password")
	}

	fetchUserStmt := `SELECT "id", "firstname", "lastname", "email", "username", "gender", "phone_number", "created_datetime" 
                      FROM users 
                      WHERE "username"=$1`
	err = db.QueryRow(fetchUserStmt, username).Scan(&registerReq.ID, &registerReq.FirstName, &registerReq.LastName, &registerReq.Email, &registerReq.UserName, &registerReq.Gender, &registerReq.Mobile, &registerReq.CreatedAt)
	if err != nil {
		return registerReq, fmt.Errorf("error fetching user data: %v", err)
	}

	return registerReq, nil
}
