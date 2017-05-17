package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan Action)            // broadcast channel

// Configure the upgrader
var upgrader = websocket.Upgrader{}

// Action object
type Action struct {
	Entity    string    `json:"entity"`
	Action    string    `json:"action"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

/**

TODO:
1. Session to store user information
2. Session ID to resume on disconnect
3. Session to be able to store any JSON object payload to use in JS
4. Create room
5. Join room by room id
6. Room listing
7. Room history save and load (only show limited row backward)
8. Room status snapshot to prevent need to read whole history to build current status

Advanced:
1. Operational Transformation?
*/

var port uint64

func handleConnections(w http.ResponseWriter, r *http.Request) {

	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	// Register our new client
	clients[ws] = true

	for {
		var msg Action
		// Read in a new Action as JSON and map it to a Action object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}

		switch msg.Action {
		case "":
			log.Printf("message: %s", msg.Message)
		case "sign_in":
			log.Printf("sign in")
		}
		// Send the newly received message to the broadcast channel
		broadcast <- msg
	}
}

func handleActions() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
		// Send it out to every client that is currently connected
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func main() {

	var err error
	if port, err = strconv.ParseUint(os.Getenv("PORT"), 10, 16); os.Getenv("PORT") != "" && err != nil {
		log.Fatalf("Unable to parse PORT: %s", err.Error())
		return
	}
	if port == 0 {
		port = 8080
	}

	// Create a simple file server
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	// Configure websocket route
	http.HandleFunc("/api.v1", handleConnections)

	// Start listening for incoming chat messages
	go handleActions()

	log.Printf("listen to port %d", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
