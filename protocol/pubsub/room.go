package pubsub

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/tomatorpg/tomatorpg/models"
)

// RoomChannel abstract
type RoomChannel struct {
	Info      models.Room
	broadcast chan Broadcast
	clients   map[*websocket.Conn]bool
	history   []models.RoomActivity
}

// NewRoom create a new room channel
func NewRoom() (room *RoomChannel) {
	room = &RoomChannel{
		broadcast: make(chan Broadcast),
		clients:   make(map[*websocket.Conn]bool),
		history:   make([]models.RoomActivity, 0, 1024),
	}
	go room.Run()
	return
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

	if len(historyCopy) > 0 {
		log.Printf("replay activities to client")
		for _, activity := range historyCopy {
			err := room.MessageTo(client, Broadcast{
				Version: "0.2",
				Entity:  "roomActivities",
				Type:    "broadcast",
				Data:    activity,
			})
			if err != nil {
				// break loop
				log.Printf("error: %v", err)
				return
			}
		}
	}
}

// Unregister the given client from the room broadcast
func (room *RoomChannel) Unregister(client *websocket.Conn) {
	delete(room.clients, client)
}

// Broadcast an activity to the room
func (room *RoomChannel) Broadcast(activity models.RoomActivity) {
	// Send the newly received message to the broadcast channel
	broadcast := Broadcast{
		Version: "0.2",
		Entity:  "roomActivities",
		Type:    "broadcast",
		Data:    activity,
	}
	room.broadcast <- broadcast
}

// MessageTo specific client
func (room *RoomChannel) MessageTo(client *websocket.Conn, msg interface{}) (err error) {
	err = client.WriteJSON(msg)
	if err != nil {
		client.Close()
		room.Unregister(client)
	}
	return
}

// Run starts the main loop to handle room broadcast
func (room *RoomChannel) Run() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-room.broadcast

		// add msg to history
		room.history = append(room.history, msg.Data)

		// Send it out to every client that is currently connected
		for client := range room.clients {
			if err := room.MessageTo(client, msg); err != nil {
				log.Printf("error: %v", err)
			}
		}
	}
}
