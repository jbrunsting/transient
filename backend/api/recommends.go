package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/jbrunsting/transient/backend/database"
)

type recommendsApi struct {
	db database.DatabaseHandler
}

const (
	userNode = 0
	postNode = 1

	upvoteEdge   = 0
	downvoteEdge = 1
	creationEdge = 2
	followEdge   = 3
)

type edgeResource struct {
	SourceId      string    `json:"sourceId"`
	DestinationId string    `json:"destinationId"`
	Type          int       `json:"type"`
	Timestamp     time.Time `json:"timestamp"`
}

type nodeResource struct {
	Id        string    `json:"id"`
	Type      int       `json:"type"`
	Timestamp time.Time `json:"timestamp"`
}

func addRecommendsEdge(e *edgeResource, bidirectional bool) {
	body, err := json.Marshal(e)
	if err != nil {
		log.Printf("Error marshalling object as json: %v\n", err)
	}

	_, err = http.Post("http://dev-recommends:4001/edge", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Error adding recommends edge: %v\n", err)
	}

	if bidirectional {
		reverse := *e
		reverse.SourceId, reverse.DestinationId = reverse.DestinationId, reverse.SourceId
		addRecommendsEdge(&reverse, false)
	}
}

func addRecommendsNode(n *nodeResource) {
	body, err := json.Marshal(n)
	if err != nil {
		log.Printf("Error marshalling object as json: %v\n", err)
	}

	_, err = http.Post("http://dev-recommends:4001/node", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Error adding recommends edge: %v\n", err)
	}
}

func (a *recommendsApi) RecommendsPostsGet(w http.ResponseWriter, r *http.Request) {
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

    resp, err := http.Get("http://dev-recommends:4001/posts/" + u.Id)
	if err != nil {
		log.Printf("Error getting recommended posts, %v\n", err)
		http.Error(w, "Could not generate post recommendations", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	postIds := []string{}
	err = json.NewDecoder(resp.Body).Decode(&postIds)
	if err != nil {
		log.Printf("Error decoding recommended posts, %v\n", err)
		http.Error(w, "Could not generate post recommendations", http.StatusServiceUnavailable)
		return
	}

	posts, err := a.db.GetPosts(postIds)
	if err != nil {
		handleDbErr(err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)
}
