package db

import (
	"database/sql"
	"fmt"
	"nephronote/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(req models.UserRegistrationForm) error {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return fmt.Errorf("error opening database connection: %v", err)
	}
	defer db.Close()

	// Hash the password before storing it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %v", err)
	}

	insertDynStmt := `INSERT INTO users("firstname", "lastname", "email", "gender", "phone_number", "username", "password", "created_datetime") 
                      VALUES($1, $2, $3, $4, $5, $6, $7, NOW())`
	_, err = db.Exec(insertDynStmt, req.FirstName, req.LastName, req.Email, req.Gender, req.Mobile, req.UserName, hashedPassword)
	if err != nil {
		return fmt.Errorf("error inserting user data: %v", err)
	}

	return nil
}
