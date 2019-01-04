package main

import (
	"log"
	"net/http"
    "encoding/json"

	"github.com/gorilla/mux"

	"github.com/jbrunsting/transient/backend/handlers"
	"github.com/jbrunsting/transient/backend/database"
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

	databaseHandler, err := database.NewDatabaseHandler()
	if err != nil {
		panic(err)
	}
	defer databaseHandler.Close()

	u := handlers.NewUserHandler(databaseHandler)

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
