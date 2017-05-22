package pubsub

import (
	"log"
	"net/http"

	"github.com/go-restit/lzjson"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	"github.com/tomatorpg/tomatorpg/models"
)

// Server implements pubsub websocket server
type Server struct {
	db       *gorm.DB
	room     *RoomChannel
	upgrader websocket.Upgrader
}

// NewServer create pubsub http handler
func NewServer(db *gorm.DB) *Server {

	// TODO: move this to serve function to dynamically create and remove
	room := NewRoom()
	go room.Run()

	return &Server{
		db:   db,
		room: room,
	}
}

// ServeHTTP implements http.Handler interface
func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := srv.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	log.Printf("%s connected", r.RemoteAddr)

	// TODO: dynamically register to room on command
	srv.room.Register(ws)
	srv.room.Replay(ws)

	for {

		jsonRequest := lzjson.NewNode()
		var rpc Request
		var activity models.RoomActivity

		// parse as JSON request for flexibility
		err := ws.ReadJSON(&jsonRequest)
		if err != nil {
			switch terr := err.(type) {
			case *websocket.CloseError:
				log.Printf("%s disconnected: %d %s",
					r.RemoteAddr,
					terr.Code,
					terr.Text,
				)
			default:
				log.Printf("error: %#v", err)
			}
			srv.room.Unregister(ws)
			break
		}

		// parse and execute the RPC
		jsonRequest.Unmarshal(&rpc)
		switch rpc.Entity {
		case "roomActivities":
			// TODO: validate payload format
			jsonRequest.Get("payload").Unmarshal(&activity)
			log.Printf("roomActivity: user-%d %s %s",
				activity.UserID, activity.Action, activity.Message)
			srv.room.Do(activity)
		case "rooms":
			if rpc.Action == "create" {
				newRoom := models.Room{}
				newRoom.ID = 0 // ensure not injecting ID
				srv.db.Create(&newRoom)
				log.Printf("rooms.create: id=%d", newRoom.ID)
				ws.WriteJSON(Response{
					Version: "0.1",
					ID:      rpc.ID,
					Type:    "response",
					Entity:  "rooms",
					Action:  "create",
					Data:    newRoom,
				})
			} else if rpc.Action == "list" {
				var rooms []models.Room
				srv.db.Order("created_at desc").Find(&rooms)
				log.Printf("rooms.list length=%d", len(rooms))
				ws.WriteJSON(Response{
					Version: "0.1",
					ID:      rpc.ID,
					Type:    "response",
					Entity:  "rooms",
					Action:  "list",
					Data:    rooms,
				})
			}
		case "":
			if rpc.Action == "" {
				log.Printf("%s pinged", r.RemoteAddr)
			}
		default:
			log.Printf("rpc: %#v", rpc)
		}
	}
}
