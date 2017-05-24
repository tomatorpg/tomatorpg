package pubsub

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

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
	return &Server{
		db:    db,
		rooms: make(map[uint64]*RoomChannel),
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

	// context variables
	var room *RoomChannel
	var user models.User

	// load user from token
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
			log.Printf("user loaded: id=%d name=%#v", user.ID, user.Name)
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
			log.Printf("roomActivity: user-%d %s in room-%d: %s",
				activity.UserID,
				activity.Action,
				room.Info.ID,
				activity.Message,
			)
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
			} else if rpc.Action == "replay" {
				// TODO: this is temp API, should do with CURD
				if room != nil {
					log.Printf("rooms.replay: id=%d", room.Info.ID)
					ws.WriteJSON(Response{
						Version: "0.1",
						ID:      rpc.ID,
						Type:    "response",
						Entity:  "rooms",
						Action:  "replay",
						Status:  "success",
						Data:    room.Info.ID,
					})
					room.Replay(ws)
				} else {
					ws.WriteJSON(Response{
						Version: "0.1",
						ID:      rpc.ID,
						Type:    "response",
						Entity:  "rooms",
						Action:  "replay",
						Status:  "error",
						Err:     fmt.Errorf("the session is not currently in a room"),
						Data:    room.Info.ID,
					})
					room.Replay(ws)
				}
			} else if rpc.Action == "join" {

				// Find the room
				idToJoin := uint(0)
				switch roomIDinJSON := jsonRequest.Get("room_id"); roomIDinJSON.Type() {
				case lzjson.TypeNumber:
					idToJoin = uint(roomIDinJSON.Int())
				case lzjson.TypeString:
					idParsed, _ := strconv.ParseFloat(roomIDinJSON.String(), 64)
					idToJoin = uint(idParsed)
				}

				roomToJoin := models.Room{}
				srv.db.Find(&roomToJoin, idToJoin)
				if roomToJoin.ID == idToJoin {

					// unregister client from old room
					if room != nil {
						room.Unregister(ws)
					}

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
						room.Info = roomToJoin
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
