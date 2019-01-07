package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gofrs/uuid"

	"github.com/jbrunsting/transient/backend/database"
	"github.com/jbrunsting/transient/backend/models"
)

type UserHandler interface {
	Get(w http.ResponseWriter, r *http.Request)
	AuthenticatedGet(w http.ResponseWriter, r *http.Request)

	Post(w http.ResponseWriter, r *http.Request)
	LoginPost(w http.ResponseWriter, r *http.Request)
	LogoutPost(w http.ResponseWriter, r *http.Request)
	InvalidatePost(w http.ResponseWriter, r *http.Request)
	DeletePost(w http.ResponseWriter, r *http.Request)
}

type userHandler struct {
	db database.DatabaseHandler
}

func NewUserHandler(db database.DatabaseHandler) UserHandler {
	return &userHandler{db: db}
}

func (h *userHandler) Get(w http.ResponseWriter, r *http.Request) {
	sessionId, err := getSessionId(r)
	if err != nil {
		http.Error(w, "Not logged in", http.StatusForbidden)
		return
	}

	u, err := h.db.GetUserFromSession(sessionId)
	if err != nil {
		handleDbErr(err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.User{
		Id:             u.Id,
		Identification: models.Identification{Username: u.Username},
		Email:          u.Email,
	})
}

func (h *userHandler) Post(w http.ResponseWriter, r *http.Request) {
	var u models.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	password, err := hashPassword(u.Password)
	if err != nil {
		log.Printf("Error hashing password: %v\n", err.Error())
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	u.Password = string(password)

	id, err := uuid.NewV4()
	if err != nil {
		log.Printf("Error generating UUID: %v\n", err.Error())
		http.Error(w, "Error generating UUID", http.StatusInternalServerError)
		return
	}
	u.Id = id.String()

	s, err := generateSession(u.Id)
	if err != nil {
		log.Printf("Error generating session ID: %v\n", err.Error())
		http.Error(w, "Error generating session ID", http.StatusInternalServerError)
		return
	}

	if err = h.db.CreateUser(u, s); err != nil {
		handleDbErr(err, w)
		return
	}
	storeSessionCookie(w, s)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *userHandler) LoginPost(w http.ResponseWriter, r *http.Request) {
	var id models.Identification
	err := json.NewDecoder(r.Body).Decode(&id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u, err := h.db.GetUserFromUsername(id.Username)
	if err != nil {
		if _, ok := err.(*database.NotFoundError); ok {
			http.Error(w, "Username or password does not match", http.StatusUnauthorized)
			return
		}
		handleDbErr(err, w)
		return
	}

	if !passwordMatches(u.Password, id.Password) {
		http.Error(w, "Username or password does not match", http.StatusUnauthorized)
		return
	}

	s, err := generateSession(u.Id)
	if err != nil {
		log.Printf("Error generating session ID: %v\n", err.Error())
		http.Error(w, "Error generating session ID", http.StatusInternalServerError)
		return
	}

	if err = h.db.CreateSession(s); err != nil {
		handleDbErr(err, w)
		return
	}
	storeSessionCookie(w, s)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *userHandler) LogoutPost(w http.ResponseWriter, r *http.Request) {
	sessionId, err := getSessionId(r)
	if err != nil {
		http.Error(w, "Not logged in", http.StatusForbidden)
		return
	}

	if err = h.db.DeleteSession(sessionId); err != nil {
		handleDbErr(err, w)
		return
	}
	deleteSessionCookie(w)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *userHandler) InvalidatePost(w http.ResponseWriter, r *http.Request) {
	sessionId, err := getSessionId(r)
	if err != nil {
		http.Error(w, "Not logged in", http.StatusForbidden)
		return
	}

	if err = h.db.DeleteOtherSessions(sessionId); err != nil {
		handleDbErr(err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *userHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	var id models.Identification
	err := json.NewDecoder(r.Body).Decode(&id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sessionId, err := getSessionId(r)
	if err != nil {
		http.Error(w, "Not logged in", http.StatusBadRequest)
		return
	}

	u, err := h.db.GetUserFromSession(sessionId)
	if err != nil {
		if _, ok := err.(*database.NotFoundError); ok {
			http.Error(w, "Username or password does not match", http.StatusUnauthorized)
			return
		}
		handleDbErr(err, w)
		return
	}

	if !passwordMatches(u.Password, id.Password) {
		http.Error(w, "Username or password does not match", http.StatusUnauthorized)
		return
	}

	if err = h.db.DeleteUser(u.Id); err != nil {
		handleDbErr(err, w)
		return
	}
	deleteSessionCookie(w)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *userHandler) AuthenticatedGet(w http.ResponseWriter, r *http.Request) {
	sessionId, err := getSessionId(r)
	if err != nil {
		http.Error(w, "Not logged in", http.StatusBadRequest)
		return
	}

	if _, err = h.db.GetUserFromSession(sessionId); err != nil {
		handleDbErr(err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
