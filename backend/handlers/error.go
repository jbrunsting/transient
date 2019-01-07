package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jbrunsting/transient/backend/database"
)

const (
	CONNECTION           = "connection"
	NOT_FOUND            = "not_found"
	DATA_VIOLATION       = "data_volation"
	UNIQUENESS_VIOLATION = "uniqueness_violation"
	UNEXPECTED           = "unexpected"
)

type httpError struct {
	Message string `json:"message"`
	Kind    string `json:"kind"`
}

func sendError(w http.ResponseWriter, code int, e httpError) {
	w.Header().Set("Content-Type", "application/json")
	log.Printf("Ssending error %v\n", e)
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(e)
}

func handleDbErr(err error, w http.ResponseWriter) {
	var code int
	var kind string
	switch e := err.(type) {
	case *database.ConnectionError:
		log.Printf("%v: %v", e.Error(), e.InternalError)
		code = http.StatusServiceUnavailable
		kind = CONNECTION
	case *database.NotFoundError:
		code = http.StatusNotFound
		kind = NOT_FOUND
	case *database.DataViolation:
		code = http.StatusBadRequest
		kind = DATA_VIOLATION
	case *database.UniquenessViolation:
		code = http.StatusBadRequest
		kind = UNIQUENESS_VIOLATION
	case *database.UnexpectedError:
		log.Printf("%v: %v", e.Error(), e.InternalError)
		code = http.StatusInternalServerError
		kind = UNEXPECTED
	}

	sendError(w, code, httpError{Message: err.Error(), Kind: kind})
}
