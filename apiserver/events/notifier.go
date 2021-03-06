package events

import (
	"sync"

	"encoding/json"

	"github.com/gorilla/websocket"
)

//Notifier represents a web sockets notifier
type Notifier struct {
	eventq  chan interface{}
	clients map[*websocket.Conn]bool
	sync.RWMutex
	//TODO: add other fields you might need
	//such as another channel or a mutex
	//(either would work)
	//remember that go maps ARE NOT safe for
	//concurrent access, so you must do something
	//to protect the `clients` map
}

//NewNotifier constructs a new Notifer.
func NewNotifier() *Notifier {
	//create, initialize and return a Notifier struct

	return &Notifier{
		eventq:  make(chan interface{}, 10),
		clients: make(map[*websocket.Conn]bool),
	}
}

//Start begins a loop that checks for new events
//and broadcasts them to all web socket clients.
//This function should be called on a new goroutine
//e.g., `go mynotifer.Start()`
func (n *Notifier) Start() {
	//check for new events written
	//to the `eventq` channel, and broadcast
	//them to all of the web socket clients
	for {
		event := <-n.eventq
		n.broadcast(event)
	}
}

//AddClient adds a new web socket client to the Notifer
func (n *Notifier) AddClient(client *websocket.Conn) {
	//TODO: implement this
	//But remember that this will be called from
	//an HTTP handler, and each HTTP request is
	//processed on its own goroutine, so your
	//implementation here MUST be safe for concurrent use
	n.Lock()
	n.clients[client] = true
	n.Unlock()

	//after you add the client to the map,
	//call n.readPump() on its own goroutine
	go n.readPump(client)
	//to proces all of the control messages sent
	//by the client to the server.
	//see https://godoc.org/github.com/gorilla/websocket#hdr-Control_Messages

}

//Notify will add a new event to the event queue
func (n *Notifier) Notify(event interface{}) {
	// add the `event` to the `eventq`
	n.eventq <- event
}

//readPump will read all messages (including control messages)
//send by the client and ignore them. This is necessary in order
//process the control messages. If you don't do this, the
//websocket will get stuck and start producing errors.
//see https://godoc.org/github.com/gorilla/websocket#hdr-Control_Messages
func (n *Notifier) readPump(client *websocket.Conn) {
	//TODO: implement this according to the notes in the
	//Control Message section of the Gorilla Web Socket docs:
	//https://godoc.org/github.com/gorilla/websocket#hdr-Control_Messages
	for {
		if _, _, err := client.NextReader(); err != nil {
			client.Close()
			break
		}
	}

}

//broadcast sends the event to all client as a JSON-encoded object
func (n *Notifier) broadcast(event interface{}) error {
	// Loop over all of the web socket clients in
	//n.clients and write the `event` parameter to the client
	//as a JSON-encoded object.
	//HINT: https://godoc.org/github.com/gorilla/websocket#Conn.WriteJSON
	//and for even better performance, try using a PreparedMessage:
	//https://godoc.org/github.com/gorilla/websocket#PreparedMessage
	//https://godoc.org/github.com/gorilla/websocket#Conn.WritePreparedMessage

	// marshall the event to a buffer
	buf, err := json.Marshal(event)
	if err != nil {
		return err
	}
	// create a prepared message to write to the clients
	pm, err := websocket.NewPreparedMessage(websocket.TextMessage, buf)
	if err != nil {
		return err
	}
	// write it all the clients
	n.Lock()
	defer n.Unlock()
	for c := range n.clients {
		//If you get an error while writing to a client,
		//the client has wandered off, so you should call
		//the `.Close()` method on the client, and delete
		//it from the n.clients map
		if err := c.WritePreparedMessage(pm); err != nil {
			c.Close()
			delete(n.clients, c)
		}
	}
	return nil
}
