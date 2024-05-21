package main

import (
	"log"
	"nephronote/internal/handlers"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter() // Add parentheses to call the function

	r.HandleFunc("/registeruser/", handlers.Register).Methods("POST")
	r.HandleFunc("/login/", handlers.Login).Methods("GET")
	r.HandleFunc("/pre_dialysis/", handlers.PreDialysisHandler).Methods("POST")
	r.HandleFunc("/post_dialysis/", handlers.PostDialysisHandler).Methods("POST")
	/*r.HandleFunc("/healthdata/", func(w http.ResponseWriter, r *http.Request) {
		// Call the TrackerAPI handler function and pass the database connection
		handler := handlers.TrackerAPI(db)
		// Call the handler function, which writes the response directly to the ResponseWriter
		handler(w, r)
	}).Methods("POST")*/
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	log.Println("Initializing server at localhost", server.Addr, "....")
	if err := server.ListenAndServe(); err != nil {
		log.Println("Error starting server:", err)
	}
}
