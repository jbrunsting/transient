package models

import (
	"time"
)

type Identification struct {
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}

type Session struct {
	Id        string    `json:"id"`
	SessionId string    `json:"sessionId"`
	Expiry    time.Time `json:"expiry"`
}

type User struct {
	Id string `json:"id"`
	Identification
	Email    string    `json:"email"`
	Sessions []Session `json:"sessions,omitempty"`
}

type PasswordChange struct {
	Password    string `json:"password"`
	NewPassword string `json:"newPassword"`
}
