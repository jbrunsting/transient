package models

import (
	"time"
)

type Post struct {
	Id       string    `json:"id"`
	Username string    `json:"username"`
	PostId   string    `json:"postId"`
	Time     time.Time `json:"time"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	PostUrl  string    `json:"postUrl"`
	ImageUrl string    `json:"imageUrl"`
}
