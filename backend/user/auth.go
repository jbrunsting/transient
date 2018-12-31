package user

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	sessionIdLength = 32
	bcryptCost      = 10
	timeFormat      = time.RFC3339
	expiryMinutes   = 86400
)

type session struct {
	SessionId string    `json:"sessionId"`
	Expiry    time.Time `json:"expiry"`
}

func hashPassword(password string) (string, error) {
	p, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	return string(p), err
}

func passwordMatches(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func generateSessionId() (string, error) {
	b := make([]byte, sessionIdLength)
	_, err := rand.Read(b)
	id := base64.URLEncoding.EncodeToString(b)[:sessionIdLength]
	return id, err
}

func generateSession() (session, error) {
	var s session
	var err error

	s.SessionId, err = generateSessionId()
	if err != nil {
		return s, err
	}

	s.Expiry = time.Now()
	s.Expiry = s.Expiry.Add(expiryMinutes * time.Minute)

	return s, nil
}

func (h *UserHandler) AuthHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := ""
		sessionId := ""
		for _, cookie := range r.Cookies() {
			if cookie.Name == usernameCookie {
				if username == "" {
					username = cookie.Value
				} else if username != cookie.Value {
					// Only allow one username cookie so we don't authenticate
					// with one username, and then use a different username
					// from the other cookie in a future handler
					http.Error(w, "Two 'username' cookies found", http.StatusBadRequest)
					return
				}
			} else if cookie.Name == sessionIdCookie {
				sessionId = cookie.Value
			}
		}

		if username == "" {
			http.Error(w, "No 'username' cookie set", http.StatusBadRequest)
			return
		} else if sessionId == "" {
			http.Error(w, "No 'sessionId' cookie set", http.StatusBadRequest)
			return
		}

		u, httpErr := h.getUser(username)
		if httpErr != nil {
			http.Error(w, httpErr.Error(), httpErr.Code)
			return
		}

		for _, s := range u.Sessions {
			if s.SessionId == sessionId {
				next(w, r)
				return
			}
		}

		http.Error(w, "Invalid session ID", http.StatusUnauthorized)
	}
}
