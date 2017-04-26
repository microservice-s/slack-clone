package passwordreset

import (
	"time"

	"github.com/aethanol/challenges-aethanol/apiserver/sessions"

	"encoding/json"

	"gopkg.in/redis.v5"
)

//redisKeyPrefix is the prefix we will use for keys
//related to session IDs. This keeps session ID keys
//separate from other keys in the shared redis key
//namespace.
const redisKeyPrefix = "token:"
const defaultAddr = "127.0.0.1:6379"

// ResetEmail represents the key for the token backed by redis
type ResetEmail string

//RedisResetStore represents a passwordreset.Store backed by redis.
type RedisResetStore struct {
	//Redis client used to talk to redis server.
	Client *redis.Client
	//Used for key expiry time on redis.
	ResetDuration time.Duration
}

//NewRedisResetStore constructs a new RedisStore, using the provided client and
//session duration. If the `client`` is nil, it will be set to redis.NewClient()
//pointing at a local redis instance. If `sessionDuration`` is negative, it will
//be set to `DefaultSessionDuration`.
func NewRedisResetStore(client *redis.Client, resetDuration time.Duration) *RedisResetStore {

	//set defaults for parameters
	//if `client` is nil, set it to a redis.NewClient()
	//pointing at a redis instance on the same machine
	//i.e., Addr is "127.0.0.1"
	if client == nil {
		roptions := redis.Options{
			Addr: defaultAddr,
		}
		client = redis.NewClient(&roptions)
	}

	//if `sessionDuration` is < 0
	//set it to DefaultSessionDuration
	if resetDuration < 0 {
		resetDuration = DefaultResetDuration
	}
	//return a new RedisStore with the Client field set to `client`
	//and the SessionDuration field set to `sessionDuration`
	return &RedisResetStore{
		Client:        client,
		ResetDuration: resetDuration,
	}
}

//Store implementation

//Save associates the provided `state` data with the provided `sid` in the store.
func (rs *RedisResetStore) Save(email ResetEmail, state interface{}) error {
	//encode the `state` into JSON
	jbuf, err := json.Marshal(state)
	if err != nil {
		return err
	}
	//use the redis client's Set() method, using `sid.getRedisKey()`
	//as the key, the JSON as the data, and the store's session duration
	//as the expiration
	status := rs.Client.Set(email.getRedisKey(), jbuf, rs.ResetDuration)
	//Set() returns a StatusCmd, which has an .Err() method that will
	//report any error that occurred; return the result of that method
	return status.Err()
}

//Get retrieves the previously saved data for the resetToken,
//and populates the `state` parameter with it.
func (rs *RedisResetStore) Get(email ResetEmail, state interface{}) error {

	// get the token from the associated email
	jbuf, err := rs.Client.Get(email.getRedisKey()).Bytes()
	if err != nil {
		if err == redis.Nil {
			return sessions.ErrStateNotFound
		}
		return err
	}

	// unmarshall the token to provided state
	err = json.Unmarshal(jbuf, state)
	if err != nil {
		return err
	}
	return nil
}

//Delete deletes all data associated with the session id from the store.
func (rs *RedisResetStore) Delete(email ResetEmail) error {
	//use the .Del() method to delete the data associated
	//with the key `sid.getRedisKey()`, and use .Err()
	//to report any errors that occurred
	err := rs.Client.Del(email.getRedisKey()).Err()
	if err != nil {
		return err
	}
	return nil
}

func (email ResetEmail) String() string {
	return string(email)
}

//returns the key to use in redis
func (email ResetEmail) getRedisKey() string {
	return redisKeyPrefix + email.String()
}
