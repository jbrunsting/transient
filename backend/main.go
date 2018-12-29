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
	usernameCookie  = "username"
	sessionIdCookie = "sessionId"
)

type Response struct {
	Message string `json:"message"`
}

var db *sql.DB

type Identification struct {
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}

type Session struct {
	SessionId string    `json:"sessionId"`
	Expiry    time.Time `json:"expiry"`
}

type User struct {
	Id string `json:"id"`
	Identification
	Email    string    `json:"email"`
	Sessions []Session `json:"sessions,omitempty"`
}

func generateSessionId() (string, error) {
	b := make([]byte, sessionIdLength)
	_, err := rand.Read(b)
	id := base64.URLEncoding.EncodeToString(b)[:sessionIdLength]
	return id, err
}

type HttpError struct {
	Msg  string
	Code int
}

func (e HttpError) Error() string {
	return e.Msg
}

func getUser(username string) (User, *HttpError) {
	var user User

	s := `
    SELECT Users.id, username, password, email, Sessions.sessionId, Sessions.expiry FROM Users
	LEFT JOIN Sessions ON Users.id = Sessions.id
	WHERE username=$1`
	rows, err := db.Query(s, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, &HttpError{Msg: "User ID not found", Code: http.StatusNotFound}
		}

		if sqlErr, ok := err.(*pq.Error); ok && sqlErr.Code.Class() == "08" {
			fmt.Printf("Error connecting to database: %v\n", err.Error())
			return user, &HttpError{Msg: "Error connecting to database", Code: http.StatusServiceUnavailable}
		}

		fmt.Printf("Error getting user from database: %v\n", err.Error())
		return user, &HttpError{Msg: "Error getting user from database", Code: http.StatusInternalServerError}
	}
	defer rows.Close()

	for rows.Next() {
		var session Session
		err = rows.Scan(&user.Id, &user.Username, &user.Password, &user.Email, &session.SessionId, &session.Expiry)
		if err != nil {
			fmt.Printf("Error parsing user from database: %v\n", err.Error())
			return user, &HttpError{Msg: "Error parsing user from database", Code: http.StatusInternalServerError}
		}
		user.Sessions = append(user.Sessions, session)
	}

	return user, nil
}

func authHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := ""
		sessionId := ""
		for _, cookie := range r.Cookies() {
			if cookie.Name == usernameCookie {
				username = cookie.Value
			} else if cookie.Name == sessionIdCookie {
				sessionId = cookie.Value
			}
		}

		if username == "" {
			http.Error(w, "No 'username' cookie set", http.StatusBadRequest)
			return
		} else if sessionId == "" {
			http.Error(w, "No 'sessionId' cookie set", http.StatusBadRequest)
			return
		}

		user, httpErr := getUser(username)
		if httpErr != nil {
			http.Error(w, httpErr.Error(), httpErr.Code)
			return
		}

		for _, session := range user.Sessions {
			if session.SessionId == sessionId {
				next(w, r)
				return
			}
		}

		http.Error(w, "Invalid session ID", http.StatusUnauthorized)
	}
}

func userPost(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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

	var session Session
	session.SessionId, err = generateSessionId()
	if err != nil {
		fmt.Printf("Error generating session ID: %v\n", err.Error())
		http.Error(w, "Error generating session ID", http.StatusInternalServerError)
		return
	}

	session.Expiry = time.Now()
	session.Expiry = session.Expiry.Add(expiryMinutes * time.Minute)

	tx, err := db.Begin()
	if err != nil {
		tx.Rollback()
		fmt.Printf("Error starting database transaction: %v\n", err.Error())
		http.Error(w, "Error starting database transaction", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(`
	INSERT INTO Users (id, username, password, email)
	VALUES ($1, $2, $3, $4)`, user.Id, user.Username, user.Password, user.Email)
	if err != nil {
		tx.Rollback()
		if sqlErr, ok := err.(*pq.Error); ok && (sqlErr.Code.Class() == "22" || sqlErr.Code.Class() == "23") {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			fmt.Printf("Error inserting user into database: %v\n", err.Error())
			http.Error(w, "Error inserting user into database", http.StatusInternalServerError)
		}
		return
	}

	_, err = tx.Exec(`
	INSERT INTO Sessions (id, sessionId, expiry)
	VALUES ($1, $2, $3)`, user.Id, session.SessionId, session.Expiry)
	if err != nil {
		tx.Rollback()
		fmt.Printf("Error inserting session into database: %v\n", err.Error())
		http.Error(w, "Error inserting session into database", http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		fmt.Printf("Error commiting database transaction: %v\n", err.Error())
		http.Error(w, "Error commiting database transaction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(session)
}

func userLoginPost(w http.ResponseWriter, r *http.Request) {
	var identification Identification
	err := json.NewDecoder(r.Body).Decode(&identification)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, httpErr := getUser(identification.Username)
	if err != nil {
		http.Error(w, httpErr.Error(), httpErr.Code)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(identification.Password))
	if err != nil {
		http.Error(w, "Password does not match", http.StatusUnauthorized)
		return
	}

	sessionId, err := generateSessionId()
	if err != nil {
		fmt.Printf("Error generating session ID: %v\n", err.Error())
		http.Error(w, "Error generating session ID", http.StatusInternalServerError)
		return
	}

	session := Session{SessionId: sessionId, Expiry: time.Now().Add(expiryMinutes * time.Minute)}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(session)
}

func handler(w http.ResponseWriter, r *http.Request) {
	response := Response{Message: "Hello world!"}
	json.NewEncoder(w).Encode(response)
}

func main() {
	var err error
	r := mux.NewRouter()

	r.HandleFunc("/", authHandler(handler))
	r.HandleFunc("/user", userPost).Methods("POST")
	r.HandleFunc("/user/login", userLoginPost).Methods("POST")

	db, err = sql.Open("postgres", "host=dev-db sslmode=disable user=transient password=password")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Println("Listening on port 3000")
	http.ListenAndServe(":3000", r)
}
