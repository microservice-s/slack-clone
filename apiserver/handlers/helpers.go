package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/aethanol/challenges-aethanol/apiserver/events"
	"github.com/aethanol/challenges-aethanol/apiserver/sessions"
)

// Respond writes data to a responseWriter
func Respond(w http.ResponseWriter, data interface{}, contentType string) {
	// add the header and encode the data as json
	w.Header().Add(headerContentType, contentTypeJSONUTF8)
	encoder := json.NewEncoder(w)
	encoder.Encode(data)
}

func (ctx *Context) authenticated(w http.ResponseWriter, r *http.Request) (*SessionState, error) {
	// Get the session state
	state := &SessionState{}

	// get the state of the browser that is accessing their page
	_, err := sessions.GetState(r, ctx.SessionKey, ctx.SessionStore, &state)
	if err != nil {
		// http.Error(w, "error getting session state "+err.Error(),
		// 	http.StatusUnauthorized)
		return nil, errors.New(http.StatusText(http.StatusUnauthorized))
	}
	return state, nil
}

func (ctx *Context) notify(dType string, data interface{}) {
	// create a new event so we can add it to the notifications queue
	event := &events.Event{
		Type: dType,
		Data: data,
	}

	ctx.Notifier.Notify(event)
}
