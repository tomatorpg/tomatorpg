package pubsub

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/tomatorpg/tomatorpg/models"
)

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

// Channel is the abstraction for a pubsub channel
type Channel interface {
	Subscribe(MessageWriter)
	Unsubscribe(MessageWriter)
	BroadcastJSON(v interface{})
}

// WebsocketChannel abstract
type WebsocketChannel struct {
	broadcast chan interface{}
	clients   map[*websocket.Conn]bool
}

// NewRoom create a new room channel
func NewRoom() Channel {
	room := &WebsocketChannel{
		broadcast: make(chan interface{}),
		clients:   make(map[*websocket.Conn]bool),
	}
	go runRoom(room)
	return room
}

// Subscribe the given client to the room broadcast
func (room *WebsocketChannel) Subscribe(client MessageWriter) {
	ws, ok := client.(*websocket.Conn)
	if !ok {
		panic(fmt.Sprintf(
			"*WebsocketChannel only allow registering *websocket.Conn, got %#v",
			client,
		))
	}
	room.clients[ws] = true
}

// Unsubscribe the given client from the room broadcast
func (room *WebsocketChannel) Unsubscribe(client MessageWriter) {
	ws, ok := client.(*websocket.Conn)
	if !ok {
		panic(fmt.Sprintf(
			"*WebsocketChannel only allow registering *websocket.Conn, got %#v",
			client,
		))
	}
	delete(room.clients, ws)
}

// BroadcastJSON an activity to the room
func (room *WebsocketChannel) BroadcastJSON(v interface{}) {
	room.broadcast <- v
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

func runRoom(room *WebsocketChannel) {
	for {
		// Grab the next message from the broadcast channel
		msg := <-room.broadcast

		// Send it out to every client that is currently connected
		for client := range room.clients {
			err := client.WriteJSON(msg)
			if err != nil {
				client.Close()
				room.Unsubscribe(client)
				log.Printf("error: %v", err)
			}
		}
	}
}
