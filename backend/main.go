package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/lib/pq"

	"golang.org/x/crypto/bcrypt"
)

const (
	sessionIdLength = 32
	expiryMinutes   = 86400
	timeFormat      = time.RFC3339
	bcryptCost      = 10
)

type Response struct {
	Message string `json:"message"`
}

var db *sql.DB

type Identification struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	Id string `json:"id"`
	Identification
	Email   string    `json:"email"`
	Session string    `json:"session"`
	Expiry  time.Time `json:"expiry"`
}

func getSessionId() (string, error) {
	b := make([]byte, sessionIdLength)
	_, err := rand.Read(b)
	id := base64.URLEncoding.EncodeToString(b)[:sessionIdLength]
	fmt.Println("Getting sesion id")
	fmt.Printf("session ID is %v\n", id)
	fmt.Println("Got sesion id")
	return id, err
}

func getUser(username string) (User, error) {
	s := `
    SELECT id, username, password, email, session, expiry FROM users WHERE username=$1`
	row := db.QueryRow(s, username)

	var user User
	err := row.Scan(&user.Id, &user.Username, &user.Password, &user.Email, &user.Session, &user.Expiry)

	return user, err
}

func userPost(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s := `
    INSERT INTO users (id, username, password, email, session, expiry)
    VALUES ($1, $2, $3, $4, $5, $6)`

	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcryptCost)
	if err != nil {
		fmt.Printf("Error hashing password: %v\n", err.Error())
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	user.Password = string(password)

	id, err := uuid.NewV4()
	if err != nil {
		fmt.Printf("Error generating UUID: %v\n", err.Error())
		http.Error(w, "Error generating UUID", http.StatusInternalServerError)
		return
	}
	user.Id = id.String()

	session, err := getSessionId()
	if err != nil {
		fmt.Printf("Error generating session ID: %v\n", err.Error())
		http.Error(w, "Error generating session ID", http.StatusInternalServerError)
		return
	}
	user.Session = session

	user.Expiry = time.Now()
	user.Expiry = user.Expiry.Add(expiryMinutes * time.Minute)

	fmt.Printf("Inserting user %v\n", user)

	_, err = db.Exec(s, user.Id, user.Username, user.Password, user.Email, user.Session, user.Expiry)
	if err != nil {
		if sqlErr, ok := err.(*pq.Error); ok {
			switch sqlErr.Code.Class() {
			case "08":
				fmt.Printf("Error connecting to database: %v\n", err.Error())
				http.Error(w, "Error connecting to database", http.StatusServiceUnavailable)
				return
			case "22", "23", "42":
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		fmt.Printf("Error inserting user into database: %v\n", err.Error())
		http.Error(w, "Error inserting user into database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func userLoginPost(w http.ResponseWriter, r *http.Request) {
	var identification Identification
	err := json.NewDecoder(r.Body).Decode(&identification)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := getUser(identification.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User ID not found", http.StatusNotFound)
			return
		}

		if sqlErr, ok := err.(*pq.Error); ok && sqlErr.Code.Class() == "08" {
			fmt.Printf("Error connecting to database: %v\n", err.Error())
			http.Error(w, "Error connecting to database", http.StatusServiceUnavailable)
			return
		}

		fmt.Printf("Error getting user from database: %v\n", err.Error())
		http.Error(w, "Error getting user from database", http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(identification.Password))
	if err != nil {
		http.Error(w, "Password does not match", http.StatusUnauthorized)
		return
	}

	session, err := getSessionId()
	if err != nil {
		fmt.Printf("Error generating session ID: %v\n", err.Error())
		http.Error(w, "Error generating session ID", http.StatusInternalServerError)
		return
	}
	user.Session = session

	user.Expiry = time.Now()
	user.Expiry = user.Expiry.Add(expiryMinutes * time.Minute)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func handler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		fmt.Println("Got error %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	fmt.Println("Iterating over rows")
	for rows.Next() {
		var id string
		var username string
		var password string
		var email string
		var session string
		var expiry string

		err = rows.Scan(&id, &username, &password, &email, &session, &expiry)
		fmt.Printf("At row id: %v, username: %v, password: %v, email: %v, session: %v, expiry: %v\n", id, password, email, username, session, expiry)
		if err != nil {
			fmt.Printf("Got error %v\n", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
	}

	response := Response{Message: "Hello world!"}
	json.NewEncoder(w).Encode(response)
}

func main() {
	var err error
	r := mux.NewRouter()

	r.HandleFunc("/", handler)
	r.HandleFunc("/user", userPost).Methods("POST")
	r.HandleFunc("/user/login", userLoginPost).Methods("POST")

	db, err = sql.Open("postgres", "host=dev-db sslmode=disable user=transient password=password")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.QueryRow("INSERT INTO users(id, username, password, email, session, expiry) VALUES('test_id', 'test_username', 'test_password', 'test_email', 'test_session', 'test_expiry')").Scan()
	if err != nil {
		fmt.Printf("Ignoring insertion error %v\n", err)
	}

	fmt.Println("Listening on port 3000")
	http.ListenAndServe(":3000", r)
}
