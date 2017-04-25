package handlers

import (
	"testing"

	"github.com/aethanol/challenges-aethanol/apiserver/models/users"
	"github.com/aethanol/challenges-aethanol/apiserver/sessions"
)

// create the handlers context for the tests
func newContext() *Context {
	return &Context{
		SessionKey:   "supersecret",
		SessionStore: sessions.NewMemStore(-1),
		UserStore:    users.NewMemStore(),
	}
}

func TestUsersHandler(t *testing.T) {
	hctx := newContext()

}

func TestSessionshandler(t *testing.T) {

}

func TestSessionsMineHandler(t *testing.T) {

}

func TestUsersMeHanlder(t *testing.T) {

}
