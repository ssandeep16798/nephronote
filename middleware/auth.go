package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("ssandeep98") // Replace with your actual secret key

type Claims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

func ValidateTokenAndGetUserID(tokenString string) (int, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return 0, errors.New("invalid token signature")
		}
		return 0, errors.New("could not parse token")
	}
	if !token.Valid {
		return 0, errors.New("invalid token")
	}

	return claims.UserID, nil
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("AuthMiddleware: Request received")
		if r.URL.Path == "/registeruser/" || r.URL.Path == "/login/" {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Println("AuthMiddleware: Missing Authorization header")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			log.Println("AuthMiddleware: Malformed Authorization header")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID, err := ValidateTokenAndGetUserID(tokenString)
		if err != nil {
			log.Println("AuthMiddleware: Invalid token:", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		log.Println("AuthMiddleware: Valid token, UserID:", userID)
		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
