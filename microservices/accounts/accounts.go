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
	"strconv"
	"time"

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

var connectionString string

func main() {
	r := mux.NewRouter()
    api := r.PathPrefix("/api/accounts").Subrouter()
    api.HandleFunc("", CreateAccount).Methods("POST")
    api.HandleFunc("", ReadAccounts).Methods("GET")
    api.HandleFunc("/{id}", ReadAccount).Methods("GET")
    api.HandleFunc("/{id}", UpdateAccount).Methods("PATCH")
    api.HandleFunc("/modify-password/{id}", ModifyPassword).Methods("PATCH")
    api.HandleFunc("/approve/{id}", ApproveAccount).Methods("PATCH")
    api.HandleFunc("/{id}", DeleteAccount).Methods("DELETE")
    http.Handle("/", r)

	fmt.Println("Accounts microservice running on http://localhost:8002/api/accounts")
	
	// CORS configuration
    corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://127.0.0.1:5502"}, // Your frontend origin
        AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowedHeaders: []string{"Content-Type"},
    })
	
	cmd := flag.String("sql", "", "")
	flag.Parse()
	connectionString = string(*cmd)

    handler := corsHandler.Handler(r)
	http.ListenAndServe(":8002", handler)
}

func CreateAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	log.Println("Entering endpoint to add new account")

	var newAccount Accounts
	err := json.NewDecoder(r.Body).Decode(&newAccount)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	sha := sha256.New()
	sha.Write([]byte(newAccount.Password))
	newAccount.Password = hex.EncodeToString(sha.Sum(nil))

	location, _ := time.LoadLocation("Asia/Singapore")
	newAccount.CreationDate = time.Now().In(location).Format("2006-01-02 15:04:05")

	newAccount.IsApproved = false
	newAccount.IsDeleted = false

	db, _ := sql.Open("mysql", connectionString)
	defer db.Close()

	_, existingUsername, err := checkInfo(db, newAccount.Username)
    if err != nil {
        log.Println(err)
        http.Error(w, "Error checking existing user", http.StatusInternalServerError)
        return
    }

    if existingUsername != "" {
        http.Error(w, "Username already exists", http.StatusConflict)
        return
    }

	result, err := db.Exec(
		`INSERT INTO tsao_accounts (Name, Username, Password, Role, CreationDate, IsApproved, IsDeleted)
		 VALUES (?, ?, ?, ?, ?, ?, ?);`,
		newAccount.Name, newAccount.Username, newAccount.Password, newAccount.Role,
		newAccount.CreationDate, newAccount.IsApproved, newAccount.IsDeleted)
	if err == nil  {
		accountID, _ := result.LastInsertId()
		newAccount.ID = int(accountID)
		
		newAccountJson, _ := json.Marshal(newAccount)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		w.Write(newAccountJson)
	} else {
		log.Println(err)
		http.Error(w, "Error creating new account", http.StatusInternalServerError)
		return
	}
}

func ReadAccounts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	log.Println("Entering endpoint to query all accounts")

	db, _ := sql.Open("mysql", connectionString)
	defer db.Close()
	
    result, err := db.Query("SELECT * FROM tsao_accounts WHERE IsDeleted=false")
	
	var accounts []Accounts
	for result.Next() {
		var account Accounts
		_ = result.Scan(
			&account.ID, &account.Name,
			&account.Username, &account.Password, &account.Role,
			&account.CreationDate, &account.IsApproved, &account.IsDeleted)
		accounts = append(accounts, account)
	}

	if err == nil  {
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(accounts)
	} else if err := result.Err(); err != sql.ErrNoRows {
		log.Println(err)
		http.Error(w, "Error iterating over rows", http.StatusInternalServerError)
		return
	}
}

func ReadAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	log.Println("Entering endpoint to read an account's information")

	accountID := mux.Vars(r)["id"]	

	db, _ := sql.Open("mysql", connectionString)
	defer db.Close()
	
	var account Accounts
    err := db.QueryRow("SELECT * FROM tsao_accounts WHERE ID=?", accountID).Scan(
		&account.ID, &account.Name,
		&account.Username, &account.Password, &account.Role,
		&account.CreationDate, &account.IsApproved, &account.IsDeleted)

	if err == nil  {
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(account)
	} else if err == sql.ErrNoRows{
		http.Error(w, "Account does not exist", http.StatusNotFound)
		return
	} else {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func UpdateAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	log.Println("Entering endpoint to update an account's information")

	accountID, _ := strconv.Atoi(mux.Vars(r)["id"])
	
	var modifiedAccount Accounts
	err := json.NewDecoder(r.Body).Decode(&modifiedAccount)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	
	db, _ := sql.Open("mysql", connectionString)
	defer db.Close()

	existingID, _, err := checkInfo(db, modifiedAccount.Username)
    if err != nil {
        log.Println(err)
        http.Error(w, "Error checking existing user", http.StatusInternalServerError)
        return
    }

    if existingID != accountID && existingID != 0 {
        http.Error(w, "Username already exists", http.StatusConflict)
        return
    }

	result, err := db.Exec(
		`UPDATE tsao_accounts
		SET Name=?, Username=?, Role=?
		WHERE ID=? AND IsDeleted=false;`,
		modifiedAccount.Name, modifiedAccount.Username, modifiedAccount.Role, accountID)
	rowsAffected, _ := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		log.Println(err)
		http.Error(w, "Account ID does not exist", http.StatusNotFound)
		return
	}  else {
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintf(w, "Account information modified for ID: %d\n", accountID)
	}
}

func ModifyPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	log.Println("Entering endpoint to modify an account's password")

	accountID := mux.Vars(r)["id"]	

	db, _ := sql.Open("mysql", connectionString)
	defer db.Close()

	var newPassword string
	err := json.NewDecoder(r.Body).Decode(&newPassword)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	sha := sha256.New()
	sha.Write([]byte(newPassword))
	newPassword = hex.EncodeToString(sha.Sum(nil))

	result, err := db.Exec(
		`UPDATE tsao_accounts SET Password=? WHERE ID=? AND IsDeleted=false;`,
		newPassword, accountID)
	rowsAffected, _ := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		log.Println(err)
		http.Error(w, "Account ID does not exist", http.StatusNotFound)
		return
	}  else {
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintf(w, "Password has been modified for ID: %s\n", accountID)
	}
}

func ApproveAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	log.Println("Entering endpoint to approve an account")

	accountID := mux.Vars(r)["id"]	

	db, _ := sql.Open("mysql", connectionString)
	defer db.Close()

	result, err := db.Exec(
		`UPDATE tsao_accounts SET IsApproved=true WHERE ID=? AND IsDeleted=false;`, accountID)
	rowsAffected, _ := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		log.Println(err)
		http.Error(w, "Account already approved", http.StatusNotFound)
		return
	}  else {
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintf(w, "Account approved for ID: %s\n", accountID)
	}
}

func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	log.Println("Entering endpoint to (soft) delete an account")
}

func checkInfo(db *sql.DB, username string) (int, string, error) {
    var retrievedID int
	var retrievedUsername string
    err := db.QueryRow(`
        SELECT ID, Username FROM tsao_accounts
        WHERE Username=?;`, username).Scan(
        &retrievedID, &retrievedUsername)

    if err != nil {
        if err == sql.ErrNoRows {
            return 0, "", nil // No existing user found
        }
        return 0, "", err
    }
    return retrievedID, retrievedUsername, nil
}