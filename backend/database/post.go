package database

import (
	"database/sql"

	"github.com/jbrunsting/transient/backend/models"
)

type postHandler struct {
	db *sql.DB
}

func (h *postHandler) CreatePost(p models.Post) error {
	_, err := h.db.Exec(`
	INSERT INTO Posts (id, postId, time, title, content, postUrl, imageUrl)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`, p.Id, p.PostId, p.Time, p.Title, p.Content, p.PostUrl, p.ImageUrl)
	if err != nil {
		return formatError(err, "post", "creating post")
	}

	return nil
}

func (h *postHandler) GetPosts(id string) ([]models.Post, error) {
	posts := []models.Post{}

	rows, err := h.db.Query(`
	SELECT id, postId, time, title, content, postUrl, imageUrl
	FROM Posts WHERE id = $1 ORDER BY time DESC`, id)
	if err != nil {
		return posts, formatError(err, "post", "creating post")
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		var content sql.NullString
		var postUrl sql.NullString
		var imageUrl sql.NullString
		if err = rows.Scan(&post.Id, &post.PostId, &post.Time, &post.Title, &content, &postUrl, &imageUrl); err != nil {
			break
		}

		post.Content = content.String
		post.PostUrl = postUrl.String
		post.ImageUrl = imageUrl.String
		posts = append(posts, post)
	}

	if rows.Err() != nil {
		err = rows.Err()
	}

	if err != nil {
		return posts, &UnexpectedError{
			Action:        "parsing user",
			InternalError: err.Error(),
		}
	}

	return posts, nil
}
