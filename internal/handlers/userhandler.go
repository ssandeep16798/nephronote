package handlers

import (
	"encoding/json"
	"fmt"
	"nephronote/internal/db"
	"nephronote/internal/models"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func Register(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Inside Register New User")

	var userRegistration models.UserRegistrationForm
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&userRegistration); err != nil {
		fmt.Println("Error decoding request body:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if db.CheckIfEmailExists(userRegistration.Email, userRegistration.UserName) {
		fmt.Println("Email or username already exists")
		var loginResponse models.LoginResponse
		loginResponse.Status = false
		loginResponse.Msg = "Email or username already exists"
		json.NewEncoder(w).Encode(loginResponse)
		return
	}

	if err := db.RegisterUser(userRegistration); err != nil {
		fmt.Println("Error registering user:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response, err := db.Fetchdata(userRegistration.UserName, userRegistration.Password)
	if err != nil {
		fmt.Println("Error fetching data:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var loginResponseWithData models.LoginResponseWithAllData
	loginResponseWithData.Status = true
	loginResponseWithData.Msg = "User registered successfully"
	loginResponseWithData.Data = response

	fmt.Println("Successfully registered user:", response.Email)
	json.NewEncoder(w).Encode(loginResponseWithData)
}

var JwtKey = []byte("my_secret_key")

type Claims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

func Login(w http.ResponseWriter, req *http.Request) {
	var loginResponse models.LoginResponse

	fmt.Println("Inside User Login")
	decoder := json.NewDecoder(req.Body)
	fmt.Println("After decoding")
	var loginRequest models.LoginRequest
	if err := decoder.Decode(&loginRequest); err != nil {
		fmt.Println("Error decoding request body:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Println("after second decoding")

	if loginRequest.UserName == "" {
		loginResponse.Status = false
		loginResponse.Msg = "Username is required"
		json.NewEncoder(w).Encode(loginResponse)
		return
	}

	user, err := db.UserDB(loginRequest.UserName, loginRequest.Password)
	if err != nil {
		fmt.Println("Error querying user database:", err)
		loginResponse.Status = false
		loginResponse.Msg = "Invalid username or password"
		json.NewEncoder(w).Encode(loginResponse)
		return
	}
	fmt.Println("after checking conditions.")

	// Token generation logic
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		fmt.Println("Error signing the token:", err)
		loginResponse.Status = false
		loginResponse.Msg = "Internal server error"
		json.NewEncoder(w).Encode(loginResponse)
		return
	}

	fmt.Println("Generated Token:", tokenString) // Debugging: Print the generated token

	var loginResponseWithData models.LoginResponseWithData
	loginResponseWithData.Status = true
	loginResponseWithData.Msg = ""
	loginResponseWithData.Data = *user // Use the user data returned by UserDB
	loginResponseWithData.Token = tokenString
	json.NewEncoder(w).Encode(loginResponseWithData)
	fmt.Println("generated response")
}
