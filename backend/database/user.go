package database

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"

	"github.com/jbrunsting/transient/backend/models"
)

type userHandler struct {
	db *sql.DB
}

func (h *userHandler) GetUserFromUsername(username string) (models.User, error) {
	return h.getUser("username = $1", username)
}

func (h *userHandler) GetUserFromSession(sessionId string) (models.User, error) {
	return h.getUser("sessionId = $1", sessionId)
}

func (h *userHandler) getUser(whereCondition string, whereArgs ...interface{}) (models.User, error) {
	var u models.User

	s := fmt.Sprintf(`
    SELECT Users.id, username, password, email, Sessions.sessionId, Sessions.expiry FROM Users
	LEFT JOIN Sessions ON Users.id = Sessions.id
	WHERE %v`, whereCondition)
	rows, err := h.db.Query(s, whereArgs...)
	if err != nil {
		return u, formatError(err, "user", "querying users")
	}
	defer rows.Close()

	if !rows.Next() {
		return u, &NotFoundError{"user"}
	}

	for {
		var sessionId sql.NullString
		var expiry pq.NullTime
		if err = rows.Scan(&u.Id, &u.Username, &u.Password, &u.Email, &sessionId, &expiry); err != nil {
			break
		}

		if sessionId.Valid && expiry.Valid {
			u.Sessions = append(u.Sessions, models.Session{
				SessionId: sessionId.String,
				Expiry:    expiry.Time,
			})
		}

		if !rows.Next() {
			break
		}
	}

	if rows.Err() != nil {
		err = rows.Err()
	}

	if err != nil {
		return u, &UnexpectedError{
			Action:        "parsing user",
			InternalError: err.Error(),
		}
	}

	return u, nil
}

func (h *userHandler) CreateUser(u models.User, s models.Session) error {
	tx, err := h.db.Begin()
	if err != nil {
		tx.Rollback()
		return formatError(err, "user", "starting database transaction")
	}

	_, err = tx.Exec(`
	INSERT INTO Users (id, username, password, email)
	VALUES ($1, $2, $3, $4)`, u.Id, u.Username, u.Password, u.Email)
	if err != nil {
		tx.Rollback()
		return formatError(err, "user", "creating user")
	}

	_, err = tx.Exec(`
	INSERT INTO Sessions (id, sessionId, expiry)
	VALUES ($1, $2, $3)`, u.Id, s.SessionId, s.Expiry)
	if err != nil {
		tx.Rollback()
		return formatError(err, "user", "creating session")
	}

	err = tx.Commit()
	return formatError(err, "user", "committing database transaction")
}

func (h *userHandler) CreateSession(s models.Session) error {
	_, err := h.db.Exec(`
	INSERT INTO Sessions (id, sessionId, expiry)
	VALUES ($1, $2, $3)`, s.Id, s.SessionId, s.Expiry)
	return formatError(err, "session", "creating session")
}

func (h *userHandler) DeleteUser(id string) error {
	_, err := h.db.Exec(`DELETE FROM Users WHERE id = $1`, id)
	return formatError(err, "user", "deleting user")
}

func (h *userHandler) DeleteSession(sessionId string) error {
	_, err := h.db.Exec(`DELETE FROM Sessions WHERE sessionId = $1`, sessionId)
	return formatError(err, "session", "deleting session")
}

func (h *userHandler) DeleteOtherSessions(currentSessionId string) error {
	s := `
    DELETE FROM Sessions
    WHERE sessionId <> $1 AND id = (SELECT id FROM Sessions WHERE sessionId = $1)`
	_, err := h.db.Exec(s, currentSessionId)
	return formatError(err, "session", "deleting other sessions")
}
