package messages

import (
	"fmt"
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

	cleanup(userStore, messageStore)
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
	// try updating the channel
	err = messageStore.UpdateChannel(update, c.ID, u)

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
			Body:      strconv.Itoa(i) + " test",
			ChannelID: channel.ID,
		}
		messageStore.InsertMessage(newMessage, u)
	}

	// delete the channel
	err = messageStore.DeleteChannel(channel.ID, u)
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
	cleanup(userStore, messageStore)
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
		UserName:     "creator",
		FirstName:    "Test",
		LastName:     "Tester",
		Password:     "password",
		PasswordConf: "password",
	}

	creator, err := userStore.Insert(nu)
	if err != nil {
		t.Errorf("error inserting user: %v", err)
	}

	if nil == creator {
		t.Fatalf("nil returned from MemStore.Insert()--you probably haven't implemented NewUser.ToUser() yet")
	}

	// create a second user that can be added to a channel
	nu2 := &users.NewUser{
		Email:        "test@test.com",
		UserName:     "tobeadded",
		FirstName:    "Test",
		LastName:     "Tester",
		Password:     "password",
		PasswordConf: "password",
	}

	added, err := userStore.Insert(nu2)
	if err != nil {
		t.Errorf("error inserting user: %v", err)
	}

	if nil == added {
		t.Fatalf("nil returned from MemStore.Insert()--you probably haven't implemented NewUser.ToUser() yet")
	}

	newChannel := &NewChannel{
		Name:        "test",
		Description: "test",
		Private:     true,
	}

	channel, err := messageStore.InsertChannel(newChannel, creator)
	if err != nil {
		t.Errorf("error inserting new channel %s", err.Error())
	}

	// check case where a user is adding themselves to a channel and THEY AREN'T ALLOWED TO
	err = messageStore.AddUserToChannel(added.ID, channel.ID, added.ID)
	if err == nil {
		t.Errorf("error checking authorization: %v", err.Error())
	}

	// check case where an authorized user is adding another user
	err = messageStore.AddUserToChannel(added.ID, channel.ID, creator.ID)
	if err != nil {
		t.Errorf("error adding other user to private channel: %v", err.Error())
	}

	// check case where a user is adding themselves to a public channel
	newChannel2 := &NewChannel{
		Name:        "test2",
		Description: "test2",
		Private:     false,
	}

	channel, err = messageStore.InsertChannel(newChannel2, creator)
	if err != nil {
		t.Errorf("error inserting new channel %s", err.Error())
	}

	err = messageStore.AddUserToChannel(added.ID, channel.ID, added.ID)
	if err != nil {
		t.Errorf("error adding other self to public channel: %v", err.Error())
	}

	cleanup(userStore, messageStore)

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
	// err = messageStore.AddUserToChannel(u, c)
	// if err != nil {
	// 	t.Errorf("error adding user to channel: %v", err.Error())
	// }

	// remove the user from the channel
	err = messageStore.RemoveUserFromChannel(u.ID, c.ID, u.ID)
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

// func TestMongoStoreGetRecentMessages(t *testing.T) {
// 	messageStore, err := NewMongoStore(nil, "test")
// 	if err != nil {
// 		t.Fatalf("error creating new message mongo store")
// 	}

// 	userStore, err := users.NewMongoStore(nil, "test")
// 	if err != nil {
// 		t.Fatalf("error creating new user mongo store")
// 	}

// 	nu := &users.NewUser{
// 		Email:        "test@test.com",
// 		UserName:     "tester",
// 		FirstName:    "Test",
// 		LastName:     "Tester",
// 		Password:     "password",
// 		PasswordConf: "password",
// 	}

// 	u, err := userStore.Insert(nu)
// 	if err != nil {
// 		t.Errorf("error inserting user: %v", err)
// 	}

// 	if nil == u {
// 		t.Fatalf("nil returned from MemStore.Insert()--you probably haven't implemented NewUser.ToUser() yet")
// 	}
// 	newChannel := &NewChannel{
// 		Name:        "test",
// 		Description: "test",
// 		Private:     true,
// 	}

// 	channel, err := messageStore.InsertChannel(newChannel, u)
// 	if err != nil {
// 		t.Errorf("error inserting new channel %s", err.Error())
// 	}

