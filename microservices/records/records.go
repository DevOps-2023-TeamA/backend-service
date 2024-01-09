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
    api := r.PathPrefix("/api/records").Subrouter()
    api.HandleFunc("", CreateRecord).Methods("POST")
    api.HandleFunc("", ReadRecords).Methods("GET")
    api.HandleFunc("/{id}", UpdateRecord).Methods("PUT")
    api.HandleFunc("/{id}", DeleteRecord).Methods("DELETE")
    http.Handle("/", r)

	fmt.Println("Capstone Entries microservice running on http://localhost:8001/api/records")
	
	// CORS configuration
    corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://127.0.0.1:5500"}, // Your frontend origin
        AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders: []string{"Content-Type"},
    })
	
	
    handler := corsHandler.Handler(r)
	http.ListenAndServe(":8001", handler)
}

func CreateRecord(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	log.Println("Entering endpoint to add new capstone entry record")
}

func ReadRecords(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	log.Println("Entering endpoint to query all capstone entries")
}

func UpdateRecord(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	log.Println("Entering endpoint to update a capstone entry record")
}

func DeleteRecord(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	log.Println("Entering endpoint to (soft) delete a capstone entry record")
}