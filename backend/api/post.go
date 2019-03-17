package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

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
		http.Error(w, "Not logged in", http.StatusUnauthorized)
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

	p.Time = time.Now()

	if err = a.db.CreatePost(p); err != nil {
		handleDbErr(err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (a *postApi) PostDelete(w http.ResponseWriter, r *http.Request) {
	sessionId, err := getSessionId(r)
	if err != nil {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}

	u, err := a.db.GetUserFromSession(sessionId)
	if err != nil {
		handleDbErr(err, w)
		return
	}

	vars := mux.Vars(r)

	postId, ok := vars["id"]
	if !ok {
		http.Error(w, "Must provide a post ID to delete", http.StatusBadRequest)
		return
	}

	post, err := a.db.GetPost(postId)
	if err != nil {
		handleDbErr(err, w)
		return
	}

	if post.Id != u.Id {
		http.Error(w, "Currently logged in user is not the owner of the post", http.StatusUnauthorized)
		return
	}

	err = a.db.DeletePost(postId)
	if err != nil {
		handleDbErr(err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (a *postApi) PostVotePost(w http.ResponseWriter, r *http.Request) {
	sessionId, err := getSessionId(r)
	if err != nil {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}

	u, err := a.db.GetUserFromSession(sessionId)
	if err != nil {
		handleDbErr(err, w)
		return
	}

	vars := mux.Vars(r)

	postId, ok := vars["id"]
	if !ok {
		http.Error(w, "Must provide a post ID to vote on", http.StatusBadRequest)
		return
	}

	var v models.Vote
	err = json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if v.Vote != models.UPVOTE && v.Vote != models.DOWNVOTE {
		http.Error(w, fmt.Sprintf("Vote must be %v or %v", models.UPVOTE, models.DOWNVOTE), http.StatusBadRequest)
		return
	}

	err = a.db.CreateVote(u.Id, postId, v.Vote)
	if err != nil {
		handleDbErr(err, w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (a *postApi) PostCommentPost(w http.ResponseWriter, r *http.Request) {
	sessionId, err := getSessionId(r)
	if err != nil {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}

	u, err := a.db.GetUserFromSession(sessionId)
	if err != nil {
		handleDbErr(err, w)
		return
	}

	vars := mux.Vars(r)

	postId, ok := vars["id"]
	if !ok {
		http.Error(w, "Must provide a post ID to comment on", http.StatusBadRequest)
		return
	}

	var c models.Comment
	err = json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

    c.Id = u.Id

	id, err := uuid.NewV4()
	if err != nil {
		log.Printf("Error generating UUID: %v\n", err.Error())
		http.Error(w, "Error generating UUID", http.StatusInternalServerError)
		return
	}
	c.CommentId = id.String()

	c.Time = time.Now()

	err = a.db.CreateComment(postId, c)
	if err != nil {
		handleDbErr(err, w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (a *postApi) PostCommentsGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, ok := vars["id"]
	if !ok {
		http.Error(w, "Must provide a post ID to get the posts for", http.StatusBadRequest)
		return
	}

	comments, err := a.db.GetComments(id)
	if err != nil {
		handleDbErr(err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(comments)
}
