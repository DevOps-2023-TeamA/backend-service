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
    api := r.PathPrefix("/api/users").Subrouter()
    api.HandleFunc("", CreateUser).Methods("POST")
    api.HandleFunc("", ReadUsers).Methods("GET")
    api.HandleFunc("/{id}", UpdateUser).Methods("PUT")
    api.HandleFunc("/{id}", DeleteUser).Methods("DELETE")
    http.Handle("/", r)

	fmt.Println("Accounts microservice running on http://localhost:8002/api/users")
	
	// CORS configuration
    corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://127.0.0.1:5500"}, // Your frontend origin
        AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders: []string{"Content-Type"},
    })
	
	
    handler := corsHandler.Handler(r)
	http.ListenAndServe(":8002", handler)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	log.Println("Entering endpoint to add new account")
}

func ReadUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	log.Println("Entering endpoint to query all accounts")
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	log.Println("Entering endpoint to update an account's information")
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	log.Println("Entering endpoint to (soft) delete an account")
}