package messages

import (
	"github.com/aethanol/challenges-aethanol/apiserver/models/users"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const defaultAddr = "127.0.0.1:27017"

// MongoStore is an implementation of MessageStore
// backed by a mongo database
type MongoStore struct {
	Session           *mgo.Session
	DatabaseName      string
	MessageCollection string
	ChannelCollection string
}

// NewMongoStore returns a new MongoStore
func NewMongoStore(session *mgo.Session, databaseName string) (*MongoStore, error) {
	// set defaults for mongo session
	// if `session` is nil set it to a mgo.Dial()
	// pointing at a mongo instance on the same machine
	var err error
	if session == nil {
		session, err = mgo.Dial(defaultAddr)
	}
	if err != nil {
		return nil, err
	}
	// if there was no databasename provided
	// default to the prod database
	if databaseName == "" {
		databaseName = "production"
	}
	// return a new mongo store and no error
	store := &MongoStore{
		Session:           session,
		DatabaseName:      databaseName,
		MessageCollection: "messages",
		ChannelCollection: "channels",
	}
	// create the index for the name field
	createIndexes(store)
	return store, nil
}

func createIndexes(ms *MongoStore) {
	// ensure index on the channel name
	chIndex := mgo.Index{
		Key:        []string{"name"},
		Unique:     true,
		Background: true,
		Sparse:     true,
	}
	ms.Session.DB(ms.DatabaseName).C(ms.ChannelCollection).EnsureIndex(chIndex)
}

func authorized(collection *mgo.Collection, query bson.M) error {
	err := collection.Find(query).One(nil)
	if err == mgo.ErrNotFound {
		return ErrUnauthorized
	}
	return nil
}

// GetChannelByID returns a channel by a given ID
func (ms *MongoStore) GetChannelByID(id interface{}) (*Channel, error) {
	// convert the ID into it's object ID so we can look up in the database
	if sID, ok := id.(string); ok {
		id = bson.ObjectIdHex(sID)
	}

	// create a channel struct to store the query into
	channel := &Channel{}
	err := ms.Session.DB(ms.DatabaseName).C(ms.ChannelCollection).FindId(id).One(channel)
	// return the error and check if it's ErrNotFound
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, ErrChannelNotFound
		}
		return nil, err
	}
	return channel, nil
}

// GetChannelByName returns a channel by a given name
func (ms *MongoStore) GetChannelByName(name string) (*Channel, error) {

	// create a channel struct to store the query into
	channel := &Channel{}
	err := ms.Session.DB(ms.DatabaseName).C(ms.ChannelCollection).Find(bson.M{"name": name}).One(channel)
	// return the error and check if it's ErrNotFound
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, ErrChannelNotFound
		}
		return nil, err
	}
	return channel, nil
}

// GetAllUserChannels returns all channels a given user is allowed to see
func (ms *MongoStore) GetAllUserChannels(user *users.User) ([]*Channel, error) {
	// convert the user ID into it's object ID so we can look up in the database
	if sID, ok := user.ID.(string); ok {
		user.ID = bson.ObjectIdHex(sID)
	}
	// create a slice of pointers to channel structs
	channels := []*Channel{}
	// search the store
	err := ms.Session.DB(ms.DatabaseName).C(ms.ChannelCollection).Find(bson.M{"$or": []bson.M{bson.M{"members": user.ID}, bson.M{"private": false}}}).All(&channels)
	// return the rror and check if it's ErrNotFound
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, ErrChannelNotFound
		}
		return nil, err
	}

	return channels, nil
}

// InsertChannel inserts a new channel into the store
// returns a Channel with a newly assigned ID
func (ms *MongoStore) InsertChannel(newChannel *NewChannel, creator *users.User) (*Channel, error) {

	// convert the creator ID into it's object ID so we can look up in the database
	if sID, ok := creator.ID.(string); ok {
		creator.ID = bson.ObjectIdHex(sID)
	}

	// validate the new channel
	err := newChannel.Validate()
	if err != nil {
		return nil, err
	}

	// convert the channel by passing in the creator
	channel, err := newChannel.ToChannel(creator)
	if err != nil {
		return nil, err
	}
	// create a new objectID for the _id
	channel.ID = bson.NewObjectId()
	// insert the chanenl to the database/collection if the channel name doesn't exist
	// it has a unique index!!
	err = ms.Session.DB(ms.DatabaseName).C(ms.ChannelCollection).Insert(channel)
	if err != nil {
		if mgo.IsDup(err) {
			return nil, ErrDuplicateKey
		}
		return nil, err
	}

	return channel, nil
}

// UpdateChannel applies ChannelUpdates to a given Channel
func (ms *MongoStore) UpdateChannel(updates *ChannelUpdates, channelID interface{}, user *users.User) error {
	// convert the channel ID into it's object ID so we can look up in the database
	if sID, ok := channelID.(string); ok {
		channelID = bson.ObjectIdHex(sID)
	}

	// convert the channel ID into it's object ID so we can look up in the database
	if sID, ok := user.ID.(string); ok {
		user.ID = bson.ObjectIdHex(sID)
	}

	// check if the user is authorized to update the channel (if they are the creator)
	authQ := bson.M{"creatorid": user.ID}
	col := ms.Session.DB(ms.DatabaseName).C(ms.ChannelCollection)
	err := authorized(col, authQ)
	// return the unauth error if we got one
	if err != nil {
		return err
	}

	// otherwise update the channel
	bUpdates := bson.M{"$set": updates}
	return col.UpdateId(channelID, bUpdates)
}

