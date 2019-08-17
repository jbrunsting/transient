package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/jbrunsting/transient/recommends/database"
	"github.com/jbrunsting/transient/recommends/models"
)

type recommendsApi struct {
	db database.DatabaseHandler
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

func (a *recommendsApi) RecommendsGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, ok := vars["id"]
	if !ok {
		http.Error(w, "Must provide an id", http.StatusBadRequest)
		return
	}

	log.Printf("ID is %s\n", id)

	graph, err := a.db.GenerateGraph()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("Graph:")
	for k, _ := range graph {
		if graph[k].Type == models.UserNode && len(graph[k].Edges) > 0 {
			log.Printf("u-%v\n", k[0:5])
			for _, edge := range graph[k].Edges {
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
	json.NewEncoder(w).Encode("test")
}
