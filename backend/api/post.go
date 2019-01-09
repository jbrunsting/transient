package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"

	"github.com/jbrunsting/transient/backend/database"
	"github.com/jbrunsting/transient/backend/models"
)

type postApi struct {
	db database.DatabaseHandler
}

func (a *postApi) PostsGet(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)

    id, ok := vars["id"]
    if !ok {
        http.Error(w, "Must provide a user ID to get the posts for", http.StatusBadRequest)
        return
    }

	posts, err := a.db.GetPosts(id)
	if err != nil {
		handleDbErr(err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)
}

func (a *postApi) PostPost(w http.ResponseWriter, r *http.Request) {
	sessionId, err := getSessionId(r)
	if err != nil {
		http.Error(w, "Not logged in", http.StatusForbidden)
		return
	}

	u, err := a.db.GetUserFromSession(sessionId)
	if err != nil {
		handleDbErr(err, w)
		return
	}

	var p models.Post
	err = json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p.Id = u.Id

	id, err := uuid.NewV4()
	if err != nil {
		log.Printf("Error generating UUID: %v\n", err.Error())
		http.Error(w, "Error generating UUID", http.StatusInternalServerError)
		return
	}
	p.PostId = id.String()

	if err = a.db.CreatePost(p); err != nil {
		handleDbErr(err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
