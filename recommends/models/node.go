package models

import (
	"time"
)

const (
	UserNode = 0
	PostNode = 1
    UpvoteEdge = 0
    DownvoteEdge = 1
    CreationEdge = 2
)

type Edge struct {
	Destination *Node
    Type int // One of Upvote, Downvote, Creation
}

type Node struct {
	Id        string
	Type      int // One of UserNode, PostNode
	Edges     []Edge
	Timestamp time.Time
	Weight    float64
}
