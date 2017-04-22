package handlers

import (
	"github.com/aethanol/challenges-aethanol/apiserver/models/users"
	"github.com/aethanol/challenges-aethanol/apiserver/sessions"
)

type Context struct {
	SessionKey   string
	SessionStore sessions.Store
	UserStore    users.Store
}
