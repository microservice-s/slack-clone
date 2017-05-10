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
	GetRecentMessages(channel *Channel, user *users.User, N int) ([]*Message, error)

	// UpdateChannel applies ChannelUpdates to a given Channel
	UpdateChannel(updates *ChannelUpdates, channel *Channel, user *users.User) error

	// DeleteChannel deletes a channel as well as all messages posted to that channel if authorized
	DeleteChannel(channel *Channel, user *users.User) error

	// AddUserToChannel adds a user to a channels Members list if authorized
	AddUserToChannel(user *users.User, channel *Channel, creator *users.User) error

	// RemoveUserFromChannel deletes a user from a Channels member list
	RemoveUserFromChannel(user *users.User, channel *Channel) error

	// InsertMessage adds a message to a channel
	InsertMessage(newMessage *NewMessage, channel *Channel) (*Message, error)

	// UpdateMessage applies MessageUpdates to a given Message
	UpdateMessage(update *MessageUpdates, message *Message) error

	//DeleteMessage removes a message from the store
	DeleteMessage(message *Message) error
}
