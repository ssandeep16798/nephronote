package main

import (
	"log"
	"nephronote/internal/handlers"
	"nephronote/middleware"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// Public routes
	r.HandleFunc("/registeruser/", handlers.Register).Methods("POST")
	r.HandleFunc("/login/", handlers.Login).Methods("GET")

	// Protected routes
	r.PathPrefix("/api/").Subrouter()
	r.Use(middleware.AuthMiddleware)
	r.HandleFunc("/pre_dialysis/", handlers.PreDialysisHandler).Methods("POST")
	r.HandleFunc("/post_dialysis/", handlers.PostDialysisHandler).Methods("POST")

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	log.Println("Initializing server at localhost", server.Addr, "....")
	if err := server.ListenAndServe(); err != nil {
		log.Println("Error starting server:", err)
	}
}
