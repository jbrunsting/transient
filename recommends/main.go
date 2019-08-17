package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/jbrunsting/transient/recommends/api"
	"github.com/jbrunsting/transient/recommends/database"
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

	r.HandleFunc("/recommends/{id}", a.RecommendsGet).Methods("GET")

	log.Println("Listening on port 4000")
	http.ListenAndServe(":4000", r)
}
