package messages

import (
	"strconv"
	"testing"

	"github.com/aethanol/challenges-aethanol/apiserver/models/users"
)

func addUser(store *users.MongoStore, name string) (*users.User, error) {
	nu := &users.NewUser{
		Email:        name + "@test.com",
		UserName:     name,
		FirstName:    name,
		LastName:     name,
		Password:     "password",
		PasswordConf: "password",
	}

	u, err := store.Insert(nu)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func addChannel(store *MongoStore, creator *users.User, name string) (*Channel, error) {
	newChannel := &NewChannel{
		Name:        name,
		Description: name,
		Private:     true,
	}

	channel, err := store.InsertChannel(newChannel, creator)
	if err != nil {
		return nil, err
	}
	return channel, nil
}

func cleanup(userStore *users.MongoStore, messageStore *MongoStore) {
	// clean up the users and messages that were added
	userStore.Session.DB(userStore.DatabaseName).C(userStore.CollectionName).RemoveAll(nil)

	messageStore.Session.DB(userStore.DatabaseName).C(messageStore.ChannelCollection).RemoveAll(nil)
	messageStore.Session.DB(userStore.DatabaseName).C(messageStore.MessageCollection).RemoveAll(nil)
}

func TestMongoStoreInsertChannel(t *testing.T) {
	userStore, err := users.NewMongoStore(nil, "test")
	if err != nil {
		t.Fatalf("error creating new user mongo store")
	}
	messageStore, err := NewMongoStore(nil, "test")
	if err != nil {
		t.Fatalf("error creating new message mongo store")
	}

	// add a new user so we can create a new channel
	u, err := addUser(userStore, "updateChannelUser")
	if err != nil {
		t.Errorf("error adding new user: %v", err.Error())
	}

	// add a channel
	_, err = addChannel(messageStore, u, "updateChan")
	if err != nil {
		t.Errorf("error adding new channel: %v", err.Error())
	}
	_, err = addChannel(messageStore, u, "updateChan")
	if err != nil {
		t.Errorf("error adding new channel: %v", err.Error())
	}
}

func TestMongoStoreUpdateChannel(t *testing.T) {
	userStore, err := users.NewMongoStore(nil, "test")
	if err != nil {
		t.Fatalf("error creating new user mongo store")
	}
	messageStore, err := NewMongoStore(nil, "test")
	if err != nil {
		t.Fatalf("error creating new message mongo store")
	}

	// add a new user so we can create a new channel
	u, err := addUser(userStore, "updateChannelUser")
	if err != nil {
		t.Errorf("error adding new user: %v", err.Error())
	}

	// add a channel
	c, err := addChannel(messageStore, u, "updateChan")
	if err != nil {
		t.Errorf("error adding new channel: %v", err.Error())
	}
	// update a channel
	update := &ChannelUpdates{
		Name:        "UPDATEDchanName",
		Description: "UPDATEDdesc",
	}
	messageStore.UpdateChannel(update, c)

	// check that it was updated
	channel, err := messageStore.GetChannelByID(c.ID)
	if err != nil {
		t.Fatalf("error getting channel by ID: %v", err.Error())
	}
	if channel.Name != "UPDATEDchanName" {
		t.Errorf("Name field not updated, got: %v, expected: `UPDATEDchanName`", channel.Name)
	}
	if channel.Description != "UPDATEDdesc" {
		t.Errorf("Description field not updated, got: %v, expected: `UPDATEDdesc`", channel.Description)
	}
	cleanup(userStore, messageStore)
}

func TestMongoStoreDeleteChannel(t *testing.T) {
	messageStore, err := NewMongoStore(nil, "test")
	if err != nil {
		t.Fatalf("error creating new message mongo store")
	}

	userStore, err := users.NewMongoStore(nil, "test")
	if err != nil {
		t.Fatalf("error creating new user mongo store")
	}

	nu := &users.NewUser{
		Email:        "test@test.com",
		UserName:     "tester",
		FirstName:    "Test",
		LastName:     "Tester",
		Password:     "password",
		PasswordConf: "password",
	}

	u, err := userStore.Insert(nu)
	if err != nil {
		t.Errorf("error inserting user: %v", err)
	}

	if nil == u {
		t.Fatalf("nil returned from MemStore.Insert()--you probably haven't implemented NewUser.ToUser() yet")
	}
	newChannel := &NewChannel{
		Name:        "test",
		Description: "test",
		Private:     true,
	}

	channel, err := messageStore.InsertChannel(newChannel, u)
	if err != nil {
		t.Errorf("error inserting new channel %s", err.Error())
	}

	// insert a few messages
	for i := 0; i < 10; i++ {
		newMessage := &NewMessage{
			Body: strconv.Itoa(i) + " test",
		}
		messageStore.InsertMessage(newMessage, channel, u)
	}

	// delete the channel
	err = messageStore.DeleteChannel(channel)
	if err != nil {
		t.Errorf("error deleting channel: %v", err.Error())
	}

	// make sure that all the messages were deleted
	mes, err := messageStore.GetRecentMessages(channel, u, 10)
	if len(mes) != 0 {
		t.Error("all messages NOT deleted")
	}

	// delete the user from the db
	if err := userStore.DeleteByID(u.ID); err != nil {
		t.Errorf("error deleting user: %v\n", err)
	}
}

func TestMongoStoreAddUserToChannel(t *testing.T) {
	messageStore, err := NewMongoStore(nil, "test")
	if err != nil {
		t.Fatalf("error creating new message mongo store")
	}

	userStore, err := users.NewMongoStore(nil, "test")
	if err != nil {
		t.Fatalf("error creating new user mongo store")
	}

	nu := &users.NewUser{
		Email:        "test@test.com",
		UserName:     "tester",
		FirstName:    "Test",
		LastName:     "Tester",
		Password:     "password",
		PasswordConf: "password",
	}

	u, err := userStore.Insert(nu)
	if err != nil {
		t.Errorf("error inserting user: %v", err)
	}

	if nil == u {
		t.Fatalf("nil returned from MemStore.Insert()--you probably haven't implemented NewUser.ToUser() yet")
	}
	newChannel := &NewChannel{
		Name:        "test",
		Description: "test",
		Private:     true,
	}

	channel, err := messageStore.InsertChannel(newChannel, u)
	if err != nil {
		t.Errorf("error inserting new channel %s", err.Error())
	}

	err = messageStore.AddUserToChannel(u, channel)
	if err != nil {
		t.Errorf("error adding user to channel: %v", err.Error())
	}

	//delete the user and channel!
	// delete the user from the db
	if err := userStore.DeleteByID(u.ID); err != nil {
		t.Errorf("error deleting user: %v\n", err)
	}

	// delete the channel
	err = messageStore.DeleteChannel(channel)
	if err != nil {
		t.Errorf("error deleting channel: %v", err.Error())
	}

}

func TestMongoStoreRemoveUserFromChannel(t *testing.T) {
	messageStore, err := NewMongoStore(nil, "test")
	if err != nil {
		t.Fatalf("error creating new message mongo store")
	}

	userStore, err := users.NewMongoStore(nil, "test")
	if err != nil {
		t.Fatalf("error creating new user mongo store")
	}

	// add a new user to the mongo
	u, err := addUser(userStore, "removeTest")
	if err != nil {
		t.Errorf("error adding user: %v", err.Error())
	}

	// add a channel to the mongo

	c, err := addChannel(messageStore, u, "removeChan")
	if err != nil {
		t.Errorf("error adding channel: %v", err.Error())
	}

	// add the user to the channel
	err = messageStore.AddUserToChannel(u, c)
	if err != nil {
		t.Errorf("error adding user to channel: %v", err.Error())
	}

	// remove the user from the channel
	err = messageStore.RemoveUserFromChannel(u, c)
	if err != nil {
		t.Errorf("error removing user from channel: %v", err.Error())
	}

	// check that the user was removed
	c2, err := messageStore.GetAllUserChannels(u)
	if err != nil {
		t.Errorf("error getting user channels: %v", err.Error())
	}
	if len(c2) != 0 {
		t.Errorf("user wasn't removed from channel")
	}

	cleanup(userStore, messageStore)
}

func TestMongoStoreGetRecentMessages(t *testing.T) {

}
func TestMongoStoreInsertMessage(t *testing.T) {
	messageStore, err := NewMongoStore(nil, "test")
	if err != nil {
		t.Fatalf("error creating new message mongo store")
	}

	userStore, err := users.NewMongoStore(nil, "test")
	if err != nil {
		t.Fatalf("error creating new user mongo store")
	}

	// add a new user to the mongo
	u, err := addUser(userStore, "insertTest")
	if err != nil {
		t.Errorf("error adding user: %v", err.Error())
	}

	// add a channel to the mongo
	c, err := addChannel(messageStore, u, "insertChan")
	if err != nil {
		t.Errorf("error adding channel: %v", err.Error())
	}

	newMessage := &NewMessage{
		Body: "insert message test",
	}
	// add a message
	message, err := messageStore.InsertMessage(newMessage, c, u)
	if err != nil {
		t.Errorf("error inserting message: %v", err.Error())
	}

	if message.Body != "insert message test" {
		t.Errorf("Message body incorrect, got: %v, expected: `insert message test`", message.Body)
	}
	cleanup(userStore, messageStore)
}
func TestMongoStoreUpdateMessageDeleteMessage(t *testing.T) {
	messageStore, err := NewMongoStore(nil, "test")
	if err != nil {
		t.Fatalf("error creating new message mongo store")
	}

	userStore, err := users.NewMongoStore(nil, "test")
	if err != nil {
		t.Fatalf("error creating new user mongo store")
	}

	// add a new user to the mongo
	u, err := addUser(userStore, "insertTest")
	if err != nil {
		t.Errorf("error adding user: %v", err.Error())
	}

	// add a channel to the mongo
	c, err := addChannel(messageStore, u, "insertChan")
	if err != nil {
		t.Errorf("error adding channel: %v", err.Error())
	}

	newMessage := &NewMessage{
		Body: "insert message test",
	}
	// add a message
	message, err := messageStore.InsertMessage(newMessage, c, u)
	if err != nil {
		t.Errorf("error inserting message: %v", err.Error())
	}

	if message.Body != "insert message test" {
		t.Errorf("Message body incorrect, got: %v, expected: `insert message test`", message.Body)
	}

	// TEST UPDATE Message
	updates := &MessageUpdates{
		Body: "UPDATED message",
	}

	err = messageStore.UpdateMessage(updates, message)

	// TEST DELETE message
	err = messageStore.DeleteMessage(message)
	cleanup(userStore, messageStore)
}
