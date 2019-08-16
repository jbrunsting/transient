package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/jbrunsting/transient/backend/database"
)

type recommendsApi struct {
	db database.DatabaseHandler
}

func (a *recommendsApi) RecommendsGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, ok := vars["id"]
	if !ok {
		http.Error(w, "Must provide an id", http.StatusBadRequest)
		return
	}

    log.Printf("ID is %s\n", id)


	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("test")
}
