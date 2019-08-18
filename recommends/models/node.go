package models

import (
	"time"
)

const (
	UserNode     = 0
	PostNode     = 1
	UpvoteEdge   = 0

	DownvoteEdge = 1
	CreationEdge = 2
)

type Edge struct {
	Destination *Node
	Type        int // One of Upvote, Downvote, Creation
}

type Node struct {
	Id        string
	Type      int // One of UserNode, PostNode
	Edges     map[string]Edge
	Timestamp time.Time
	Weights   map[string]float64
}

type EdgeResource struct {
	SourceId      string `json:"sourceId"`
	DestinationId string `json:"destinationId"`
	Type          int    `json:"type"`
}

type NodeResource struct {
	Id        string    `json:"id"`
	Type      int       `json:"type"`
	Timestamp time.Time `json:"timestamp"`
}

func (n *Node) AddEdge(e Edge) {
	n.Edges[e.Destination.Id] = e
}
