package pubsub

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/tomatorpg/tomatorpg/models"
)

// RoomChannel abstract
type RoomChannel struct {
	broadcast chan Broadcast
	clients   map[*websocket.Conn]bool
	history   []Broadcast
}

// NewRoom create a new room channel
func NewRoom() *RoomChannel {
	return &RoomChannel{
		broadcast: make(chan Broadcast),
		clients:   make(map[*websocket.Conn]bool),
		history:   make([]Broadcast, 0, 1024),
	}
}

// Register the given client to the room broadcast
func (room *RoomChannel) Register(client *websocket.Conn) {
	room.clients[client] = true
}

// Replay play back the action history stack to a newly connected user
// TODO: allow to playback partially
func (room *RoomChannel) Replay(client *websocket.Conn) {
	historyCopy := make([]Broadcast, len(room.history))
	copy(historyCopy, room.history)

	for _, broadcast := range historyCopy {
		err := client.WriteJSON(broadcast)
		log.Printf("replay: %s", broadcast.Data.Message)
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
	// Send the newly received message to the broadcast channel
	broadcast := Broadcast{
		Version: "0.1",
		Entity:  "roomActivities",
		Data:    activity,
	}
	room.broadcast <- broadcast
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
