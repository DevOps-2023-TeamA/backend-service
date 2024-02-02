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

var connectionString	string

func main() {
	r := mux.NewRouter()
    api := r.PathPrefix("/api/records").Subrouter()
    api.HandleFunc("", CreateRecord).Methods("POST")
    api.HandleFunc("", ReadRecords).Methods("GET")
    api.HandleFunc("/{id}", ReadRecord).Methods("GET")
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

	acadYear := r.URL.Query().Get("ay")
	title := r.URL.Query().Get("title")

	db, _ := sql.Open("mysql", connectionString)
	defer db.Close()
	var result *sql.Rows
	var err error

	if acadYear == "" && title == "" {
		fmt.Println("No queries")
    	result, err = db.Query(`SELECT * FROM tsao_records WHERE IsDeleted=false`)
	} else if acadYear == "" {
		fmt.Println("No acadYear, have title")
		result, err = db.Query(`SELECT * FROM tsao_records WHERE Title LIKE ? AND IsDeleted=false`, "%"+title+"%")
		} else if title == "" {
		fmt.Println("No title, have acadYear")
    	result, err = db.Query(`SELECT * FROM tsao_records WHERE AcadYear=? AND IsDeleted=false`, acadYear)
	} else {
		fmt.Println("Have acadYear and title")
		result, err = db.Query(`SELECT * FROM tsao_records WHERE AcadYear=? AND Title LIKE ? AND IsDeleted=false`, acadYear, "%"+title+"%")
	}

	var records []Records
	for result.Next() {
		var record Records
		_ = result.Scan(
			&record.ID, &record.AccountID, &record.ContactRole,
			&record.StudentCount, &record.AcadYear, &record.Title,
			&record.CompanyName, &record.CompanyPOC, &record.Description, &record.CreationDate, &record.IsDeleted)
		records = append(records, record)
	}

	if err == nil  {
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(records)
	} else if err := result.Err(); err != sql.ErrNoRows {
		log.Println(err)
		http.Error(w, "Error iterating over rows", http.StatusInternalServerError)
		return
	}
}

func ReadRecord(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	log.Println("Entering endpoint to query single capstone entries")

	recordID := mux.Vars(r)["id"]	

	db, _ := sql.Open("mysql", connectionString)
	defer db.Close()
	
	var record Records
    err := db.QueryRow("SELECT * FROM tsao_records WHERE ID=?", recordID).Scan(
		&record.ID, &record.AccountID, &record.ContactRole,
			&record.StudentCount, &record.AcadYear, &record.Title,
			&record.CompanyName, &record.CompanyPOC, &record.Description, &record.CreationDate, &record.IsDeleted)

	if err == nil  {
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(record)
	} else if err == sql.ErrNoRows{
		http.Error(w, "Record does not exist", http.StatusNotFound)
		return
	} else {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func UpdateRecord(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	log.Println("Entering endpoint to update a capstone entry record")

	recordID := mux.Vars(r)["id"]	
	
	var modifiedRecord Records
	err := json.NewDecoder(r.Body).Decode(&modifiedRecord)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	
	db, _ := sql.Open("mysql", connectionString)
	defer db.Close()

	result, err := db.Exec(
		`UPDATE tsao_records
		SET
			AccountID=?, ContactRole=?, StudentCount=?, AcadYear=?, 
			Title=?, CompanyName=?, CompanyPOC=?, Description=?
		WHERE ID=?;`,
		modifiedRecord.AccountID, modifiedRecord.ContactRole, modifiedRecord.StudentCount, modifiedRecord.AcadYear, 
		modifiedRecord.Title, modifiedRecord.CompanyName, modifiedRecord.CompanyPOC, modifiedRecord.Description, recordID)
	rowsAffected, _ := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		log.Println(err)
		http.Error(w, "Record ID does not exist", http.StatusNotFound)
		return
	}  else {
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintf(w, "Record information modified for ID: %s\n", recordID)
	}
}

func DeleteRecord(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	log.Println("Entering endpoint to (soft) delete a capstone entry record")

	recordID := mux.Vars(r)["id"]	
	
	db, _ := sql.Open("mysql", connectionString)
	defer db.Close()
	result, err := db.Exec(
		`UPDATE tsao_records SET IsDeleted=true WHERE ID=? AND IsDeleted=false;`, recordID)
	rowsAffected, _ := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		log.Println(err)
		http.Error(w, "Unable to delete user", http.StatusNotFound)
		return
	} else {
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintf(w, "User deleted for ID: %s\n", recordID)
	}
}