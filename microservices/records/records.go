package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Records struct {
	ID					int 		`json:"ID"`
	AccountID			int 		`json:"AccountID"`
	ContactRole			string 		`json:"ContactRole"`
	StudentCount		int 		`json:"StudentCount"`	
	AcadYear			string 		`json:"AcadYear"`
	Title				string 		`json:"Title"`
	CompanyName			string 		`json:"CompanyName"`
	CompanyPOC			string 		`json:"CompanyPOC"`
	Description			string 		`json:"Description"`
	CreationDate        string 		`json:"CreationDate"`
	IsDeleted			bool 		`json:"IsDeleted"`
}

var connectionString	string

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
		AllowedOrigins: []string{"http://127.0.0.1:5502"}, // Your frontend origin
        AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders: []string{"Content-Type"},
    })
		
	cmd := flag.String("sql", "", "")
	flag.Parse()
	connectionString = string(*cmd)
	
    handler := corsHandler.Handler(r)
	http.ListenAndServe(":8001", handler)
}

func CreateRecord(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	log.Println("Entering endpoint to add new capstone entry record")

	var newRecords Records
	err := json.NewDecoder(r.Body).Decode(&newRecords)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	location, _ := time.LoadLocation("Asia/Singapore")
	newRecords.CreationDate = time.Now().In(location).Format("2006-01-02 15:04:05")
	newRecords.IsDeleted = false

	db, _ := sql.Open("mysql", connectionString)
	defer db.Close()

	result, err := db.Exec(
		`INSERT INTO tsao_records (AccountID, ContactRole, StudentCount, AcadYear, Title, CompanyName, CompanyPOC, Description, CreationDate, IsDeleted)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`,
		newRecords.AccountID, newRecords.ContactRole, newRecords.StudentCount, newRecords.AcadYear,
		newRecords.Title, newRecords.CompanyName, newRecords.CompanyPOC, newRecords.Description, newRecords.CreationDate, newRecords.IsDeleted)
	if err == nil  {
		recordID, _ := result.LastInsertId()
		newRecords.ID = int(recordID)
		
		newRecordsJson, _ := json.Marshal(newRecords)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		w.Write(newRecordsJson)
	} else {
		log.Println(err)
		http.Error(w, "Error creating new capstone record", http.StatusInternalServerError)
		return
	}
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