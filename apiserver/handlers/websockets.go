package handlers

import "net/http"

//TriggerEvent triggers a new MessageEvent. This is just a handy
//way to create new events for demo purposes. In a real app, you
//would create and broacast events in response to various handler
//actions, e.g., new user sign-up, post of a new message, etc.
func (ctx *Context) TriggerEvent(w http.ResponseWriter, r *http.Request) {
	// //CORS headers to allow cross-origin requests
	// w.Header().Add("Access-Control-Allow-Origin", "*")
	// w.Header().Add("Access-Control-Request-Method", "POST")
	// w.Header().Add("Access-Control-Request-Headers", "Content-Type")

	//TODO: create a new MessageEvent with a hard-coded message
	//and the current time for CreatedAt
	//Then pass the MessageEvent to the `.Notify()` method of your notifier
	//so that the event gets broadcasted to all web socket clients
}

//WebSocketUpgradeHandler handles websocket upgrade requests
func (ctx *Context) WebSocketUpgradeHandler(w http.ResponseWriter, r *http.Request) {
	//TODO: upgrade this request to a web socket connection
	//see https://godoc.org/github.com/gorilla/websocket#hdr-Overview
	//NOTE that by default, the websocket package will reject
	//cross-origin upgrade requests, so make sure you set the
	//CheckOrigin field of the Upgrader to allow upgrades from
	//any origin.
	//See https://godoc.org/github.com/gorilla/websocket#hdr-Origin_Considerations

	//after upgrading, use the `.AddClient()` method on your notifier
	//to add the new client to your notifier's map of clients

}
