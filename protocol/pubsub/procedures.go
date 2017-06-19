package pubsub

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/go-restit/lzjson"
	"github.com/tomatorpg/tomatorpg/models"
)

func ping(ctx context.Context, req interface{}) (resp interface{}, err error) {
	resp = "pong"
	return
}

type whoamiResp struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func whoami(ctx context.Context, req interface{}) (resp interface{}, err error) {
	sess := GetSession(ctx)
	if sess == nil {
		err = fmt.Errorf("session not found")
		return
	}
	resp = whoamiResp{
		ID:   sess.User.ID,
		Name: sess.User.Name,
	}
	return
}

func createRoom(ctx context.Context, req interface{}) (resp interface{}, err error) {
	db := GetDB(ctx)
	// TODO: read request payload for room data
	newRoom := models.Room{}
	newRoom.ID = 0 // ensure not injecting ID
	db.Create(&newRoom)
	log.Printf("rooms.create: id=%d", newRoom.ID)
	resp = newRoom
	return
}

func listRooms(ctx context.Context, req interface{}) (resp interface{}, err error) {
	db := GetDB(ctx)
	var rooms []models.Room
	db.Order("created_at desc").Find(&rooms)
	log.Printf("rooms.list length=%d", len(rooms))
	resp = rooms
	return
}

func replayRoom(ctx context.Context, req interface{}) (resp interface{}, err error) {

	db := GetDB(ctx)

	// TODO: this is temp API, should do with CURD
	//       should rewrite Replay as normal crud listing
	//       to be independent from websocket
	sess := GetSession(ctx)
	if sess == nil {
		err = fmt.Errorf("session not found")
		return
	}
	if sess.Conn == nil {
		err = fmt.Errorf("socket not found")
		return
	}
	if sess.Room == nil {
		err = fmt.Errorf("the session is not currently in a room")
		return
	}

	log.Printf("rooms.replay: id=%d", sess.RoomInfo.ID)
	resp = sess.RoomInfo.ID

	// replay history (TODO: rewrite as pure CRUD)
	historyCopy := make([]models.RoomActivity, 0, 100)
	db.Find(&historyCopy, "room_id = ?", sess.RoomInfo.ID)
	if len(historyCopy) > 0 {
		log.Printf("replay activities to client")
		for _, activity := range historyCopy {
			err := sess.Conn.WriteJSON(Broadcast{
				Version: "0.2",
				Entity:  "roomActivities",
				Type:    "broadcast",
				Data:    activity,
			})
			if err != nil {
				sess.Conn.Close()
				sess.Room.Unsubscribe(sess.Conn)
				log.Printf("error: %v", err)
				break
			}
		}
	}
	return
}

func createRoomActivity(ctx context.Context, req interface{}) (resp interface{}, err error) {
	// TODO: rewrite to pure crud
	sess := GetSession(ctx)
	if sess == nil {
		err = fmt.Errorf("session not found")
		return
	}
	if sess.Conn == nil {
		err = fmt.Errorf("socket not found in session")
		return
	}
	if sess.Room == nil {
		err = fmt.Errorf("room not found in session")
		return
	}

	// get raw json request from context
	// TODO: remove the need to do json decode here
	// TODO: validate payload format
	var activity models.RoomActivity
	jsonRequest := GetJSONReq(ctx)
	if jsonRequest == nil {
		err = fmt.Errorf("jsonRequest not found in context")
		return
	}
	jsonRequest.Get("payload").Unmarshal(&activity)
	activity.UserID = sess.User.ID     // enforce user session
	activity.RoomID = sess.RoomInfo.ID // enforce room id of sesion
	if activity.Action == "" {
		activity.Action = "message"
	}
	log.Printf("roomActivity: user-%d %s in room-%d: %s",
		activity.UserID,
		activity.Action,
		sess.RoomInfo.ID,
		activity.Message,
	)

	// create activity in DB
	// TODO: handle db error
	db := GetDB(ctx)
	db.Create(&activity)

	resp = nil
	BroadcastActivity(sess.Room, activity)
	return
}

func joinRoom(ctx context.Context, req interface{}) (resp interface{}, err error) {

	sess := GetSession(ctx)
	if sess == nil {
		err = fmt.Errorf("session not found")
		return
	}
	if sess.Conn == nil {
		err = fmt.Errorf("socket not found in session")
		return
	}

	jsonRequest := GetJSONReq(ctx)
	if jsonRequest == nil {
		err = fmt.Errorf("jsonRequest not found in context")
		return
	}

	// Find the room
	idToJoin := uint(0)
	switch roomIDinJSON := jsonRequest.Get("room_id"); roomIDinJSON.Type() {
	case lzjson.TypeNumber:
		idToJoin = uint(roomIDinJSON.Int())
	case lzjson.TypeString:
		idParsed, _ := strconv.ParseFloat(roomIDinJSON.String(), 64)
		idToJoin = uint(idParsed)
	}

	db := GetDB(ctx)
	srv := GetServer(ctx)

	roomToJoin := models.Room{}
	db.Find(&roomToJoin, idToJoin)
	if roomToJoin.ID != idToJoin {
		log.Printf("%s failed to join room %d",
			sess.HTTPRequest.RemoteAddr,
			roomToJoin.ID,
		)
		err = fmt.Errorf("room (id=%d) not found", idToJoin)
		return
	}

	// unregister client from old room
	if sess.Room != nil {
		sess.Room.Unsubscribe(sess.Conn)
	}

	// attach the client to the room
	if _, ok := srv.rooms[uint64(roomToJoin.ID)]; ok {
		log.Printf("%s joinned room %d",
			sess.HTTPRequest.RemoteAddr,
			roomToJoin.ID,
		)
		sess.RoomInfo = roomToJoin
		sess.Room = srv.rooms[uint64(roomToJoin.ID)]
	} else {
		log.Printf("%s reactivated and joinned room %d",
			sess.HTTPRequest.RemoteAddr,
			roomToJoin.ID,
		)
		sess.RoomInfo = roomToJoin
		sess.Room = NewRoom()
		srv.rooms[uint64(roomToJoin.ID)] = sess.Room
	}

	// register client to new room
	sess.Room.Subscribe(sess.Conn)
	sess.RoomInfo = roomToJoin
	resp = sess.RoomInfo
	return
}
