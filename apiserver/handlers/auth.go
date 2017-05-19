package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"time"

	"github.com/aethanol/challenges-aethanol/apiserver/models/users"
	"github.com/aethanol/challenges-aethanol/apiserver/sessions"
)

const headerAuthorization = "Authorization"

// UsersHandler allows new users to sign-up (POST) or returns all users (GET)
func (ctx *Context) UsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		// decode the request body into a newUser struct
		decoder := json.NewDecoder(r.Body)
		newUser := &users.NewUser{}
		if err := decoder.Decode(newUser); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		// validate the new user
		if err := newUser.Validate(); err != nil {
			http.Error(w, "error validating user: "+err.Error(),
				http.StatusBadRequest)
			return
		}

		// Ensure there isn't already a user in the UserStore with the same email address
		// by just checking UserNotFound err, then any
		if _, err := ctx.UserStore.GetByEmail(newUser.Email); err == nil {
			http.Error(w, "Error: email already exists in database",
				http.StatusBadRequest)
			return
		} else if err != users.ErrUserNotFound {
			// return the internal service error if it's not the UserNotFound error << in this case not an err
			http.Error(w, "Error:"+err.Error(), http.StatusInternalServerError)
		}

		// Ensure there isn't already a user in the UserStore with the same user name
		if _, err := ctx.UserStore.GetByUserName(newUser.UserName); err == nil {
			http.Error(w, "Error: user name already exists in database",
				http.StatusBadRequest)
			return
		} else if err != users.ErrUserNotFound {
			// return the internal service error if it's not the UserNotFound error << in this case not an err
			http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
		}

		// Insert the new user into the UserStore
		var user *users.User
		var err error
		if user, err = ctx.UserStore.Insert(newUser); err != nil {
			http.Error(w, "error inserting new user: "+err.Error(),
				http.StatusInternalServerError)
			return
		}

		// Begin a new session with the context session signing key and a new session state
		state := &SessionState{
			BeganAt:    time.Now(),
			ClientAddr: r.RemoteAddr,
			User:       user,
		}
		if _, err := sessions.BeginSession(ctx.SessionKey, ctx.SessionStore,
			state, w); err != nil {
			http.Error(w, "error beginning new session: "+err.Error(),
				http.StatusInternalServerError)
			return
		}

		// notify the clients of the new user
		ctx.notify("new user", user)
		// Respond to the client with the models.User struct encoded as a JSON object
		Respond(w, user, contentTypeJSONUTF8)
	case "GET":
		// Get all users from the UserStore and write them to the response
		// as a JSON-encoded array
		// TODO: make sure they're authenticated!
		users, err := ctx.UserStore.GetAll()
		if err != nil {
			http.Error(w, "error getting all users: "+err.Error(),
				http.StatusInternalServerError)
			return
		}
		// Respond to the user
		Respond(w, users, contentTypeJSONUTF8)
	}
}

// SessionsHandler allows existing users to sign-in
func (ctx *Context) SessionsHandler(w http.ResponseWriter, r *http.Request) {
	// The request method must be "POST"
	if r.Method != "POST" {
		http.Error(w, "request method must be POST", http.StatusMethodNotAllowed)
		return
	}
	// Decode the request body into a users.Credentials struct
	decoder := json.NewDecoder(r.Body)
	creds := &users.Credentials{}
	if err := decoder.Decode(creds); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// Get the user with the provided email from the UserStore; if not found, respond with an http.StatusUnauthorized
	user, err := ctx.UserStore.GetByEmail(creds.Email)
	if err != nil {
		http.Error(w, "error authenticating user", http.StatusUnauthorized)
		return
	}

	// TODO: implement sign in attempts ec here

	// Authenticate the user using the provided password; if that fails, respond with an http.StatusUnauthorized
	err = user.Authenticate(creds.Password)
	if err != nil {
		http.Error(w, "error authenticating user", http.StatusUnauthorized)
		return
	}
	// Begin a new session by getting the session state
	state := &SessionState{
		BeganAt:    time.Now(),
		ClientAddr: r.RemoteAddr,
		User:       user,
	}
	// Begin a new session with the context session signing key and the state
	// Adds the header authorization to the response
	if _, err := sessions.BeginSession(ctx.SessionKey, ctx.SessionStore,
		state, w); err != nil {
		http.Error(w, "error beginning new session: "+err.Error(),
			http.StatusInternalServerError)
		return
	}

	// Respond to the client with the models.User struct encoded as a JSON object
	Respond(w, user, contentTypeJSONUTF8)
}

// SessionsMineHandler allows authenticated users to sign-out
func (ctx *Context) SessionsMineHandler(w http.ResponseWriter, r *http.Request) {
	// The request method must be "DELETE"
	if r.Method != "DELETE" {
		http.Error(w, "request method must be DELETE", http.StatusMethodNotAllowed)
		return
	}
	// End the session by getting the sessionID from the request and deleting from redis
	sid, err := sessions.GetSessionID(r, ctx.SessionKey)
	if err != nil {
		http.Error(w, "error getting sessionID: "+err.Error(),
			http.StatusBadRequest)
		return
	}
	// delete the session from the store
	if err := ctx.SessionStore.Delete(sid); err != nil {
		http.Error(w, "error deleting session"+err.Error(),
			http.StatusInternalServerError)
		return
	}
	// Respond to the client with a simple message saying that the user has been signed out
	io.WriteString(w, "user signed out\n")
}

// UsersMeHanlder allows a users to get their current session state
func (ctx *Context) UsersMeHanlder(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// Get the session state
		state := &SessionState{}

		// get the state of the browser that is accessing their page
		_, err := sessions.GetState(r, ctx.SessionKey, ctx.SessionStore, &state)
		if err != nil {
			http.Error(w, "error getting session state "+err.Error(),
				http.StatusForbidden)
			return
		}

		// Respond to the client with the session state's User field, encoded as a JSON object
		Respond(w, state.User, contentTypeJSONUTF8)
	case "PATCH":
		// allow the client to set the FirstName and/or LastName fields for the currently-authenticated user.
		// Use the users.UserUpdates struct when decoding the post body, and pass that to the .Update() method
		// of the UserStore on your handler context.

		// Get the session state
		state := &SessionState{}

		// get the state of the browser that is accessing their page
		sid, err := sessions.GetState(r, ctx.SessionKey, ctx.SessionStore, &state)
		if err != nil {
			http.Error(w, "error getting session state "+err.Error(),
				http.StatusForbidden)
			return
		}

		// Decode the request body into a users.UserUpdates struct
		decoder := json.NewDecoder(r.Body)
		updates := &users.UserUpdates{}
		if err := decoder.Decode(updates); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		// update the user in the mongostore
		if err := ctx.UserStore.Update(updates, state.User); err != nil {
			http.Error(w, "error updating user "+err.Error(),
				http.StatusInternalServerError)
			return
		}

		// update the session state with the user
		state.User.FirstName = updates.FirstName
		state.User.LastName = updates.LastName
		if err := ctx.SessionStore.Save(sid, state); err != nil {
			http.Error(w, "error updating user "+err.Error(),
				http.StatusInternalServerError)
			return
		}

		// Respond to the client with a simple message saying that the user has been updated
		io.WriteString(w, "user updated\n")
	}

}
