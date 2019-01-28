package database

import (
	"database/sql"

	"github.com/jbrunsting/transient/backend/models"
)

type DatabaseHandler interface {
	GetUserFromUsername(username string) (models.User, error)
	GetUserFromSession(sessionId string) (models.User, error)
	GetUserFromId(sessionId string) (models.User, error)
	CreateUser(u models.User, s models.Session) error
	CreateSession(s models.Session) error
	DeleteOtherSessions(currentSessionId string) error
	DeleteSession(sessionId string) error
	DeleteUser(id string) error
	SearchUsers(search string, limit int) ([]models.User, error)

	GetPosts(id string) ([]models.Post, error)
	GetPost(postId string) (models.Post, error)
	CreatePost(p models.Post) error
	DeletePost(postId string) error

	CreateFollowing(id, followingId string) error
	GetFollowings(id string) ([]string, error)
	DeleteFollowing(id, followingId string) error

    Close()
}

type databaseHandler struct {
	db *sql.DB
	userHandler
	postHandler
	followingHandler	
}

func NewDatabaseHandler() (DatabaseHandler, error) {
	db, err := sql.Open("postgres", "host=dev-db sslmode=disable user=transient password=password")
	if err != nil {
		return nil, err
	}
	return &databaseHandler{db: db, userHandler: userHandler{db}, postHandler: postHandler{db}, followingHandler: followingHandler{db}}, nil
}

func (h *databaseHandler) Close() {
	h.db.Close()
}
