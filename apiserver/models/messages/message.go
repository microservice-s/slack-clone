package messages

import (
	"errors"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/aethanol/challenges-aethanol/apiserver/models/users"
)

//MessageID defines the type for message IDs
type MessageID interface{}

// Message represents a Message in the store
type Message struct {
	ID        MessageID    `json:"id" bson:"_id"`
	ChannelID ChannelID    `json:"channelID"`
	Body      string       `json:"body"`
	CreatedAt time.Time    `json:"createdAt"`
	CreatorID users.UserID `json:"creatorID"`
	EditedAt  time.Time    `json:"editedAt"`
}

// NewMessage represents a new message when created
type NewMessage struct {
	ChannelID ChannelID `json:"channelID"`
	Body      string    `json:"body"`
}

// MessageUpdates represents message updates that can be applied to a message
type MessageUpdates struct {
	Body string `json:"body"`
}

// Validate validates a new message
func (nm *NewMessage) Validate() error {
	if len(nm.Body) == 0 {
		return errors.New("Error: body is zero length")
	}

	if nm.ChannelID == nil {
		return errors.New("Error: no channel specified")
	}

	return nil
}

// ToMessage converts a NewMessage to a Message
func (nm *NewMessage) ToMessage(creator *users.User) (*Message, error) {
	// make sure that the creatorID is a bson ID
	if sID, ok := creator.ID.(string); ok {
		creator.ID = bson.ObjectIdHex(sID)
	}

	// make sure that the ChannelID is a bson ID
	if sID, ok := nm.ChannelID.(string); ok {
		nm.ChannelID = bson.ObjectIdHex(sID)
	}

	// return a new message
	// EditedAt will be null and then can be used to check to display *(edited sym)
	return &Message{
		ChannelID: nm.ChannelID,
		Body:      nm.Body,
		CreatedAt: time.Now(),
		CreatorID: creator.ID,
	}, nil
}
