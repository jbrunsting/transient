package api

import (
	"net/http"

	"github.com/jbrunsting/transient/backend/database"
)

type Api interface {
	RecommendsGet(w http.ResponseWriter, r *http.Request)
}

type api struct {
	recommendsApi
}

func NewApi(db database.DatabaseHandler) Api {
    return &api{recommendsApi: recommendsApi{db: db}}
}
