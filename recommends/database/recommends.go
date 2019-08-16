package database

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"

	"github.com/jbrunsting/transient/backend/models"
)

type recommendsHandler struct {
	db *sql.DB
}

func (h *recommendsHandler) GetUserFromId(id string) (models.User, error) {
	return h.getUser("Users.id = $1", id)
}

func (h *recommendsHandler) getUser(whereCondition string, whereArgs ...interface{}) (models.User, error) {
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
