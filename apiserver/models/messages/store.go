package messages

import (
	"errors"

	"github.com/aethanol/challenges-aethanol/apiserver/models/users"
)

// ErrMessageNotFound is returned when the requested message is not found in the store
var ErrMessageNotFound = errors.New("message not found")

// ErrChannelNotFound is returned when the requested channel is not found in the store
var ErrChannelNotFound = errors.New("channel not found")

// ErrDuplicateKey is returned when a duplicate field is inserted
var ErrDuplicateKey = errors.New("duplicate key")

// ErrUnauthorized is returned when a user is unable to see a field
var ErrUnauthorized = errors.New("user unauthorized")

// Store represents an abstract store for messages.Channel and messages.Message objects.
// This interface is used by the HTTP handlers to insert new Messages, channels
// get and update. This interface can be implemented for any persistent database.
type Store interface {
	// GetAllUserChannels returns all channels a given user is allowed to see
	GetAllUserChannels(user *users.User) ([]*Channel, error)

	// InsertChannel inserts a new channel into the store
	// returns a Channel with a newly assigned ID
	InsertChannel(newChannel *NewChannel, creator *users.User) (*Channel, error)

	// GetChannelByName returns a channel by a given name
	GetChannelByName(name string) (*Channel, error)

	// GetChannelByID returns a channel by a given ID
	GetChannelByID(id interface{}) (*Channel, error)

	// GetRecentMessages gets the most recent N messages
	// posted to a particular channel if a user is authorized
	GetRecentMessages(channelID interface{}, user *users.User, N int) ([]*Message, error)

	// UpdateChannel applies ChannelUpdates to a given Channel
	UpdateChannel(updates *ChannelUpdates, channelID interface{}, user *users.User) error

	// DeleteChannel deletes a channel as well as all messages posted to that channel if authorized
	DeleteChannel(channelID interface{}, user *users.User) error

	// AddUserToChannel adds a user to a channels Members list if authorized
	AddUserToChannel(userID interface{}, channelID interface{}, creatorID interface{}) error

	// RemoveUserFromChannel deletes a user from a Channels member list
	RemoveUserFromChannel(userID interface{}, channelID interface{}, creatorID interface{}) error

	// GetMessageByID returns a message by a given ID
	GetMessageByID(id interface{}) (*Message, error)

	// InsertMessage adds a message to a channel
	InsertMessage(newMessage *NewMessage, creator *users.User) (*Message, error)

	// UpdateMessage applies MessageUpdates to a given Message
	UpdateMessage(updates *MessageUpdates, messageID interface{}, user *users.User) error

	//DeleteMessage removes a message from the store
	DeleteMessage(messageID interface{}, user *users.User) error
}
