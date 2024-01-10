package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Accounts struct {
	ID       		int    	`json:"ID"`
	Name     		string 	`json:"Name"`
	Username     	string 	`json:"Username"`
	Password		string 	`json:"Password"`
	Role			string 	`json:"Role"`
	CreationDate	string	`json:"CreationDate"`
	IsApproved		bool 	`json:"IsApproved"`
	IsDeleted		bool 	`json:"IsDeleted"`
}

var connectionString	string

func main() {
	r := mux.NewRouter()
    api := r.PathPrefix("/api/auth").Subrouter()
    api.HandleFunc("", Login).Methods("POST")
    http.Handle("/", r)

	fmt.Println("Auth microservice running on http://localhost:8000/api/auth")
	
	// CORS configuration
    corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://127.0.0.1:8080"}, // Your frontend origin
        AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders: []string{"Content-Type"},
    })
	
	cmd := flag.String("sql", "", "")
	flag.Parse()
	connectionString = string(*cmd)
	
    handler := corsHandler.Handler(r)
	http.ListenAndServe(":8000", handler)
}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	log.Println("Entering endpoint to add new user account")
	
	var loginAcc Accounts
	err := json.NewDecoder(r.Body).Decode(&loginAcc)
	if err != nil {
        http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
        return
    }
	sha := sha256.New()
	sha.Write([]byte(loginAcc.Password))
	encyrptedPassword := hex.EncodeToString(sha.Sum(nil))

	db, _ := sql.Open("mysql", connectionString)
	defer db.Close()
	
	var acc Accounts
    err = db.QueryRow("SELECT * FROM tsao_accounts WHERE Username=? AND Password=? AND IsDeleted=0", loginAcc.Username, encyrptedPassword).Scan(
		&acc.ID, &acc.Name, 
		&acc.Username, &acc.Password, 
		&acc.Role, &acc.CreationDate, &acc.IsApproved, &acc.IsDeleted)
	fmt.Println(acc)
		
	if err == nil  {
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(acc)
	} else if err == sql.ErrNoRows{
		http.Error(w, "Account does not exist / Invalid credentials", http.StatusInternalServerError)
		return
	} 
}