package models

import (
	"time"
)

const (
    UPVOTE = 1
    DOWNVOTE = -1
)

type Vote struct {
	Id     string    `json:"id"`
	PostId string    `json:"postId"`
	Time   time.Time `json:"time"`
	Vote   int       `json:"vote"`
}
