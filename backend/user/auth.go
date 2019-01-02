package user

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	sessionIdCookie = "sessionId"
	sessionIdLength = 32
	bcryptCost      = 10
	timeFormat      = time.RFC3339
	expiryMinutes   = 86400
)

type session struct {
	Id        string    `json:"id"`
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

func generateSession(id string) (session, error) {
	var s session
	var err error

	s.Id = id

	s.SessionId, err = generateSessionId()
	if err != nil {
		return s, err
	}

	s.Expiry = time.Now()
	s.Expiry = s.Expiry.Add(expiryMinutes * time.Minute)

	return s, nil
}

func getSessionId(r *http.Request) (string, error) {
	for _, cookie := range r.Cookies() {
		if cookie.Name == sessionIdCookie {
			return cookie.Value, nil
		}
	}

	return "", errors.New("No session ID cookie")
}

func storeSessionCookie(w http.ResponseWriter, s session) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionIdCookie,
		Path:     "/",
		Value:    s.SessionId,
		Expires:  s.Expiry,
		HttpOnly: true,
		Domain:   "",
	})
}

func deleteSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionIdCookie,
		Path:     "/",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		Domain:   "",
	})
}
