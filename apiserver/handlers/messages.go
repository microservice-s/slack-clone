package handlers

import (
	"encoding/json"
	"net/http"

	"errors"

	"github.com/aethanol/challenges-aethanol/apiserver/models/messages"
	"github.com/aethanol/challenges-aethanol/apiserver/sessions"
)

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

func (ctx *Context) authorizedQuery(w http.ResponseWriter, r *http.Request) error {

	return nil
}

// ChannelsHandler allows a user to (GET) their valid channels and (POST) add a user to a channels member list
func (ctx *Context) ChannelsHandler(w http.ResponseWriter, r *http.Request) {
	// check the authentication
	state, err := ctx.authenticated(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	switch r.Method {
	// GET the channels for the authenticated user
	case "GET":
		// get the channels
		channels, err := ctx.MessageStore.GetAllUserChannels(state.User)
		if err != nil {
			http.Error(w, "error getting user channels: "+err.Error(),
				http.StatusInternalServerError)
			return
		}

		// write the channels to the user
		Respond(w, channels, contentTypeJSONUTF8)
	// POST new channels to the store
	case "POST":
		// decode the request body into a newChannel struct
		decoder := json.NewDecoder(r.Body)
		newChannel := &messages.NewChannel{}
		if err := decoder.Decode(newChannel); err != nil {
			http.Error(w, "Error: invalid JSON", http.StatusBadRequest)
			return
		}

		// validate the channel
		if err := newChannel.Validate(); err != nil {
			http.Error(w, "error validating channel: "+err.Error(),
				http.StatusBadRequest)
			return
		}

		// ensure there isn't already a channel with the same name
		// if _, err := ctx.MessageStore.GetChannelByName(newChannel.Name); err == nil {
		// 	http.Error(w, "Error: channel name already exists",
		// 		http.StatusBadRequest)
		// 	return
		// } else if err != messages.ErrChannelNotFound {
		// 	// return the internal service error if it's not the UserNotFound error << in this case not an err
		// 	http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
		// }

		// insert the channel to the store
		channel, err := ctx.MessageStore.InsertChannel(newChannel, state.User)
		if err != nil {
			http.Error(w, "error inserting channel: "+err.Error(),
				http.StatusInternalServerError)
			return
		}

		// write the channel to the user
		Respond(w, channel, contentTypeJSONUTF8)
	}
}

// SpecificChannelHandler allows a user to
func (ctx *Context) SpecificChannelHandler(w http.ResponseWriter, r *http.Request) {
	// check the authentication
	_, err := ctx.authenticated(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	switch r.Method {
	case "GET":
	case "PATCH":
	case "DELETE":
	case "LINK":
	case "UNLINK":
	}
}

// MessagesHandler
func (ctx *Context) MessagesHandler(w http.ResponseWriter, r *http.Request) {
	// check the authentication
	_, err := ctx.authenticated(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	switch r.Method {
	case "POST":
	}
}

// SpecificMessageHandler
func (ctx *Context) SpecificMessageHandler(w http.ResponseWriter, r *http.Request) {
	// check the authentication
	_, err := ctx.authenticated(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	switch r.Method {
	case "PATCH":
	case "DELETE":
	}
}
