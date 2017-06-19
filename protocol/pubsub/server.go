package pubsub

import (
	"context"
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
	chans    ChanColl
	upgrader websocket.Upgrader
	router   *Router
}

// NewServer create pubsub http handler
func NewServer(db *gorm.DB, coll ChanColl) *Server {
	router := NewRouter()
	router.Add("crud", "rooms", "create", createRoom)
	router.Add("crud", "rooms", "list", listRooms)
	router.Add("crud", "roomActivities", "create", createRoomActivity)
	router.Add("pubsub", "", "ping", ping)
	router.Add("pubsub", "rooms", "replay", replayRoom)
	router.Add("pubsub", "rooms", "join", joinRoom)
	router.Add("pubsub", "", "whoami", whoami)
	return &Server{
		db:    db,
		chans: coll,
		upgrader: websocket.Upgrader{
			Subprotocols: []string{
				"tomatorpc-v1",
			},
		},
		router: router,
	}
}

// LoadOrNewChan load or creates a new channel for a given room id
func (srv *Server) LoadOrNewChan(id uint) Channel {
	return srv.chans.LoadOrOpen(id)
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
	user := models.User{
		Name: "Visitor",
	}

	// load user from token
	if c, err := r.Cookie("tomatorpg-token"); err != nil {
		log.Printf("error reading token from cookie: %s", err.Error())
	} else if token, err := ParseToken("abcdef", c.Value); err != nil {
		log.Printf("error parsing / validating token: %s", err.Error())
	} else {
		// get user of the id
		srv.db.Find(&user, token.Claims().Get("id"))
		if user.ID != 0 {
			log.Printf("user loaded: id=%d name=%#v", user.ID, user.Name)
		}
	}

	// session to be used and modified by procedures
	sess := &Session{
		HTTPRequest: r,
		User:        user,
		Conn:        ws,
	}

	// build common procedure context
	ctx := context.Background()
	ctx = WithDB(ctx, srv.db)
	ctx = WithServer(ctx, srv)
	ctx = WithSession(ctx, sess)

	for {

		// parse as JSON request for flexibility
		jsonRequest := lzjson.NewNode()
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
			if sess.RoomChan != nil {
				sess.RoomChan.Unsubscribe(sess.Conn)
			}
			break
		}

		// parse and execute the RPC
		var req Request
		jsonRequest.Unmarshal(&req)
		reqCtx := WithJSONReq(ctx, jsonRequest)

		// handle all routes similarly
		resp, err := srv.router.ServeRequest(reqCtx, req)
		if err != nil {
			log.Printf("error: %s", err.Error())
			ws.WriteJSON(ErrorResponseTo(req, err))
			return
		}
		ws.WriteJSON(SuccessResponseTo(req, resp))
	}
}
