package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"

	"github.com/jbrunsting/transient/backend/models"
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
			nodes[node.Id] = &node
		}
	}

	lookback := time.Now().AddDate(0, 0, -lookbackDays)
	s = `SELECT Posts.id, Posts.postId, Posts.time, Votes.id, Votes.vote FROM Posts
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
	for voteRows.Next() {
		if err = voteRows.Scan(&posterId, &postId, &postTime, &voterId, &vote); err != nil {
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
				}

			}

			if node, ok := nodes[posterId]; !ok {
				edge := models.Edge{
					Destination: nodes[postId],
					Type:        models.CreationEdge,
				}
				node.Edges = append(node.Edges, edge)
			}

			if voterId.Valid && vote.Valid {
				var edgeType int
				if vote.Int64 == 1 {
					edgeType = models.UpvoteEdge
				} else if vote.Int64 == -1 {
					edgeType = models.DownvoteEdge
				} else {
					return nodes, fmt.Errorf("Got unknown vote type %v", vote.Int64)
				}

				edge := models.Edge{
					Destination: nodes[postId],
					Type:        edgeType,
				}

				if node, ok := nodes[voterId.String]; ok {
					node.Edges = append(node.Edges, edge)
				} else {
					log.Printf("Unknown user id from vote, got id %v\n", voterId)
				}
			}
		}
	}

	return nodes, nil
}

func (h *recommendsHandler) GetUserFromId(id string) (models.User, error) {
	return h.getUser("Users.id = $1", id)
}

func (h *recommendsHandler) getUser(whereCondition string, whereArgs ...interface{}) (models.User, error) {
	var u models.User

	s := fmt.Sprintf(`
    SELECT Users.id, username, password, email, Sessions.sessionId, Sessions.expiry FROM Users
	LEFT JOIN Sessions ON Users.id = Sessions.id
	WHERE %v`, whereCondition)
	rows, err := h.db.Query(s, whereArgs...)
	if err != nil {
		return u, formatError(err, "user", "querying users")
	}
	defer rows.Close()

	if !rows.Next() {
		return u, &NotFoundError{"user"}
	}

	for {
		var sessionId sql.NullString
		var expiry pq.NullTime
		if err = rows.Scan(&u.Id, &u.Username, &u.Password, &u.Email, &sessionId, &expiry); err != nil {
			break
		}

		if sessionId.Valid && expiry.Valid {
			u.Sessions = append(u.Sessions, models.Session{
				SessionId: sessionId.String,
				Expiry:    expiry.Time,
			})
		}

		if !rows.Next() {
			break
		}
	}

	if rows.Err() != nil {
		err = rows.Err()
	}

	if err != nil {
		return u, &UnexpectedError{
			Action:        "parsing user",
			InternalError: err.Error(),
		}
	}

	return u, nil
}
