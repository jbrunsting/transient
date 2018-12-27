package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

type Response struct {
	Message string `json:"message"`
}

var db *sql.DB

type User struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Session  string `json:"session"`
	Expiry   string `json:"expiry"`
}

func userPost(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s := `
    INSERT INTO users (id, email, username, password, session, expiry)
    VALUES ($1, $2, $3, $4, $5, $6)`

	id, err := uuid.NewV4()
	if err != nil {
		fmt.Printf("Error creating UUID: %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = db.Exec(s, id.String(), user.Email, user.Username, user.Password, "", "")
	if err != nil {
		if sqlErr, ok := err.(*pq.Error); ok {
			switch sqlErr.Code.Class() {
			case "08":
				fmt.Printf("Error inserting user: %v", err.Error())
				http.Error(w, err.Error(), http.StatusServiceUnavailable)
			case "22", "23", "42":
				http.Error(w, err.Error(), http.StatusBadRequest)
			default:
				fmt.Printf("Error inserting user: %v", err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		} else {
			fmt.Printf("Error inserting user: %v", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
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
		var email string
		var username string
		var password string
		var session string
		var expiry string

		err = rows.Scan(&id, &email, &username, &password, &session, &expiry)
		fmt.Printf("At row id: %v, email: %v, username: %v, password: %v, session: %v, expiry: %v\n", id, email, username, password, session, expiry)
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

	db, err = sql.Open("postgres", "host=dev-db sslmode=disable user=transient password=password")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.QueryRow("INSERT INTO users(id, email, username, password, session, expiry) VALUES('test_id', 'test_username', 'test_email', 'test_password', 'test_session', 'test_expiry')").Scan()
	if err != nil {
		fmt.Printf("Ignoring insertion error %v\n", err)
	}

	fmt.Println("Listening on port 3000")
	http.ListenAndServe(":3000", r)
}
