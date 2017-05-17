package events

// Event defines a event that is transmitted via websocket to a client
type Event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// type NewUser struct{

// }

// type NewChannel struct{

// }

// type ChannelUpdate struct{

// }

// type ChannelDelete struct {

// }

// type UserJoin struct {

// }

// type UserLeft struct {

// }

// type NewMessage struct {

// }

// type MessageUpdate struct {

// }

// type MessageDelete struct {

// }
