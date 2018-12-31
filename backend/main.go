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

	r.HandleFunc("/user", u.Get).Methods("GET")
	r.HandleFunc("/user", u.Post).Methods("POST")
	r.HandleFunc("/user/login", u.LoginPost).Methods("POST")
	r.HandleFunc("/user/logout", u.LogoutPost).Methods("POST")
	r.HandleFunc("/user/invalidate", u.InvalidatePost).Methods("POST")
	r.HandleFunc("/user/delete", u.DeletePost).Methods("POST")
	r.HandleFunc("/user/authenticated", u.AuthenticatedGet).Methods("GET")

	log.Println("Listening on port 3000")
	http.ListenAndServe(":3000", r)
}
