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
	SELECT Posts.id, Users.username, postId, time, title, content, postUrl, imageUrl
	FROM Posts
	INNER JOIN Users on Users.id = Posts.id
    WHERE Posts.id = $1 ORDER BY time DESC`, id)
	if err != nil {
		return posts, formatError(err, "post", "getting posts")
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		var content sql.NullString
		var postUrl sql.NullString
		var imageUrl sql.NullString
		if err = rows.Scan(&post.Id, &post.Username, &post.PostId, &post.Time, &post.Title, &content, &postUrl, &imageUrl); err != nil {
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

func (h *postHandler) GetPost(postId string) (models.Post, error) {
	var post models.Post

	rows, err := h.db.Query(`
	SELECT Posts.id, Users.username, postId, time, title, content, postUrl, imageUrl
	FROM Posts
	INNER JOIN Users ON Users.id = Posts.id
    WHERE postId = $1
    `, postId)
	if err != nil {
		return post, formatError(err, "post", "retrieving post")
	}
	defer rows.Close()

	if rows.Next() {
		var content sql.NullString
		var postUrl sql.NullString
		var imageUrl sql.NullString
		err = rows.Scan(&post.Id, &post.Username, &post.PostId, &post.Time, &post.Title, &content, &postUrl, &imageUrl)
		post.Content = content.String
		post.PostUrl = postUrl.String
		post.ImageUrl = imageUrl.String
	} else if rows.Err() == nil {
		return post, &NotFoundError{"post"}
	}

	if rows.Err() != nil {
		err = rows.Err()
	}

	if err != nil {
		return post, &UnexpectedError{
			Action:        "parsing user",
			InternalError: err.Error(),
		}
	}

	return post, nil
}

func (h *postHandler) DeletePost(postId string) error {
	_, err := h.db.Exec(`DELETE FROM Posts WHERE postId = $1`, postId)
	if err != nil {
		return formatError(err, "post", "deleting post")
	}

	return nil
}

func (h *postHandler) GetFollowingsPosts(id string) ([]models.Post, error) {
	posts := []models.Post{}

	rows, err := h.db.Query(`
	SELECT Posts.id, Users.username, postId, time, title, content, postUrl, imageUrl FROM Posts
	INNER JOIN Followings on Followings.followingId = Posts.id
	INNER JOIN Users on Users.id = Posts.id
	WHERE Followings.id = $1
	ORDER BY time DESC`, id)
	if err != nil {
		return posts, formatError(err, "post", "getting posts")
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		var content sql.NullString
		var postUrl sql.NullString
		var imageUrl sql.NullString
		if err = rows.Scan(&post.Id, &post.Username, &post.PostId, &post.Time, &post.Title, &content, &postUrl, &imageUrl); err != nil {
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
