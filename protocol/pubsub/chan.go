package pubsub

import (
	"fmt"
	"io"

	"github.com/gorilla/websocket"
	"github.com/tomatorpg/tomatorpg/models"
)

// ChanColl is the abstraction of a collection of channels
type ChanColl interface {
	LoadOrOpen(id uint) Channel
	Close(id uint)
}

// Channel is the abstraction for a pubsub channel
type Channel interface {
	Subscribe(MessageWriteCloser)
	Unsubscribe(MessageWriteCloser)
	BroadcastJSON(v interface{})
}

// MessageWriter is the abstraction to writing websocket
// messages into a websocket.
type MessageWriter interface {
	// WriteMessage writes raw bytes to websocket
	// messageType is integer number defined in RFC6455
	// https://tools.ietf.org/html/rfc6455#section-11.7
	WriteMessage(messageType int, p []byte) error

	// WriteJSON encode the v value into JSON and send throught websocket
	WriteJSON(v interface{}) error
}

// MessageWriteCloser is the abstraction to writing websocket
// messages into a websocket with a close method.
type MessageWriteCloser interface {
	MessageWriter
	io.Closer
}

// WebsocketChanColl implements ChanColl for *WebsocketChannel
type WebsocketChanColl map[uint]Channel

// LoadOrOpen implements ChanColl
func (coll WebsocketChanColl) LoadOrOpen(id uint) Channel {
	if _, ok := coll[id]; !ok {
		coll[id] = NewWebsocketChan()
	}
	return coll[id]
}

// Close implements ChanColl
func (coll WebsocketChanColl) Close(id uint) {
	if _, ok := coll[id]; ok {
		delete(coll, id)
	}
}

// WebsocketChan abstract
type WebsocketChan struct {
	broadcast chan interface{}
	clients   map[*websocket.Conn]bool
}

// NewWebsocketChan create a new room channel
func NewWebsocketChan() Channel {
	wsChan := &WebsocketChan{
		broadcast: make(chan interface{}),
		clients:   make(map[*websocket.Conn]bool),
	}
	go wsChan.run()
	return wsChan
}

// Subscribe the given client to the room broadcast
func (wsChan *WebsocketChan) Subscribe(client MessageWriteCloser) {
	ws, ok := client.(*websocket.Conn)
	if !ok {
		panic(fmt.Sprintf(
			"*WebsocketChannel only allow registering *websocket.Conn, got %T(%#v)",
			client,
			client,
		))
	}
	wsChan.clients[ws] = true
}

// Unsubscribe the given client from the room broadcast
func (wsChan *WebsocketChan) Unsubscribe(client MessageWriteCloser) {
	ws, ok := client.(*websocket.Conn)
	if !ok {
		panic(fmt.Sprintf(
			"*WebsocketChannel only allow unregistering *websocket.Conn, got %T(%#v)",
			client,
			client,
		))
	}
	delete(wsChan.clients, ws)
}

// BroadcastJSON an activity to the room
func (wsChan *WebsocketChan) BroadcastJSON(v interface{}) {
	wsChan.broadcast <- v
}

func (wsChan *WebsocketChan) run() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-wsChan.broadcast

		// Send it out to every client that is currently connected
		for client := range wsChan.clients {
			err := messageTo(wsChan, client, msg)
			_ = err // TODO: record / log the error by another channel
		}
	}
}

// BroadcastActivity broadcast RoomActivity to the given channel
func BroadcastActivity(ch Channel, activity models.RoomActivity) {
	// Send the newly received message to the broadcast channel
	broadcast := Broadcast{
		Version: "0.2",
		Entity:  "roomActivities",
		Type:    "broadcast",
		Data:    activity,
	}
	ch.BroadcastJSON(broadcast)
}

func messageTo(wsChan Channel, client MessageWriteCloser, msg interface{}) (err error) {
	err = client.WriteJSON(msg)
	if err != nil {
		client.Close()
		wsChan.Unsubscribe(client)
	}
	return
}
