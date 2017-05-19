package handlers

import (
	"net/http"

	"log"

	"github.com/gorilla/websocket"
)

// get the upgrader for websocket upgrading
// make sure you set the CheckOrigin field of the Upgrader to allow upgrades from any origin
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

//WebSocketUpgradeHandler handles websocket upgrade requests
func (ctx *Context) WebSocketUpgradeHandler(w http.ResponseWriter, r *http.Request) {

	// ensure the user is authenticated
	_, err := ctx.authenticated(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	//upgrade this request to a web socket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error adding client: %s", err.Error())
		return
	}
	//after upgrading, use the `.AddClient()` method on your notifier
	//to add the new client to your notifier's map of clients
	ctx.Notifier.AddClient(conn)

}
