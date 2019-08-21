package models

import (
	"sort"
	"time"
)

const (
	UserNode = 0
	PostNode = 1

	UpvoteEdge   = 0
	DownvoteEdge = 1
	CreationEdge = 2
	FollowEdge   = 3 // TODO: Still have to add followings to graph

	// Edges will be given priority over edges which are hourDiffForPriority
	// hours older, regardless of type
	hourDiffForPriority = 100
)

var edgeRankings = []int{FollowEdge, CreationEdge, DownvoteEdge, UpvoteEdge}

type Edge struct {
	Destination *Node
	Type        int // One of Upvote, Downvote, Creation
	Timestamp   time.Time
}

type Node struct {
	Id           string
	Type         int // One of UserNode, PostNode
	Edges        []Edge
	Destinations map[string]bool
	Timestamp    time.Time
	Weights      map[string]float64 // TODO: Use map that allows concurrent access
}

type EdgeResource struct {
	SourceId      string    `json:"sourceId"`
	DestinationId string    `json:"destinationId"`
	Type          int       `json:"type"`
	Timestamp     time.Time `json:"timestamp"`
}

type NodeResource struct {
	Id        string    `json:"id"`
	Type      int       `json:"type"`
	Timestamp time.Time `json:"timestamp"`
}

func (n *Node) AddEdge(e Edge) {
	if _, ok := n.Destinations[e.Destination.Id]; !ok {
		if n.Edges == nil {
			n.Edges = []Edge{}
		}
		n.Edges = append(n.Edges, e)
	}

	// TODO: This is probably slow, since we are sorting every time we add an
	// edge, instead of just inserting in sorted order
	sort.Slice(n.Edges, func(i, j int) bool {
		ei := n.Edges[i]
		ej := n.Edges[j]

		timeDiff := ei.Timestamp.Sub(ej.Timestamp).Hours()
		if timeDiff > hourDiffForPriority {
			return true
		} else if timeDiff < -hourDiffForPriority {
			return false
		}

		for edgeType := range edgeRankings {
			if ei.Type == edgeType && ej.Type != edgeType {
				return true
			}

			if ej.Type == edgeType && ei.Type != edgeType {
				return false
			}
		}

		return false
	})
}
