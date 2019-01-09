package database

import (
	"database/sql"

	"github.com/jbrunsting/transient/backend/models"
)

type DatabaseHandler interface {
	GetUserFromUsername(username string) (models.User, error)
	GetUserFromSession(sessionId string) (models.User, error)
	CreateUser(u models.User, s models.Session) error
	CreateSession(s models.Session) error
	DeleteOtherSessions(currentSessionId string) error
	DeleteSession(sessionId string) error
	DeleteUser(id string) error

	GetPosts(id string) ([]models.Post, error)
	CreatePost(p models.Post) error

	Close()
}

type databaseHandler struct {
	db *sql.DB
	userHandler
	postHandler
}

func NewDatabaseHandler() (DatabaseHandler, error) {
	db, err := sql.Open("postgres", "host=dev-db sslmode=disable user=transient password=password")
	if err != nil {
		return nil, err
	}
	return &databaseHandler{db: db, userHandler: userHandler{db}, postHandler: postHandler{db}}, nil
}

func (h *databaseHandler) Close() {
	h.db.Close()
}
