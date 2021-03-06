package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"path"

	"github.com/aethanol/challenges-aethanol/apiserver/models/messages"
	"github.com/aethanol/challenges-aethanol/apiserver/models/users"
)

// ChannelsHandler allows a user to (GET) their valid channels and (POST) add a channel to the store
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

		// notify the clients of the new channel
		ctx.notify("new channel", channel)
		// write the channel to the user
		Respond(w, channel, contentTypeJSONUTF8)
	}
}

// SpecificChannelHandler allows a user to GET the most recent messages of a channel, PATCH to update a channel
func (ctx *Context) SpecificChannelHandler(w http.ResponseWriter, r *http.Request) {
	// check the authentication
	state, err := ctx.authenticated(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	// get the channelID
	_, cID := path.Split(r.URL.Path)
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

		// notify the clients of the updated channel
		ctx.notify("updated channel", channel)

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
	// add a user to a channel
	case "LINK":
		// check if there is a Link header in the request
		headLink := r.Header.Get("Link")
		// case where someone is adding a user to a channel
		if len(headLink) != 0 {
			if err := ctx.MessageStore.AddUserToChannel(headLink, cID, state.User.ID); err != nil {
				http.Error(w, "error linking user: "+err.Error(),
					http.StatusForbidden)
				return
			}
			// notify the clients of the new user joining the channel
			d := struct {
				UserID    users.UserID       `json:"userid"`
				ChannelID messages.ChannelID `json:"channelid"`
			}{
				headLink,
				cID,
			}
			ctx.notify("user joined", d)

			// user is adding themselves to a channel
		} else {
			if err := ctx.MessageStore.AddUserToChannel(state.User.ID, cID, state.User.ID); err != nil {
				http.Error(w, "error linking user: "+err.Error(),
					http.StatusForbidden)
				return
			}
			// notify the clients of the new user joining the channel
			d := struct {
				UserID    users.UserID       `json:"userid"`
				ChannelID messages.ChannelID `json:"channelid"`
			}{
				state.User.ID,
				cID,
			}
			ctx.notify("user joined", d)
		}

		// otherwise respond with a simple message that the channel was deleted
		io.WriteString(w, "user added to channel\n")
		// delete a user from a channel
	case "UNLINK":
		// check if there is a Link header in the request
		headLink := r.Header.Get("Link")
		// case where someone is adding a user to a channel
		if len(headLink) != 0 {
			if err := ctx.MessageStore.AddUserToChannel(headLink, cID, state.User.ID); err != nil {
				http.Error(w, "error linking user: "+err.Error(),
					http.StatusForbidden)
				return
			}
			// notify the clients of the new user joining the channel
			d := struct {
				UserID    users.UserID       `json:"userid"`
				ChannelID messages.ChannelID `json:"channelid"`
			}{
				headLink,
				cID,
			}
			ctx.notify("user left", d)

			// user is adding themselves to a channel
		} else {
			if err := ctx.MessageStore.AddUserToChannel(state.User.ID, cID, state.User.ID); err != nil {
				http.Error(w, "error linking user: "+err.Error(),
					http.StatusForbidden)
				return
			}
			// notify the clients of the new user joining the channel
			d := struct {
				UserID    users.UserID       `json:"userid"`
				ChannelID messages.ChannelID `json:"channelid"`
			}{
				state.User.ID,
				cID,
			}
			ctx.notify("user left", d)
		}
		// otherwise respond with a simple message that the channel was deleted
		io.WriteString(w, "user removed from channel\n")
	}
}

// MessagesHandler handles all requests to /v1/messages (POST) will add messages to a specified channel
func (ctx *Context) MessagesHandler(w http.ResponseWriter, r *http.Request) {
	// check the authentication
	state, err := ctx.authenticated(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	switch r.Method {
	// add a message to a specified channel
	case "POST":
		// decode the request body into a newChannel struct
		decoder := json.NewDecoder(r.Body)
		newMessage := &messages.NewMessage{}
		if err := decoder.Decode(newMessage); err != nil {
			http.Error(w, "Error: invalid JSON", http.StatusBadRequest)
			return
		}

		// validate the message
		if err := newMessage.Validate(); err != nil {
			http.Error(w, "error validating message: "+err.Error(),
				http.StatusBadRequest)
			return
		}

		// insert the message to the store and check if it was
		message, err := ctx.MessageStore.InsertMessage(newMessage, state.User)
		if err == messages.ErrUnauthorized {
			http.Error(w, "Error adding message: "+err.Error(),
				http.StatusForbidden)
			return
		} else if err != nil {
			http.Error(w, "error inserting message: "+err.Error(),
				http.StatusInternalServerError)
			return
		}

		// notify the clients of the new message
		ctx.notify("new message", message)

		// write the message to the user
		Respond(w, message, contentTypeJSONUTF8)
	}
}

// SpecificMessageHandler handles all requests made to the /v1/messages/<message-id> (PATCH) updates messages
// (DELETE) deletes messages authed
func (ctx *Context) SpecificMessageHandler(w http.ResponseWriter, r *http.Request) {
	// check the authentication
	state, err := ctx.authenticated(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	// get the message id
	_, mID := path.Split(r.URL.Path)
	switch r.Method {
	// allow a user to update a specified message if they are the creator
	case "PATCH":
		// Decode the request body into a messages.MessageUpdate struct
		decoder := json.NewDecoder(r.Body)
		updates := &messages.MessageUpdates{}
		if err := decoder.Decode(updates); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		// update the message with the channelID, the updates and the current user
		err := ctx.MessageStore.UpdateMessage(updates, mID, state.User)
		// if we got an error write it back to the user that they are unauthorized
		if err != nil {
			http.Error(w, "error updating message: "+err.Error(),
				http.StatusForbidden)
			return
		}
		// write the updated message back to the user
		message, err := ctx.MessageStore.GetMessageByID(mID)
		if err != nil {
			http.Error(w, "error updating message: "+err.Error(),
				http.StatusInternalServerError)
			return
		}

		// notify the clients of the message update
		ctx.notify("message update", message)

		// respond
		Respond(w, message, contentTypeJSONUTF8)

	// allow a user to delete a message if they are the message creator
	case "DELETE":
		// delete the message and check the id
		err := ctx.MessageStore.DeleteMessage(mID, state.User)
		if err != nil {
			http.Error(w, "error deleting message: "+err.Error(),
				http.StatusForbidden)
			return
		}

		// notify the clients of the message update
		ctx.notify("message deleted", mID)
		// otherwise respond with a simple message that the message was deleted
		io.WriteString(w, "message deleted\n")
	}
}
