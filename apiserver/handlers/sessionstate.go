package handlers

import (
	"time"

	"github.com/aethanol/challenges-aethanol/apiserver/models/users"
)

// SessionState represents a user and their current session
type SessionState struct {
	BeganAt    time.Time
	ClientAddr string
	User       *users.User
}
