package pubsub

import (
	"context"
	"fmt"
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
	logger := GetLogContext(ctx)
	// TODO: read request payload for room data
	newRoom := models.Room{}
	newRoom.ID = 0 // ensure not injecting ID
	db.Create(&newRoom)
	logger.Log(
		"action", "rooms.create",
		"room.id", newRoom.ID,
	)
	resp = newRoom
	return
}

func listRooms(ctx context.Context, req interface{}) (resp interface{}, err error) {
	db := GetDB(ctx)
	var rooms []models.Room
	db.Order("created_at desc").Find(&rooms)
	logger := GetLogContext(ctx)
	logger.Log(
		"action", "rooms.list",
		"len(room)", len(rooms),
	)
	resp = rooms
	return
}

func replayRoom(ctx context.Context, req interface{}) (resp interface{}, err error) {

	db := GetDB(ctx)
	logger := GetLogContext(ctx)

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
	if sess.RoomChan == nil {
		err = fmt.Errorf("the session is not currently in a room")
		return
	}

	logger.Log(
		"action", "rooms.replay",
		"room.id", sess.RoomInfo.ID,
	)
	resp = sess.RoomInfo.ID

	// replay history (TODO: rewrite as pure CRUD)
	historyCopy := make([]models.RoomActivity, 0, 100)
	db.Find(&historyCopy, "room_id = ?", sess.RoomInfo.ID)
	if len(historyCopy) > 0 {
		logger.Log(
			"message", "start replaying activities to client",
			"room.id", sess.RoomInfo.ID,
		)
		for _, activity := range historyCopy {
			err := sess.Conn.WriteJSON(Broadcast{
				Version: "0.2",
				Entity:  "roomActivities",
				Type:    "broadcast",
				Data:    activity,
			})
			if err != nil {
				sess.Conn.Close()
				sess.RoomChan.Unsubscribe(sess.Conn)
				logger.Log(
					"message", "error writing JSON to socket",
					"error", err.Error(),
				)
				break
			}
		}
	}
	return
}

func createRoomActivity(ctx context.Context, req interface{}) (resp interface{}, err error) {

	logger := GetLogContext(ctx)

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
	if sess.RoomChan == nil {
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
	logger.Log(
		"message", "roomActivity",
		"user.id", activity.UserID,
		"room.id", sess.RoomInfo.ID,
		"activity.action", activity.Action,
		"activity.message", activity.Message,
	)

	// create activity in DB
	// TODO: handle db error
	db := GetDB(ctx)
	db.Create(&activity)

	resp = nil
	BroadcastActivity(sess.RoomChan, activity)
	return
}

func joinRoom(ctx context.Context, req interface{}) (resp interface{}, err error) {

	logger := GetLogContext(ctx)

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
		logger.Log(
			"message", "failed to join room, room not found",
			"room.id", roomToJoin.ID,
		)
		err = fmt.Errorf("room (id=%d) not found", idToJoin)
		return
	}

	// unregister client from old room
	if sess.RoomChan != nil && sess.Conn != nil {
		sess.RoomChan.Unsubscribe(sess.Conn)
	}

	// attach the client to the room
	sess.RoomInfo = roomToJoin
	sess.RoomChan = srv.chans.LoadOrOpen(roomToJoin.ID)

	// register client to new room
	sess.RoomChan.Subscribe(sess.Conn)
	sess.RoomInfo = roomToJoin
	resp = sess.RoomInfo
	return
}
