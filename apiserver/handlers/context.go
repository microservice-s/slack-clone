package handlers

import (
	"github.com/aethanol/challenges-aethanol/apiserver/models/messages"
	"github.com/aethanol/challenges-aethanol/apiserver/models/users"
	"github.com/aethanol/challenges-aethanol/apiserver/passwordreset"
	"github.com/aethanol/challenges-aethanol/apiserver/sessions"
	"github.com/aethanol/challenges-aethanol/apiserver/websockets"
)

// Context contains the stores for the server
type Context struct {
	SessionKey   string
	SessionStore sessions.Store
	UserStore    users.Store
	MessageStore messages.Store
	ResetStore   passwordreset.Store
	EmailPass    string
	Notifier     *websockets.Notifier
}
