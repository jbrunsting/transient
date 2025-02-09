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
	FollowEdge   = 3

	// Edges will be given priority over edges which are hourDiffForPriority
	// hours older, regardless of type
	hourDiffForPriority = 100
)

var edgeRankings = []int{FollowEdge, CreationEdge, DownvoteEdge, UpvoteEdge}

type Edge struct {
	Source      *Node
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

func (n *Node) SortEdges() {
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

func AddEdge(e Edge) {
	if _, ok := e.Source.Destinations[e.Destination.Id]; !ok {
		if e.Source.Edges == nil {
			e.Source.Edges = []Edge{}
			e.Source.Destinations = map[string]bool{}
		}
		e.Source.Edges = append(e.Source.Edges, e)
		e.Source.Destinations[e.Destination.Id] = true
	}

	e.Source.SortEdges()
}
