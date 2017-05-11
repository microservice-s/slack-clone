package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"errors"

	"strings"

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

// ChannelsHandler allows a user to (GET) their valid channels and (POST) add a user to a channels member list
func (ctx *Context) ChannelsHandler(w http.ResponseWriter, r *http.Request) {
	// check the authentication
	state, err := ctx.authenticated(w, r)
	fmt.Println(state.User)
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

		// insert the channel to the store and check if it was a duplicate channel
		channel, err := ctx.MessageStore.InsertChannel(newChannel, state.User)
		if err == messages.ErrDuplicateKey {
			http.Error(w, "error duplicate channel name: "+err.Error(),
				http.StatusBadRequest)
			return
		} else if err != nil {
			http.Error(w, "error inserting channel: "+err.Error(),
				http.StatusInternalServerError)
			return
		}

		// write the channel to the user
		Respond(w, channel, contentTypeJSONUTF8)
	}
}

// SpecificChannelHandler allows a user to GET the most recent messages of a channel, PATCH to update a channel
//
func (ctx *Context) SpecificChannelHandler(w http.ResponseWriter, r *http.Request) {
	// check the authentication
	state, err := ctx.authenticated(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	// get the channelID
	cID := strings.TrimPrefix(r.URL.Path, "/v1/channels/")
	switch r.Method {
	// get the most recent 500 recent messages of a specific channel
	case "GET":
		// get the most recent 500 messages
		messages, err := ctx.MessageStore.GetRecentMessages(cID, state.User, 500)
		if err != nil {
			http.Error(w, "Error getting messages: "+err.Error(), http.StatusForbidden)
			return
		}

		// Write the messages to the user
		Respond(w, messages, contentTypeJSONUTF8)
	// update the specified channel if the current user is the channel creator
	case "PATCH":
		// Decode the request body into a messages.ChannelUpdate struct
		decoder := json.NewDecoder(r.Body)
		updates := &messages.ChannelUpdates{}
		if err := decoder.Decode(updates); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		// update the channel with the channelID, the updates and the current user
		err := ctx.MessageStore.UpdateChannel(updates, cID, state.User)
		// if we got an error write it back to the user that they are unauthorized
		if err != nil {
			http.Error(w, "error updating channel: "+err.Error(),
				http.StatusForbidden)
			return
		}
		// write the updated channel back to the user
		channel, err := ctx.MessageStore.GetChannelByID(cID)
		if err != nil {
			http.Error(w, "error updating channel: "+err.Error(),
				http.StatusInternalServerError)
			return
		}
		// respond
		Respond(w, channel, contentTypeJSONUTF8)
	// delete the channel specified
	case "DELETE":
		// delete the channel and check the id
		err := ctx.MessageStore.DeleteChannel(cID, state.User)
		if err != nil {
			http.Error(w, "error deleting channel: "+err.Error(),
				http.StatusForbidden)
			return
		}
		// otherwise respond with a simple message that the channel was deleted
		io.WriteString(w, "channel deleted\n")
	// add a
	case "LINK":
		// check if there is a Link header in the request
		headLink := r.Header.Get("Link")
		// case where someone is adding a user to a channel
		var err error
		if len(headLink) != 0 {
			err = ctx.MessageStore.AddUserToChannel(headLink, cID, state.User.ID)
			// user is adding themselves to a channel
		} else {
			err = ctx.MessageStore.AddUserToChannel(state.User.ID, cID, state.User.ID)
		}
		if err != nil {
			http.Error(w, "error linking user: "+err.Error(),
				http.StatusForbidden)
			return
		}
		// otherwise respond with a simple message that the channel was deleted
		io.WriteString(w, "user added to channel\n")
	case "UNLINK":
		// check if there is a Link header in the request
		headLink := r.Header.Get("Link")
		// case where someone is adding a user to a channel
		var err error
		if len(headLink) != 0 {
			err = ctx.MessageStore.RemoveUserFromChannel(headLink, cID, state.User.ID)
			// user is adding themselves to a channel
		} else {
			err = ctx.MessageStore.RemoveUserFromChannel(state.User.ID, cID, state.User.ID)
		}
		if err != nil {
			http.Error(w, "error linking user: "+err.Error(),
				http.StatusForbidden)
			return
		}
		// otherwise respond with a simple message that the channel was deleted
		io.WriteString(w, "user added to channel\n")
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