// DeleteChannel deletes a channel as well as all messages posted to that channel if they are the creator
func (ms *MongoStore) DeleteChannel(channelID interface{}, user *users.User) error {
	// convert the channel ID into it's object ID so we can look up in the database
	if sID, ok := channelID.(string); ok {
		channelID = bson.ObjectIdHex(sID)
	}
	// check if the user is authorized to update the channel (if they are the creator)
	authQ := bson.M{"creatorid": user.ID}
	col := ms.Session.DB(ms.DatabaseName).C(ms.ChannelCollection)
	err := authorized(col, authQ)
	// return the unauth error if we got one
	if err != nil {
		return err
	}

	// delete all messages that are in the channel from the messages collection
	_, err = ms.Session.DB(ms.DatabaseName).C(ms.MessageCollection).RemoveAll(bson.M{"channelid": channelID})
	if err != nil {
		return err
	}
	// delete the channel from the channel collection
	err = ms.Session.DB(ms.DatabaseName).C(ms.ChannelCollection).RemoveId(channelID)
	if err != nil {
		return err
	}
	return nil
}

// AddUserToChannel adds a user to a channels Members list
func (ms *MongoStore) AddUserToChannel(userID interface{}, channelID interface{}, creatorID interface{}) error {
	// convert the user ID into it's object ID so we can look up in the database
	if sID, ok := userID.(string); ok {
		userID = bson.ObjectIdHex(sID)
	}
	// convert the channel ID into it's object ID so we can look up in the database
	if sID, ok := channelID.(string); ok {
		channelID = bson.ObjectIdHex(sID)
	}

	// convert the channel ID into it's object ID so we can look up in the database
	if sID, ok := creatorID.(string); ok {
		creatorID = bson.ObjectIdHex(sID)
	}

	// check the authorization of the creator if they are the creator OR if the channel is public
	authQ := bson.M{"$and": []bson.M{bson.M{"_id": channelID}, bson.M{"$or": []bson.M{bson.M{"creatorid": creatorID}, bson.M{"private": false}}}}}
	col := ms.Session.DB(ms.DatabaseName).C(ms.ChannelCollection)
	err := authorized(col, authQ)
	// return the unauth error if we got one
	if err != nil {
		return err
	}
	// upsert the user to the array in the mongostore (will only add a user if they aren't in the array already!)
	_, err = ms.Session.DB(ms.DatabaseName).C(ms.ChannelCollection).UpsertId(channelID, bson.M{"$addToSet": bson.M{"members": userID}})
	if err != nil {
		return err
	}
	return nil
}

// RemoveUserFromChannel deletes a user from a Channels member list
func (ms *MongoStore) RemoveUserFromChannel(userID interface{}, channelID interface{}, creatorID interface{}) error {
	// convert the user ID into it's object ID so we can look up in the database
	if sID, ok := userID.(string); ok {
		userID = bson.ObjectIdHex(sID)
	}
	// convert the channel ID into it's object ID so we can look up in the database
	if sID, ok := channelID.(string); ok {
		channelID = bson.ObjectIdHex(sID)
	}
	// convert the channel ID into it's object ID so we can look up in the database
	if sID, ok := creatorID.(string); ok {
		creatorID = bson.ObjectIdHex(sID)
	}

	// check the authorization of the creator if they are the creator OR if the channel is public
	authQ := bson.M{"$and": []bson.M{bson.M{"_id": channelID}, bson.M{"$or": []bson.M{bson.M{"creatorid": userID}, bson.M{"private": false}}}}}
	col := ms.Session.DB(ms.DatabaseName).C(ms.ChannelCollection)
	err := authorized(col, authQ)

	// pull the user from the list of members
	err = ms.Session.DB(ms.DatabaseName).C(ms.ChannelCollection).UpdateId(channelID, bson.M{"$pull": bson.M{"members": userID}})
	if err != nil {
		return err
	}
	return nil
}