// 	// insert a few messages
// 	for i := 0; i < 10; i++ {
// 		newMessage := &NewMessage{
// 			Body: strconv.Itoa(i) + " test",
// 		}
// 		messageStore.InsertMessage(newMessage, channel, u)
// 	}

// 	unAuthU := &users.User{
// 		ID: 1111,
// 	}
// 	messages, _ := messageStore.GetRecentMessages(channel, unAuthU, 100)
// 	fmt.Println(messages)
// 	cleanup(userStore, messageStore)
// 	// fmt.Println(err)
// }

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
		Body:      "insert message test",
		ChannelID: c.ID,
	}
	// add a message
	message, err := messageStore.InsertMessage(newMessage, u)
	if err != nil {
		t.Errorf("error inserting message: %v", err.Error())
	}

	if message.Body != "insert message test" {
		t.Errorf("Message body incorrect, got: %v, expected: `insert message test`", message.Body)
	}

	// try inserting a message with no channel
	nm := &NewMessage{
		Body: "test wit no channelID",
	}
	message, err = messageStore.InsertMessage(nm, u)
	if err == nil {
		t.Errorf("Error: message not validated properly %v", err.Error())
	}

	// try inserting to a channel you are not a part of
	// delete the user from the channel
	messageStore.RemoveUserFromChannel(u.ID, c.ID, u.ID)
	// try inserting a message with no channel
	nm = &NewMessage{
		Body:      "test inserting when not part of the channel",
		ChannelID: c.ID,
	}
	message, err = messageStore.InsertMessage(nm, u)
	if err == nil {
		t.Errorf("Error: message insert not authenticating properly %v", err.Error())
	}
	cleanup(userStore, messageStore)
}

//TestMongoStoreUPdateMessage tests updating messages
func TestMongoStoreUpdateMessage(t *testing.T) {
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
		Body:      "insert message test",
		ChannelID: c.ID,
	}
	// add a message
	message, err := messageStore.InsertMessage(newMessage, u)
	if err != nil {
		t.Errorf("error inserting message: %v", err.Error())
	}

	if message.Body != "insert message test" {
		t.Errorf("Message body incorrect, got: %v, expected: `insert message test`", message.Body)
	}

	//TEST UPDATE Message normal
	updates := &MessageUpdates{
		Body: "UPDATED message",
	}

	err = messageStore.UpdateMessage(updates, message.ID, u)
	if err != nil {
		t.Errorf("error updating message from authenticated user: %v", err.Error())
	}
	fmt.Println("got past updating normal message")

	// test updating a message that the user is not the creator
	// add a new user to the mongo
	u2, err := addUser(userStore, "notAuthed")
	if err != nil {
		t.Errorf("error adding user: %v", err.Error())
	}

	err = messageStore.UpdateMessage(updates, message.ID, u2)
	if err == nil {
		t.Errorf("error users can modify others messages: %v", err.Error())
	}

	cleanup(userStore, messageStore)
}

func TestMongoStoreDeleteMessage(t *testing.T) {
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
		Body:      "insert message test",
		ChannelID: c.ID,
	}
	// add a message
	message, err := messageStore.InsertMessage(newMessage, u)
	if err != nil {
		t.Errorf("error inserting message: %v", err.Error())
	}

	if message.Body != "insert message test" {
		t.Errorf("Message body incorrect, got: %v, expected: `insert message test`", message.Body)
	}

	// test deleting the message with the user that created it
	err = messageStore.DeleteMessage(message.ID, u)
	_, err = messageStore.GetMessageByID(message.ID)
	if err == nil {
		t.Errorf("error deleting message: %v", err.Error())
	}

	// test deleting the message with a different user
	// add a message with the first user and try to delete it with another
	message, err = messageStore.InsertMessage(newMessage, u)
	if err != nil {
		t.Errorf("error inserting message: %v", err.Error())
	}

	if message.Body != "insert message test" {
		t.Errorf("Message body incorrect, got: %v, expected: `insert message test`", message.Body)
	}

	// add a new user to the mongo
	u2, err := addUser(userStore, "otherUser")
	if err != nil {
		t.Errorf("error adding user: %v", err.Error())
	}

	err = messageStore.DeleteMessage(message.ID, u2)
	if err == nil {
		t.Errorf("error: users can delete others messages : %v", err.Error())
	}
}
