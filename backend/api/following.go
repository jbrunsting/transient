package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/jbrunsting/transient/backend/database"
)

type followingApi struct {
	db database.DatabaseHandler
}

func (a *followingApi) FollowingsGet(w http.ResponseWriter, r *http.Request) {
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

	followings, err := a.db.GetFollowings(u.Id)
	if err != nil {
		handleDbErr(err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(followings)
}

func (a *followingApi) FollowingPost(w http.ResponseWriter, r *http.Request) {
	sessionId, err := getSessionId(r)
	if err != nil {
		http.Error(w, "Not logged in", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)

	id, ok := vars["id"]
	if !ok {
		http.Error(w, "Must provide a user ID to follow", http.StatusBadRequest)
		return
	}

	u, err := a.db.GetUserFromSession(sessionId)
	if err != nil {
		handleDbErr(err, w)
		return
	}

	err = a.db.CreateFollowing(u.Id, id)
	if err != nil {
		handleDbErr(err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (a *followingApi) FollowingDelete(w http.ResponseWriter, r *http.Request) {
	sessionId, err := getSessionId(r)
	if err != nil {
		http.Error(w, "Not logged in", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)

	id, ok := vars["id"]
	if !ok {
		http.Error(w, "Must provide a following to delete", http.StatusBadRequest)
		return
	}

	u, err := a.db.GetUserFromSession(sessionId)
	if err != nil {
		handleDbErr(err, w)
		return
	}

	err = a.db.DeleteFollowing(u.Id, id)
	if err != nil {
		handleDbErr(err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
