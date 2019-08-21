package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/jbrunsting/transient/recommends/models"
)

const (
	lookbackDays = 20
)

type recommendsHandler struct {
	db *sql.DB
}

func (h *recommendsHandler) GenerateGraph() (map[string]*models.Node, error) {
	nodes := make(map[string]*models.Node, 100000)

	s := `SELECT id FROM Users`
	userRows, err := h.db.Query(s)
	if err != nil {
		return nodes, formatError(err, "user", "querying user ID's")
	}
	defer userRows.Close()

	for userRows.Next() {
		var node models.Node
		if err = userRows.Scan(&node.Id); err != nil {
			// Don't return error because we don't want a single error to abort
			// the whole operation. In the future, we may want to handle this
			// better
			log.Printf("Error reading user row: %s\n", err)
		} else {
            node.Weights = map[string]float64{}
			nodes[node.Id] = &node
		}
	}

	lookback := time.Now().AddDate(0, 0, -lookbackDays)
	s = `SELECT Posts.id, Posts.postId, Posts.time, Votes.id, Votes.vote, Votes.time FROM Posts
    LEFT JOIN Votes ON Votes.postId = Posts.postId WHERE Posts.time > $1`
	voteRows, err := h.db.Query(s, lookback)
	if err != nil {
		return nodes, formatError(err, "vote", "querying votes")
	}
	defer voteRows.Close()

	var posterId string
	var postId string
	var postTime time.Time
	var voterId sql.NullString
	var vote sql.NullInt64
    var voteTime *time.Time
	for voteRows.Next() {
		if err = voteRows.Scan(&posterId, &postId, &postTime, &voterId, &vote, &voteTime); err != nil {
			// Don't return error because we don't want a single error to abort
			// the whole operation. In the future, we may want to handle this
			// better
			log.Printf("Error reading post row: %s\n", err)
		} else {
			if _, ok := nodes[postId]; !ok {
				nodes[postId] = &models.Node{
					Id:        postId,
					Type:      models.PostNode,
					Timestamp: postTime,
					Weights: map[string]float64{},
				}
			}

			if posterNode, ok := nodes[posterId]; ok {
				forwardCreationEdge := models.Edge{
					Destination: nodes[postId],
					Type:        models.CreationEdge,
                    Timestamp:   postTime,
				}
				posterNode.AddEdge(forwardCreationEdge)

				reverseCreationEdge := models.Edge{
					Destination: posterNode,
					Type:        models.CreationEdge,
                    Timestamp:   postTime,
				}
				nodes[postId].AddEdge(reverseCreationEdge)

				if voterId.Valid && vote.Valid {
					var edgeType int
					if vote.Int64 == 1 {
						edgeType = models.UpvoteEdge
					} else if vote.Int64 == -1 {
						edgeType = models.DownvoteEdge
					} else {
						return nodes, fmt.Errorf("Got unknown vote type %v", vote.Int64)
					}

					forwardVoteEdge := models.Edge{
						Destination: nodes[postId],
						Type:        edgeType,
                        Timestamp:   *voteTime,
					}

					if voterNode, ok := nodes[voterId.String]; ok {
						voterNode.AddEdge(forwardVoteEdge)

						reverseVoteEdge := models.Edge{
							Destination: voterNode,
							Type:        edgeType,
                            Timestamp:   *voteTime,
						}

						nodes[postId].AddEdge(reverseVoteEdge)
					} else {
						log.Printf("Unknown user id from vote, got id %v\n", voterId)
					}
				}
			} else {
				log.Printf("Unknown poster id, got id %v\n", posterId)
			}
		}
	}

	return nodes, nil
}
