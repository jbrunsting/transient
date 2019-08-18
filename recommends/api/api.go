package api

import (
	"net/http"

	"github.com/jbrunsting/transient/recommends/models"
)

type Api interface {
	NodePost(w http.ResponseWriter, r *http.Request)
	EdgePost(w http.ResponseWriter, r *http.Request)
	RecommendsGet(w http.ResponseWriter, r *http.Request)
}

type api struct {
	recommendsApi
}

func NewApi(graph map[string]*models.Node) Api {
	return &api{recommendsApi: recommendsApi{graph: graph}}
}
