package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	r := mux.NewRouter()
    api := r.PathPrefix("/api/auth").Subrouter()
    api.HandleFunc("", Login).Methods("POST")
    http.Handle("/", r)

	fmt.Println("Auth microservice running on http://localhost:8000/api/auth")
	
	// CORS configuration
    corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://127.0.0.1:5500"}, // Your frontend origin
        AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders: []string{"Content-Type"},
    })
	
	
    handler := corsHandler.Handler(r)
	http.ListenAndServe(":8000", handler)
}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	log.Println("Entering endpoint to add new user account")
}