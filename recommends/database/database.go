package database

import (
	"database/sql"

	"github.com/jbrunsting/transient/recommends/models"
)

type DatabaseHandler interface {
	GenerateGraph() (map[string]*models.Node, error)

	Close()
}

type databaseHandler struct {
	db *sql.DB
	recommendsHandler
}

func NewDatabaseHandler() (DatabaseHandler, error) {
	db, err := sql.Open("postgres", "host=dev-db sslmode=disable user=transient password=password")
	if err != nil {
		return nil, err
	}
	return &databaseHandler{db: db, recommendsHandler: recommendsHandler{db}}, nil
}

func (h *databaseHandler) Close() {
	h.db.Close()
}
