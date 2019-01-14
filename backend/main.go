package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/jbrunsting/transient/backend/api"
	"github.com/jbrunsting/transient/backend/database"
)

type response struct {
	Message string `json:"message"`
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

	a := api.NewApi(databaseHandler)

	r.HandleFunc("/user", a.SelfGet).Methods("GET")
	r.HandleFunc("/user/{username}", a.UserGet).Methods("GET")
	r.HandleFunc("/user", a.UserPost).Methods("POST")
	r.HandleFunc("/user/login", a.UserLoginPost).Methods("POST")
	r.HandleFunc("/user/logout", a.UserLogoutPost).Methods("POST")
	r.HandleFunc("/user/invalidate", a.UserInvalidatePost).Methods("POST")
	r.HandleFunc("/user/delete", a.UserDeletePost).Methods("POST")
	r.HandleFunc("/authenticated", a.UserAuthenticatedGet).Methods("GET")
	r.HandleFunc("/posts/{id}", a.PostsGet).Methods("GET")
	r.HandleFunc("/post", a.PostPost).Methods("POST")

	log.Println("Listening on port 3000")
	http.ListenAndServe(":3000", r)
}
