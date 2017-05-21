package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/tomatorpg/tomatorpg/assets"
)

var broadcast = make(chan Action)            // broadcast channel
var clients = make(map[*websocket.Conn]bool) // connected clients
var port uint64
var tplIndex *template.Template
var webpackDevHost string

// Configure the upgrader
var upgrader = websocket.Upgrader{}

func init() {

	var err error

	tplBin, err := assets.Asset("html/index.html")
	if err != nil {
		log.Fatalf("cannot find index.html in assets")
	}

	tplIndex, err = template.New("index").Parse(string(tplBin))
	if err != nil {
		log.Fatal(err)
	}

	// load dot env file, if exists
	if _, err = os.Stat(".env"); err == nil {
		if err = godotenv.Load(); err != nil {
			log.Fatalf("Unable to load .env, %#v", err)
		}
	}

	// load port
	if port, err = strconv.ParseUint(os.Getenv("PORT"), 10, 16); os.Getenv("PORT") != "" && err != nil {
		log.Fatalf("Unable to parse PORT: %s", err.Error())
		return
	}
	if port == 0 {
		port = 8080
	}

	// check if in development mode
	// if so, try to load webpack dev server host
	if os.Getenv("NODE_ENV") == "development" {
		if webpackDevHost = os.Getenv("WEBPACK_DEV_SERVER_HOST"); webpackDevHost == "" {
			webpackDevHost = "http://localhost:8081" // default, if not set
		}
	}
}

// Action object
type Action struct {
	Entity    string    `json:"entity"`
	Action    string    `json:"action"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

func handlePage(scriptPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			ScriptPath string
		}{
			ScriptPath: scriptPath,
		}
		tplIndex.Execute(w, data)
	}
}

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
