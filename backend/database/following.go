package database

import (
	"database/sql"
)

type followingHandler struct {
	db *sql.DB
}

func (h *followingHandler) CreateFollowing(id, followingId string) error {
	_, err := h.db.Exec(`
	INSERT INTO Followings (id, followingId)
	VALUES ($1, $2)`, id, followingId)
	if err != nil {
		return formatError(err, "following", "creating following")
	}

	return nil
}

func (h *followingHandler) GetFollowings(id string) ([]string, error) {
	followings := []string{}

	rows, err := h.db.Query(`
	SELECT followingId
	FROM Followings WHERE id = $1`, id)
	if err != nil {
		return followings, formatError(err, "followings", "getting followings")
	}
	defer rows.Close()

	for rows.Next() {
		var followingId sql.NullString
		if err = rows.Scan(&followingId); err != nil {
			break
		}

		if followingId.String != "" {
			followings = append(followings, followingId.String)
		}
	}

	if rows.Err() != nil {
		err = rows.Err()
	}

	if err != nil {
		return followings, &UnexpectedError{
			Action:        "parsing followings",
			InternalError: err.Error(),
		}
	}

	return followings, nil
}

func (h *followingHandler) DeleteFollowing(id, followingId string) error {
	_, err := h.db.Exec(`DELETE FROM Followings WHERE id = $1 AND followingId = $2`, id, followingId)
	if err != nil {
		return formatError(err, "following", "deleting following")
	}

	return nil
}
