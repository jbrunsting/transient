package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

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

	_, err = http.Post("http://dev-recommends:4001/edge", "image/jpeg", bytes.NewBuffer(body))
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

	_, err = http.Post("http://localhost:4001/node", "image/jpeg", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Error adding recommends edge: %v\n", err)
	}
}
