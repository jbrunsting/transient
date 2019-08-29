package database

import (
	"database/sql"
	"fmt"
	"time"

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

func (h *postHandler) GetUserPosts(id string) ([]models.Post, error) {
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
	posts, err := h.GetPosts([]string{postId})
	if err != nil {
		return models.Post{}, err
	}
	if len(posts) == 0 {
		return posts[0], &NotFoundError{"post"}
	}
	return posts[0], nil
}

func (h *postHandler) GetPosts(postIds []string) ([]models.Post, error) {
	posts := []models.Post{}

	if len(postIds) == 0 {
		return posts, nil
	}

	postIdsInterface := make([]interface{}, len(postIds))
	for i, postId := range postIds {
		postIdsInterface[i] = postId
	}

	inQuery := "$1"
	for i := 2; i < len(postIds)+1; i++ {
		inQuery += fmt.Sprintf(", $%v", i)
	}

	rows, err := h.db.Query(`
	SELECT Posts.id, Users.username, postId, time, title, content, postUrl, imageUrl
	FROM Posts
	INNER JOIN Users ON Users.id = Posts.id
    WHERE postId IN (`+inQuery+`)
	`, postIdsInterface...)
	if err != nil {
		return posts, formatError(err, "post", "retrieving posts")
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		var content sql.NullString
		var postUrl sql.NullString
		var imageUrl sql.NullString
		err = rows.Scan(&post.Id, &post.Username, &post.PostId, &post.Time, &post.Title, &content, &postUrl, &imageUrl)
		if err != nil {
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
			Action:        "parsing posts",
			InternalError: err.Error(),
		}
	}

	return posts, nil
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
			Action:        "parsing posts",
			InternalError: err.Error(),
		}
	}

	return posts, nil
}

func (h *postHandler) CreateVote(id string, postId string, vote int) error {
	_, err := h.db.Exec(`
	INSERT INTO Votes (id, postId, time, vote)
	VALUES ($1, $2, $3, $4)
    ON CONFLICT ON CONSTRAINT Votes_pkey DO UPDATE SET vote = $4`, id, postId, time.Now(), vote)
	if err != nil {
		return formatError(err, "vote", "creating vote")
	}

	return nil
}

func (h *postHandler) CreateComment(postId string, p models.Comment) error {
	_, err := h.db.Exec(`
	INSERT INTO Comments (id, postId, commentId, time, content)
    VALUES ($1, $2, $3, $4, $5)`, p.Id, postId, p.CommentId, p.Time, p.Content)
	if err != nil {
		return formatError(err, "vote", "creating vote")
	}

	return nil
}

func (h *postHandler) GetComments(postId string) ([]models.Comment, error) {
	comments := []models.Comment{}

	rows, err := h.db.Query(`
	SELECT id, commentId, time, content
	FROM Comments WHERE postId = $1
	ORDER BY time DESC`, postId)
	if err != nil {
		return comments, formatError(err, "comments", "getting comments")
	}
	defer rows.Close()

	for rows.Next() {
		var comment models.Comment
		if err = rows.Scan(&comment.Id, &comment.CommentId, &comment.Time, &comment.Content); err != nil {
			break
		}

		comment.PostId = postId
		comments = append(comments, comment)
	}

	if rows.Err() != nil {
		err = rows.Err()
	}

	if err != nil {
		return comments, &UnexpectedError{
			Action:        "parsing comments",
			InternalError: err.Error(),
		}
	}

	return comments, nil
}
