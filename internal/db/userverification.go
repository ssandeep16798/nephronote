package db

import (
	"database/sql"
	"errors"
	"fmt"
	"nephronote/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func CheckIfEmailExists(email string, username string) bool {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	defer db.Close()

	var count int
	checkEmailQuery := `SELECT COUNT(*) FROM users WHERE email=$1 OR username=$2`
	err = db.QueryRow(checkEmailQuery, email, username).Scan(&count)
	if err != nil {
		CheckError(err)
	}

	fmt.Println(count)
	return count > 0
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}

}

func UserDB(username, password string) (*models.UserRegistrationForm, error) {
	db, err := getDBConnection()
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %v", err)
	}
	defer db.Close()

	var user models.UserRegistrationForm
	var storedHashedPassword string

	fetchPasswordStmt := `SELECT "password", "id", "firstname", "lastname", "email", "username", "gender", "phone_number", "created_datetime" 
                          FROM users 
                          WHERE "username"=$1`
	err = db.QueryRow(fetchPasswordStmt, username).Scan(&storedHashedPassword, &user.ID, &user.FirstName, &user.LastName, &user.Email, &user.UserName, &user.Gender, &user.Mobile, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid username or password")
		}
		return nil, fmt.Errorf("error fetching user data: %v", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(password))
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	return &user, nil
}
