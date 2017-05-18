package messages

import (
	"time"

	"gopkg.in/mgo.v2/bson"

	"errors"

	"github.com/aethanol/challenges-aethanol/apiserver/models/users"
)

//ChannelID defines the type for channel IDs
type ChannelID interface{}

type Channel struct {
	ID          ChannelID      `json:"id" bson:"_id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"createdAt"`
	CreatorID   users.UserID   `json:"creatorID"`
	Members     []users.UserID `json:"members"`
	Private     bool           `json:"private"`
}

type NewChannel struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Members     []users.UserID `json:"members,omitempty"`
	Private     bool           `json:"private,omitempty"`
}

type ChannelUpdates struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Validate validates a new channel
func (nc *NewChannel) Validate() error {
	// check that there was a provided name field
	if len(nc.Name) == 0 {
		return errors.New("Error: name is zero length")
	}
	return nil
}

// ToChannel converst the NewChannel to a Channel
func (nc *NewChannel) ToChannel(creator *users.User) (*Channel, error) {
	// make sure that the creatorID is a bson ID
	if sID, ok := creator.ID.(string); ok {
		creator.ID = bson.ObjectIdHex(sID)
	}
	// create a new Channel struct to convert to
	channel := &Channel{
		Name:        nc.Name,
		Description: nc.Description,
		CreatedAt:   time.Now(),
		CreatorID:   creator.ID,
		Private:     nc.Private,
	}
	// Initialize the members slice
	var members []users.UserID
	if nc.Members == nil {
		// make a new slice with the user as the only member
		members = make([]users.UserID, 0, 10)
		members = append(members, creator.ID)
	} else {
		members = nc.Members
	}
	channel.Members = members

	return channel, nil
}
