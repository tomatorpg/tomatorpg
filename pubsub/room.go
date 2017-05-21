package pubsub

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/tomatorpg/tomatorpg/models"
)

// RoomChannel abstract
type RoomChannel struct {
	broadcast chan models.RoomActivity
	clients   map[*websocket.Conn]bool
	history   []models.RoomActivity
}

// NewRoom create a new room channel
func NewRoom() *RoomChannel {
	return &RoomChannel{
		broadcast: make(chan models.RoomActivity),
		clients:   make(map[*websocket.Conn]bool),
		history:   make([]models.RoomActivity, 0, 1024),
	}
}

// Register the given client to the room broadcast
func (room *RoomChannel) Register(client *websocket.Conn) {
	room.clients[client] = true
}

// Replay play back the action history stack to a newly connected user
// TODO: allow to playback partially
func (room *RoomChannel) Replay(client *websocket.Conn) {
	historyCopy := make([]models.RoomActivity, len(room.history))
	copy(historyCopy, room.history)

	for _, msg := range historyCopy {
		err := client.WriteJSON(msg)
		log.Printf("replay: %s", msg.Message)
		if err != nil {
			log.Printf("error: %v", err)
			return
		}
	}
}

// Unregister the given client from the room broadcast
func (room *RoomChannel) Unregister(client *websocket.Conn) {
	delete(room.clients, client)
}

// Do action given to the room
func (room *RoomChannel) Do(activity models.RoomActivity) {
	switch activity.Action {
	case "":
		// Send the newly received message to the broadcast channel
		room.broadcast <- activity
	case "sign_in":
		log.Printf("sign in")
	}
}

// Run starts the main loop to handle room broadcast
func (room *RoomChannel) Run() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-room.broadcast

		// add msg to history
		room.history = append(room.history, msg)
		log.Printf("room.history.length %d", len(room.history))

		// Send it out to every client that is currently connected
		for client := range room.clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(room.clients, client)
			}
		}
	}
}
