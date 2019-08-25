package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
    "time"

	"github.com/gorilla/mux"

	"github.com/jbrunsting/transient/recommends/models"
)

const (
    // Only use the first 5000 edges
	maxEdges = 5000
	startingWeight = 1000000000
	maxAgeHours = 720
    recommendsIterations = 10
)

var typeFractions = map[int]float64{
	models.CreationEdge: 0.1,
	models.UpvoteEdge:   0.004,
	models.DownvoteEdge: -0.02,
	models.FollowEdge:   0.2, // TODO: support
}

type recommendsApi struct {
	graph map[string]*models.Node
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

	if sourceNode, ok := a.graph[e.SourceId]; ok {
		if destinationNode, ok := a.graph[e.DestinationId]; ok {
			var edge models.Edge
			edge.Source = sourceNode
			edge.Destination = destinationNode
			edge.Type = e.Type
			edge.Timestamp = e.Timestamp
			models.AddEdge(edge)
            edge.Destination, edge.Source = edge.Source, edge.Destination
            models.AddEdge(edge)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
		} else {
			http.Error(w, "Invalid destination id", http.StatusBadRequest)
		}
	} else {
		http.Error(w, "Invalid source id", http.StatusBadRequest)
	}
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
		} else if edge.Type == models.FollowEdge {
			t = "f"
		} else {
			t = "?"
		}

		edgeStrings = append(edgeStrings, fmt.Sprintf("%v -> %v-%v", t, d, edge.Destination.Id[0:5]))
	}

	return "[" + strings.Join(edgeStrings, ", ") + "]"
}

func updateWeights(nodes []*models.Node, id string, now time.Time) []*models.Node {
	updatedWeights := map[*models.Node]float64{}
	for _, node := range nodes {
		for i := 0; i < len(node.Edges) && i < maxEdges; i++ {
			edge := node.Edges[i]

			if _, ok := updatedWeights[edge.Destination]; !ok {
				updatedWeights[edge.Destination] = 0
			}

			updatedWeights[edge.Destination] += node.Weights[id] * typeFractions[edge.Type]
		}
	}

    // Apply time penalty to posts
	for node, _ := range updatedWeights {
        if node.Type == models.PostNode {
            node.Weights[id] *= (maxAgeHours - now.Sub(node.Timestamp).Hours()) / maxAgeHours
            if node.Weights[id] < 0 {
                node.Weights[id] = 0
            }
        }
    }

	updated := []*models.Node{}
	for node, weight := range updatedWeights {
		updated = append(updated, node)
		if node.Weights[id] < weight {
			node.Weights[id] = weight
		}
	}

	return updated
}

func GenerateRecommends(start *models.Node) []*models.Node {
	id := start.Id

	nodes := []*models.Node{start}
	seen := map[*models.Node]bool{start: true}

	start.Weights[id] = startingWeight

    // TODO: Can use smaller number of iterations initially, and then in
    // background do more iterations to get more recommendations
	now := time.Now()
    for i := 0; i < recommendsIterations; i++ {
		nodes = updateWeights(nodes, id, now)
		for _, node := range nodes {
			seen[node] = true
		}
    }

	// Eliminate any nodes which have already been voted on
	for _, edge := range start.Edges {
		if edge.Destination.Type == models.PostNode {
			edge.Destination.Weights[id] = 0
		}
	}

	recommends := []*models.Node{}
	for node, _ := range seen {
        if node.Type == models.PostNode && node.Weights[id] > 0 {
            recommends = append(recommends, node)
        }
	}

	sort.Slice(recommends, func(i, j int) bool {
		return recommends[i].Weights[start.Id] < recommends[j].Weights[start.Id]
	})

	return recommends
}

func (a *recommendsApi) RecommendsGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, ok := vars["id"]
	if !ok {
		http.Error(w, "Must provide an id", http.StatusBadRequest)
		return
	}

	node, ok := a.graph[id]
	if !ok {
		http.Error(w, "Unknown ID", http.StatusBadRequest)
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
					} else if edge.Type == models.FollowEdge {
						t = "f"
					} else {
						t = "?"
					}
				}

				log.Printf("  %v -> %v%v\n", t, edge.Destination.Id[0:5], p)
			}
		}
	}

	recommends := GenerateRecommends(node)
	log.Printf("Recommends:")
	for _, node := range recommends {
		log.Printf(node.Id[0:5])
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
