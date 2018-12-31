package user

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/lib/pq"
)

const (
	usernameCookie  = "username"
	sessionIdCookie = "sessionId"
)

type UserHandler struct {
	DB *sql.DB
}

type identification struct {
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}

type user struct {
	Id string `json:"id"`
	identification
	Email    string    `json:"email"`
	Sessions []session `json:"sessions,omitempty"`
}

type httpError struct {
	Msg  string
	Code int
}

func (e httpError) Error() string {
	return e.Msg
}

func (h *UserHandler) getUser(username string) (user, *httpError) {
	var u user

	s := `
    SELECT Users.id, username, password, email, Sessions.sessionId, Sessions.expiry FROM Users
	LEFT JOIN Sessions ON Users.id = Sessions.id
	WHERE username=$1`
	rows, err := h.DB.Query(s, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return u, &httpError{Msg: "User ID not found", Code: http.StatusNotFound}
		}

		if sqlErr, ok := err.(*pq.Error); ok && sqlErr.Code.Class() == "08" {
			log.Printf("Error connecting to database: %v\n", err.Error())
			return u, &httpError{Msg: "Error connecting to database", Code: http.StatusServiceUnavailable}
		}

		log.Printf("Error getting user from database: %v\n", err.Error())
		return u, &httpError{Msg: "Error getting user from database", Code: http.StatusInternalServerError}
	}
	defer rows.Close()

	if !rows.Next() {
		if rows.Err() != nil {
			log.Printf("Error iterating over user rows: %v\n", rows.Err().Error())
			return u, &httpError{Msg: "Error iterating over user rows", Code: http.StatusInternalServerError}
		}

		return u, &httpError{Msg: "User not found", Code: http.StatusNotFound}
	}

	// Only call "rows.Next()" at the end, because we previously called it to
	// check if the result set was empty
	for {
		var s session
		err = rows.Scan(&u.Id, &u.Username, &u.Password, &u.Email, &s.SessionId, &s.Expiry)
		if err != nil {
			log.Printf("Error parsing user from database: %v\n", err.Error())
			return u, &httpError{Msg: "Error parsing user from database", Code: http.StatusInternalServerError}
		}
		if s.SessionId != "" {
			u.Sessions = append(u.Sessions, s)
		}

		if !rows.Next() {
			break
		}
	}

	if rows.Err() != nil {
		log.Printf("Error iterating over user rows: %v\n", rows.Err().Error())
		return u, &httpError{Msg: "Error iterating over user rows", Code: http.StatusInternalServerError}
	}

	return u, nil
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	username := ""
	for _, cookie := range r.Cookies() {
		if cookie.Name == usernameCookie {
			username = cookie.Value
		}
	}

    if username == "" {
        http.Error(w, "No 'username' cookie set", http.StatusBadRequest)
		return
	}

	u, httpErr := h.getUser(username)
	if httpErr != nil {
		http.Error(w, httpErr.Error(), httpErr.Code)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(u)
}

func (h *UserHandler) Post(w http.ResponseWriter, r *http.Request) {
	var u user
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

	s, err := generateSession()
	if err != nil {
		log.Printf("Error generating session ID: %v\n", err.Error())
		http.Error(w, "Error generating session ID", http.StatusInternalServerError)
		return
	}

	tx, err := h.DB.Begin()
	if err != nil {
		tx.Rollback()
		log.Printf("Error starting database transaction: %v\n", err.Error())
		http.Error(w, "Error starting database transaction", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(`
	INSERT INTO Users (id, username, password, email)
	VALUES ($1, $2, $3, $4)`, u.Id, u.Username, u.Password, u.Email)
	if err != nil {
		tx.Rollback()
		if sqlErr, ok := err.(*pq.Error); ok && (sqlErr.Code.Class() == "22" || sqlErr.Code.Class() == "23") {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			log.Printf("Error inserting user into database: %v\n", err.Error())
			http.Error(w, "Error inserting user into database", http.StatusInternalServerError)
		}
		return
	}

	_, err = tx.Exec(`
	INSERT INTO Sessions (id, sessionId, expiry)
	VALUES ($1, $2, $3)`, u.Id, s.SessionId, s.Expiry)
	if err != nil {
		tx.Rollback()
		log.Printf("Error inserting session into database: %v\n", err.Error())
		http.Error(w, "Error inserting session into database", http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Error commiting database transaction: %v\n", err.Error())
		http.Error(w, "Error commiting database transaction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(s)
}

func (h *UserHandler) LoginPost(w http.ResponseWriter, r *http.Request) {
	var id identification
	err := json.NewDecoder(r.Body).Decode(&id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u, httpErr := h.getUser(id.Username)
	if httpErr != nil {
		http.Error(w, httpErr.Error(), httpErr.Code)
		return
	}

	if !passwordMatches(u.Password, id.Password) {
		http.Error(w, "Password does not match", http.StatusUnauthorized)
		return
	}

	s, err := generateSession()
	if err != nil {
		log.Printf("Error generating session ID: %v\n", err.Error())
		http.Error(w, "Error generating session ID", http.StatusInternalServerError)
		return
	}

	_, err = h.DB.Exec(`
	INSERT INTO Sessions (id, sessionId, expiry)
	VALUES ($1, $2, $3)`, u.Id, s.SessionId, s.Expiry)
	if err != nil {
		log.Printf("Error inserting session into database: %v\n", err.Error())
		http.Error(w, "Error inserting session into database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(s)
}
