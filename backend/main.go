package main

import (
	"log"
	"database/sql"
	"net/http"
    "encoding/json"

	"github.com/gorilla/mux"

	"github.com/jbrunsting/transient/backend/user"
)

type response struct {
	Message string `json:"message"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(response{Message: "Hello world!"})
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var err error
	r := mux.NewRouter()

	db, err := sql.Open("postgres", "host=dev-db sslmode=disable user=transient password=password")
	if err != nil {
		panic(err)
	}
	defer db.Close()

    u := user.UserHandler{DB: db}

	r.HandleFunc("/", u.AuthHandler(handler))
	r.HandleFunc("/user", u.Post).Methods("POST")
	r.HandleFunc("/user/login", u.LoginPost).Methods("POST")

	log.Println("Listening on port 3000")
	http.ListenAndServe(":3000", r)
}
