package handlers

import (
	"github.com/aethanol/challenges-aethanol/apiserver/models/users"
	"github.com/aethanol/challenges-aethanol/apiserver/passwordreset"
	"github.com/aethanol/challenges-aethanol/apiserver/sessions"
)

// Context contains the stores for the server
type Context struct {
	SessionKey   string
	SessionStore sessions.Store
	UserStore    users.Store
	ResetStore   passwordreset.Store
	EmailPass    string
}
