package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/jbrunsting/transient/recommends/models"
)

type recommendsApi struct {
	graph map[string]*models.Node
}

func formatEdges(edges []models.Edge) string {
	edgeStrings := []string{}
	for _, edge := range edges {
		d := ""
		if edge.Destination.Type == models.UserNode {
			d = "u"
		} else if edge.Destination.Type == models.PostNode {
			d = "p"
		} else {
			d = "?"
		}

		t := ""
		if edge.Type == models.CreationEdge {
			t = "c"
		} else if edge.Type == models.UpvoteEdge {
			t = "+"
		} else if edge.Type == models.DownvoteEdge {
			t = "-"
		} else {
			t = "?"
		}

		edgeStrings = append(edgeStrings, fmt.Sprintf("%v -> %v-%v", t, d, edge.Destination.Id[0:5]))
	}

	return "[" + strings.Join(edgeStrings, ", ") + "]"
}

func (a *recommendsApi) NodePost(w http.ResponseWriter, r *http.Request) {
	var n models.NodeResource
	err := json.NewDecoder(r.Body).Decode(&n)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if n.Type != models.UserNode && n.Type != models.PostNode {
		http.Error(w,
			fmt.Sprintf("Invalid type, must be one of [%v, %v]", models.UserNode, models.PostNode),
			http.StatusBadRequest)
		return
	}

	var node models.Node
	node.Id = n.Id
	node.Type = n.Type
	node.Timestamp = n.Timestamp
	a.graph[node.Id] = &node

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (a *recommendsApi) EdgePost(w http.ResponseWriter, r *http.Request) {
	var e models.EdgeResource
	err := json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, ok := a.graph[e.DestinationId]; !ok {
		http.Error(w, "Invalid destination id", http.StatusBadRequest)
		return
	}

	if e.Type != models.UpvoteEdge && e.Type != models.DownvoteEdge && e.Type != models.CreationEdge {
		http.Error(w,
			fmt.Sprintf("Invalid type, must be one of [%v, %v, %v]", models.UpvoteEdge, models.DownvoteEdge, models.CreationEdge),
			http.StatusBadRequest)
		return
	}

	if _, ok := a.graph[e.SourceId]; ok {
		var edge models.Edge
		edge.Destination = a.graph[e.DestinationId]
		edge.Type = e.Type
		a.graph[e.SourceId].Edges = append(a.graph[e.SourceId].Edges, edge)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Invalid source id", http.StatusBadRequest)
	}
}

func (a *recommendsApi) RecommendsGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, ok := vars["id"]
	if !ok {
		http.Error(w, "Must provide an id", http.StatusBadRequest)
		return
	}

	log.Printf("ID is %s\n", id)

	log.Printf("Graph:")
	for k, _ := range a.graph {
		if a.graph[k].Type == models.UserNode && len(a.graph[k].Edges) > 0 {
			log.Printf("u-%v\n", k[0:5])
			for _, edge := range a.graph[k].Edges {
				t := ""
				p := ""
				if edge.Type == models.CreationEdge {
					t = "c"
				} else {
					p = fmt.Sprintf(" -> %v", formatEdges(edge.Destination.Edges))
					if edge.Type == models.UpvoteEdge {
						t = "+"
					} else if edge.Type == models.DownvoteEdge {
						t = "-"
					} else {
						t = "?"
					}
				}

				log.Printf("  %v -> %v%v\n", t, edge.Destination.Id[0:5], p)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
