package pubsub

import (
	"fmt"
	"log"
	"net/http"

	"gopkg.in/jose.v1/crypto"
	"gopkg.in/jose.v1/jws"

	"github.com/go-restit/lzjson"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	"github.com/tomatorpg/tomatorpg/models"
)

// Server implements pubsub websocket server
type Server struct {
	db       *gorm.DB
	rooms    map[uint64]*RoomChannel
	upgrader websocket.Upgrader
}

// NewServer create pubsub http handler
func NewServer(db *gorm.DB) *Server {

	rooms := make(map[uint64]*RoomChannel)

	// TODO: remove dummy room
	// dummy room
	rooms[0] = NewRoom()

	return &Server{
		db:    db,
		rooms: rooms,
		upgrader: websocket.Upgrader{
			Subprotocols: []string{
				"tomatorpc-v1",
			},
		},
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
	room := srv.rooms[0]
	room.Register(ws)
	room.Replay(ws)

	// load user from token
	var user models.User
	if c, err := r.Cookie("tomatorpg-token"); err != nil {
		log.Printf("error reading token from cookie: %s", err.Error())
	} else {
		serializedToken := []byte(c.Value)
		token, _ := jws.ParseJWT(serializedToken)
		if err = token.Validate([]byte("abcdef"), crypto.SigningMethodHS256); err != nil {
			log.Printf("error validating token: %s", err.Error())
		}

		// TODO: further validate token (e.g. expires)

		// get user of the id
		srv.db.Find(&user, token.Claims().Get("id"))
		if user.ID != 0 {
			log.Printf("user loaded: %#v", user)
		}
	}

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
			room.Unregister(ws)
			break
		}

		// parse and execute the RPC
		jsonRequest.Unmarshal(&rpc)
		switch rpc.Entity {
		case "roomActivities":
			// TODO: validate payload format
			jsonRequest.Get("payload").Unmarshal(&activity)
			activity.UserID = user.ID // enforce user session
			log.Printf("roomActivity: user-%d %s %s",
				activity.UserID, activity.Action, activity.Message)
			ws.WriteJSON(Response{
				Version: "0.1",
				ID:      rpc.ID,
				Type:    "response",
				Entity:  "roomActivity",
				Action:  "create",
				Status:  "success",
			})
			room.Broadcast(activity)
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
					Status:  "success",
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
					Status:  "success",
					Data:    rooms,
				})
			} else if rpc.Action == "join" {

				// Find the room
				idToJoin := uint(jsonRequest.Get("room_id").Int())
				roomToJoin := models.Room{}
				srv.db.Find(&roomToJoin, idToJoin)
				if roomToJoin.ID == idToJoin {

					// unregister client from old room
					room.Unregister(ws)

					// attach the client to the room
					if _, ok := srv.rooms[uint64(roomToJoin.ID)]; ok {
						log.Printf("%s joinned room %d",
							r.RemoteAddr,
							roomToJoin.ID,
						)
						room = srv.rooms[uint64(roomToJoin.ID)]
					} else {
						log.Printf("%s reactivated and joinned room %d",
							r.RemoteAddr,
							roomToJoin.ID,
						)
						room = NewRoom()
						srv.rooms[uint64(roomToJoin.ID)] = room
					}

					// register client to new room
					room.Register(ws)

					ws.WriteJSON(Response{
						Version: "0.1",
						ID:      rpc.ID,
						Type:    "response",
						Entity:  "rooms",
						Action:  "join",
						Status:  "success",
					})

					// replay message after join
					room.Replay(ws)

				} else {
					log.Printf("%s failed to join room %d",
						r.RemoteAddr,
						roomToJoin.ID,
					)
					ws.WriteJSON(Response{
						Version: "0.1",
						ID:      rpc.ID,
						Type:    "response",
						Entity:  "rooms",
						Action:  "join",
						Status:  "error",
						Err:     fmt.Errorf("room (id=%d) not found", idToJoin),
					})
				}
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