// GetRecentMessages gets the most recent N messages
// posted to a particular channel if it is public or the user is a member
func (ms *MongoStore) GetRecentMessages(channelID interface{}, user *users.User, N int) ([]*Message, error) {
	// convert the channel ID into it's object ID so we can look up in the database
	if sID, ok := channelID.(string); ok {
		channelID = bson.ObjectIdHex(sID)
	}

	// convert the user ID into it's object ID so we can look up in the database
	if sID, ok := user.ID.(string); ok {
		user.ID = bson.ObjectIdHex(sID)
	}

	// var result []struct {
	// 	ID          ChannelID      `json:"id" bson:"_id"`
	// 	Name        string         `json:"name"`
	// 	Description string         `json:"description"`
	// 	CreatedAt   time.Time      `json:"createdAt"`
	// 	CreatorID   users.UserID   `json:"creatorID"`
	// 	Members     []users.UserID `json:"members"`
	// 	Private     bool           `json:"private"`
	// 	Messages    []Message      `json:"messages"`
	// }
	// query mongo for the messages for the given channel and where the user is a member
	col := ms.Session.DB(ms.DatabaseName).C(ms.ChannelCollection)
	// check if the user is a member of the channel OR if it is public
	err := col.Find(bson.M{"$or": []bson.M{bson.M{"members": user.ID}, bson.M{"private": false}}}).One(nil)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, ErrUnauthorized
		}
		return nil, err
	}
	messages := []*Message{}
	col = ms.Session.DB(ms.DatabaseName).C(ms.MessageCollection)
	col.Find(bson.M{"channelid": channelID}).Sort("-createdat").Limit(N).All(&messages)

	// pipe := col.Pipe([]bson.M{{"$match": bson.M{"_id": channel.ID}},
	// 	bson.M{"$match": bson.M{"$or": []bson.M{bson.M{"members": user.ID}, bson.M{"private": false}}}},
	// 	bson.M{"$lookup": bson.M{
	// 		"from":         "messages",
	// 		"localField":   "_id",
	// 		"foreignField": "channelid",
	// 		"as":           "messages"}}})
	// err := pipe.Iter().All(&result)
	// query := col.Find(bson.M{"$and": []bson.M{bson.M{"channelid": channel.ID}, bson.M{"$or": []bson.M{bson.M{"members": user.ID}, bson.M{"private": false}}}}})
	// err := query.Sort("-createdat").Limit(N).All(&messages)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, ErrMessageNotFound
		}
		return nil, err
	}

	return messages, nil
}

// InsertMessage adds a new message to the database
func (ms *MongoStore) InsertMessage(newMessage *NewMessage, creator *users.User) (*Message, error) {

	// convert the creator ID into it's object ID so we can look up in the database
	if sID, ok := creator.ID.(string); ok {
		creator.ID = bson.ObjectIdHex(sID)
	}

	// validate the new message
	err := newMessage.Validate()
	if err != nil {
		return nil, err
	}

	// convert the message by passing the creator and channel
	message, err := newMessage.ToMessage(creator)
	if err != nil {
		return nil, err
	}

	// check that the user is a member of the channel that they are trying to post to
	cCol := ms.Session.DB(ms.DatabaseName).C(ms.ChannelCollection)
	authQ := bson.M{"$and": []bson.M{bson.M{"_id": message.ChannelID}, bson.M{"members": creator.ID}}}
	err = authorized(cCol, authQ)
	if err != nil {
		return nil, err
	}

	message.ID = bson.NewObjectId()
	// insert the message to the database
	err = ms.Session.DB(ms.DatabaseName).C(ms.MessageCollection).Insert(message)
	if err != nil {
		return nil, err
	}
	return message, nil
}

// GetMessageByID returns a message by a given ID
func (ms *MongoStore) GetMessageByID(id interface{}) (*Message, error) {
	// convert the ID into it's object ID so we can look up in the database
	if sID, ok := id.(string); ok {
		id = bson.ObjectIdHex(sID)
	}

	// create a channel struct to store the query into
	message := &Message{}
	err := ms.Session.DB(ms.DatabaseName).C(ms.MessageCollection).FindId(id).One(message)
	// return the error and check if it's ErrNotFound
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, ErrChannelNotFound
		}
		return nil, err
	}
	return message, nil
}

// UpdateMessage applies MessageUpdates to a given Message
func (ms *MongoStore) UpdateMessage(updates *MessageUpdates, messageID interface{}, user *users.User) error {
	// convert the message ID into it's object ID so we can look up in the database
	if sID, ok := messageID.(string); ok {
		messageID = bson.ObjectIdHex(sID)
	}
	// convert the user ID into it's object ID so we can look up in the database
	if sID, ok := user.ID.(string); ok {
		user.ID = bson.ObjectIdHex(sID)
	}
	// check if the user is authorized to update the message (if they are the creator)
	authQ := bson.M{"creatorid": user.ID}
	col := ms.Session.DB(ms.DatabaseName).C(ms.MessageCollection)
	err := authorized(col, authQ)
	// return the unauth error if we got one
	if err != nil {
		return err
	}

	// otherwise apply the updates
	bUpdates := bson.M{"$set": updates}
	return col.UpdateId(messageID, bUpdates)
}

//DeleteMessage removes a message from the store
func (ms *MongoStore) DeleteMessage(messageID interface{}, user *users.User) error {
	//convert the message ID into it's object ID so we can look up in the database
	if sID, ok := messageID.(string); ok {
		messageID = bson.ObjectIdHex(sID)
	}

	// check if the user is authorized to delete the message (if they are the creator)
	authQ := bson.M{"creatorid": user.ID}
	col := ms.Session.DB(ms.DatabaseName).C(ms.MessageCollection)
	err := authorized(col, authQ)
	// return the unauth error if we got one
	if err != nil {
		return err
	}

	// delete it by it's id
	return ms.Session.DB(ms.DatabaseName).C(ms.MessageCollection).RemoveId(messageID)
}
