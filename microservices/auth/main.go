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
	"os"
	"path/filepath"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

var connectionString	string

func main() {
	r := mux.NewRouter()
    api := r.PathPrefix("/api/auth").Subrouter()
    api.HandleFunc("", Login).Methods("POST")
    http.Handle("/", r)

	fmt.Println("Auth microservice running on http://localhost:8000/api/auth")
	
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
	http.ListenAndServe(":8000", handler)
}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	log.Println("Entering endpoint to validate user credentials")
	
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
    err = db.QueryRow("SELECT * FROM tsao_accounts WHERE Username=? AND Password=?", loginAcc.Username, encyrptedPassword).Scan(
		&acc.ID, &acc.Name, 
		&acc.Username, &acc.Password, 
		&acc.Role, &acc.CreationDate, &acc.IsApproved, &acc.IsDeleted)
		
	if acc.IsDeleted {
		http.Error(w, "Account does not exist", http.StatusInternalServerError)
		return
	} else if err == sql.ErrNoRows{
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	} else if err == nil  {
		token, err := generateJWT(acc.Username)
		if err != nil {
			http.Error(w, "Error generating JWT token", http.StatusInternalServerError)
			return
		}
		
		http.SetCookie(w, &http.Cookie{
			Name:  "jwtToken",
			Value: token,
			HttpOnly: true,
			Secure: false,
			SameSite: http.SameSiteNoneMode,
			Domain: "127.0.0.1",
			Path: "/",
			MaxAge: 60*60*24*30,
		})
		
		resData := Response{
			ID: acc.ID,
			Name: acc.Name,
			Role: acc.Role,
			Token: token,
		}
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(resData)
	} else {
		http.Error(w, "Unexpected error occured", http.StatusNotFound)
	}
}

func generateJWT(username string) (string, error) {
	pwd, _ := os.Getwd()
	envDir := filepath.Dir(filepath.Dir(pwd))
	envPath := filepath.Join(envDir, ".env")
	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatal("Error loading .env file.")
	}
	
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Fatal("SECRET_KEY not found in .env file")
	}
	secretKeyBytes := []byte(secretKey)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, 
        jwt.MapClaims{ 
			"username": username, 
			"exp": time.Now().Add(time.Minute * 5).Unix(), 
        })

    tokenString, err := token.SignedString(secretKeyBytes)
	fmt.Println("Token string: ", tokenString)
    if err != nil {
    	return "", err
    }

	return tokenString, nil
}