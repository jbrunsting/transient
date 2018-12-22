package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

type Response struct {
	Message string `json:"message"`
}

var db *sql.DB

func handler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
        fmt.Println("Got error %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return;
	}

	fmt.Println("Iterating over rows")
	for rows.Next() {
		var id string
		var username string
		var password string
		var session string
		var expiry string

		err = rows.Scan(&id, &username, &password, &session, &expiry)
		fmt.Printf("At row id: %v, username: %v, password: %v, session: %v, expiry: %v\n", id, username, password, session, expiry)
		if err != nil {
            fmt.Printf("Got error %v\n", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return;
		}
	}

	response := Response{Message: "Hello world!"}
	json.NewEncoder(w).Encode(response)
}

func main() {
	var err error
	r := mux.NewRouter()

	r.HandleFunc("/", handler)

	db, err = sql.Open("postgres", "host=dev-db sslmode=disable user=transient password=password")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.QueryRow("INSERT INTO users(id, username, password, session, expiry) VALUES('test_id', 'test_username', 'test_password', 'test_session', 'test_expiry')").Scan()
	if err != nil {
		fmt.Printf("Ignoring insertion error %v\n", err)
	}

	fmt.Println("Listening on port 3000")
	http.ListenAndServe(":3000", r)
}
