package passwordreset

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

//InvalidSessionID represents an empty, invalid session ID
const InvalidResetToken ResetToken = ""

const tokenLength = 32
const signedLength = tokenLength + sha256.Size

//ResetToken represents a valid, digitally-signed resetToken
type ResetToken string

//ErrInvalidID is returned when an invalid session id is passed to ValidateID()
var ErrInvalidToken = errors.New("Invalid Token")

//NewSessionID creates and returns a new digitally-signed session ID,
//using `signingKey` as the HMAC signing key. An error is returned only
//if there was an error generating random bytes for the session ID
func NewResetToken(signingKey string) (ResetToken, error) {
	//make a byte slice of length `signedLength`
	buf := make([]byte, signedLength)

	//use the crypto/rand package to read `idLength`
	//random bytes into the first part of that byte slice
	//this will be our new session ID
	//if you get an error, return InvalidSessionID and
	//the error
	_, err := rand.Read(buf)
	if err != nil {
		return InvalidResetToken, err
	}

	// //use the crypto/hmac package to generate a new
	// //Message Authentication Code (MAC) for the new
	// //session ID, using the provided signing key,
	// //and put it in the last part of the byte slice
	// mac := hmac.New(sha256.New, []byte(signingKey))
	// mac.Write(buf[:signedLength])
	// sig := mac.Sum(nil)
	// copy(buf[signedLength:], sig)

	//use the encoding/base64 package to encode the
	//byte slice into a base64.URLEncoding
	//and return the result as a new SessionID
	token := ResetToken(base64.URLEncoding.EncodeToString(buf))
	return token, nil
}

//String returns a string representation of the ResetToken
func (token ResetToken) String() string {
	//just return the `token` as a string
	//HINT: https://tour.golang.org/basics/13
	return string(token)
}
