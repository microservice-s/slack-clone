package users

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const defaultAddr = "127.0.0.1:27017"

// MongoStore is an implementation of UserStore
// backed by a mongo database
type MongoStore struct {
	Session        *mgo.Session
	DatabaseName   string
	CollectionName string
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
	return &MongoStore{
		Session:        session,
		DatabaseName:   databaseName,
		CollectionName: "users",
	}, nil
}

//GetAll returns all users
func (ms *MongoStore) GetAll() ([]*User, error) {
	// create a new slice of pointers to user structs
	users := []*User{}
	// return all into the provided slice
	err := ms.Session.DB(ms.DatabaseName).C(ms.CollectionName).Find(nil).All(&users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

//GetByID returns the User with the given ID
func (ms *MongoStore) GetByID(id interface{}) (*User, error) {
	// type assert that the given id is a string and convert to bson
	if sID, ok := id.(string); ok {
		id = bson.ObjectIdHex(sID)
	}
	// create empty user struct
	// and store the result of the query to it
	user := &User{}
	err := ms.Session.DB(ms.DatabaseName).C(ms.CollectionName).FindId(id).One(user)
	return user, err
}

//GetByEmail returns the User with the given email
func (ms *MongoStore) GetByEmail(email string) (*User, error) {
	// create empty user struct
	// and store the result of the query to it
	user := &User{}
	err := ms.Session.DB(ms.DatabaseName).C(ms.CollectionName).Find(bson.M{"email": email}).One(user)
	return user, err
}

//GetByUserName returns the User with the given user name
func (ms *MongoStore) GetByUserName(name string) (*User, error) {
	// create empty user struct
	// and store the result of the query to it
	user := &User{}
	err := ms.Session.DB(ms.DatabaseName).C(ms.CollectionName).Find(bson.M{"username": name}).One(user)
	return user, err
}

//Insert inserts a new NewUser into the store
//and returns a User with a newly-assigned ID
func (ms *MongoStore) Insert(newUser *NewUser) (*User, error) {
	user, err := newUser.ToUser()
	if err != nil {
		return nil, err
	}
	user.ID = bson.NewObjectId()
	// write to the database/collection
	err = ms.Session.DB(ms.DatabaseName).C(ms.CollectionName).Insert(user)
	return user, err
}

//Update applies UserUpdates to the currentUser
func (ms *MongoStore) Update(updates *UserUpdates, currentuser *User) error {
	if sID, ok := currentuser.ID.(string); ok {
		currentuser.ID = bson.ObjectIdHex(sID)
	}
	col := ms.Session.DB(ms.DatabaseName).C(ms.CollectionName)
	bUpdates := bson.M{"$set": updates}
	return col.UpdateId(currentuser.ID, bUpdates)
}

// DeleteByID deletes a user from the db by id
func (ms *MongoStore) DeleteByID(id interface{}) error {
	// type assert that the given id is a string and convert to bson
	if sID, ok := id.(string); ok {
		id = bson.ObjectIdHex(sID)
	}
	err := ms.Session.DB(ms.DatabaseName).C(ms.CollectionName).RemoveId(id)
	return err
}
