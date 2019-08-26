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

    graph, err := databaseHandler.GenerateGraph()
    if err != nil {
        panic(err)
    }

	a := api.NewApi(graph)

	r.HandleFunc("/posts/{id}", a.PostsGet).Methods("GET")
	r.HandleFunc("/followings/{id}", a.FollowingsGet).Methods("GET")
	r.HandleFunc("/edge", a.EdgePost).Methods("POST")
	r.HandleFunc("/node", a.NodePost).Methods("POST")

	log.Println("Listening on port 4000")
	http.ListenAndServe(":4000", r)
}
