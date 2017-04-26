package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"

	"strings"

	"github.com/aethanol/challenges-aethanol/apiserver/models/users"
	"github.com/aethanol/challenges-aethanol/apiserver/passwordreset"
)

type Email struct {
	Email string `json:"email"`
}

type UpdatePassword struct {
	Token        passwordreset.ResetToken `json:"token"`
	Password     string                   `json:"password"`
	PasswordConf string                   `json:"passwordConf"`
}

func send(body string, to string) error {
	fmt.Printf("body: %v, to: %v", body, to)
	from := "resetburner8080@gmail.com"
	password := "Contact1"

	msg := "From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: Password reset information" + "\r\n\r\n" +
		body + "\r\n"

	err := smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", from, password, "smtp.gmail.com"), from, []string{to}, []byte(msg))
	if err != nil {
		return err
	}
	return nil
}

// gmail: resetburner8080@gmail.com
func (ctx *Context) ResetCodesHandler(w http.ResponseWriter, r *http.Request) {

	// The request method must be "POST"
	if r.Method != "POST" {
		http.Error(w, "request method must be POST", http.StatusMethodNotAllowed)
		return
	}

	//decode the request body into an Email struct
	decoder := json.NewDecoder(r.Body)
	email := &Email{}
	if err := decoder.Decode(email); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// Check that there is an email in the db
	// by just checking UserNotFound err, then any
	if _, err := ctx.UserStore.GetByEmail(email.Email); err != nil {
		http.Error(w, "Error: email doesn't exist in the database",
			http.StatusBadRequest)
		return
	}

	// get a new reset token and add it to the cache
	token, err := passwordreset.NewResetToken(ctx.SessionKey)
	if err != nil {
		http.Error(w, "error generating new reset token"+err.Error(),
			http.StatusInternalServerError)
		return
	}

	if err := ctx.ResetStore.Save(passwordreset.ResetEmail(email.Email), token); err != nil {
		http.Error(w, "error saving new reset token"+err.Error(),
			http.StatusInternalServerError)
		return
	}

	// send the reset email to the user
	send(string(token), email.Email)

}

// PasswordResethandler handles requests to reset user passwords
// accepts a JSON-encoded object containing the one-time use reset
// code obtained from the previous API, a new password, and a confirmation of that new password.
func (ctx *Context) PasswordResethandler(w http.ResponseWriter, r *http.Request) {
	//     PUT /v1/passwords/email-address: accepts a JSON-encoded object containing the one-time use reset
	//code obtained from the previous API, a new password, and a confirmation of that new password.
	//If the reset code is valid, and if the new password and password conifrmation fields match,
	//use the email-address from the URL to find the user account in the database and reset its password.
	//Respond with a simple confirmation message. It's up to you if you want to automatically start a new
	//authenticated session after a successful password reset: some systems do that, but others make the user
	//explicitly sign-in using the new password, just to reinfroce it.
	// The request method must be "PUT"
	if r.Method != "PUT" {
		http.Error(w, "request method must be PUT", http.StatusMethodNotAllowed)
		return
	}

	//decode the request body into an UpdatePassword struct
	decoder := json.NewDecoder(r.Body)
	reset := UpdatePassword{}
	if err := decoder.Decode(reset); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// get the email from the URL string
	email := strings.Trim(r.URL.String(), "/v1/passwords/")

	// get the token into the gToken
	var gToken passwordreset.ResetToken
	if err := ctx.ResetStore.Get(passwordreset.ResetEmail(email), &gToken); err != nil {
		http.Error(w, "error getting reset token"+err.Error(),
			http.StatusBadRequest)
		return
	}

	// check if the tokens aren't matching
	if gToken != reset.Token {
		http.Error(w, "token's don't match",
			http.StatusBadRequest)
		return
	}

	// check that the passwords are valid
	if err := users.ValidatePassword(reset.Password, reset.PasswordConf); err != nil {
		http.Error(w, "error resetting password: "+err.Error(),
			http.StatusBadRequest)
		return
	}

	// update the mongoStore with the updated user
	if err := ctx.UserStore.ResetPassword(email, reset.Password); err != nil {
		http.Error(w, "error resetting password: "+err.Error(),
			http.StatusInternalServerError)
		return
	}

	// make sure to REMOVE the reset token after we reset
	if err := ctx.ResetStore.Delete(passwordreset.ResetEmail(email)); err != nil {
		http.Error(w, "error removing token: "+err.Error(),
			http.StatusInternalServerError)
		return
	}
}
