package pubsub

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/go-restit/lzjson"
	"github.com/tomatorpg/tomatorpg/models"
)

func ping(ctx context.Context, req Request) Response {
	return SuccessResponseTo(req, "pong")
}

func createRoom(ctx context.Context, req Request) Response {
	db := GetDB(ctx)
	// TODO: read request payload for room data
	newRoom := models.Room{}
	newRoom.ID = 0 // ensure not injecting ID
	db.Create(&newRoom)
	log.Printf("rooms.create: id=%d", newRoom.ID)
	return SuccessResponseTo(req, nil)
}

func listRooms(ctx context.Context, req Request) Response {
	db := GetDB(ctx)
	var rooms []models.Room
	db.Order("created_at desc").Find(&rooms)
	log.Printf("rooms.list length=%d", len(rooms))
	return SuccessResponseTo(req, rooms)
}

func replayRoom(ctx context.Context, req Request) Response {

	// TODO: this is temp API, should do with CURD
	//       should rewrite Replay as normal crud listing
	//       to be independent from websocket
	sess := GetSession(ctx)
	if sess == nil {
		return ErrorResponseTo(
			req,
			fmt.Errorf("session not found"),
		)
	}
	if sess.Conn == nil {
		return ErrorResponseTo(
			req,
			fmt.Errorf("socket not found"),
		)
	}
	if sess.Room == nil {
		return ErrorResponseTo(
			req,
			fmt.Errorf("the session is not currently in a room"),
		)
	}

	log.Printf("rooms.replay: id=%d", sess.Room.Info.ID)
	resp := SuccessResponseTo(req, sess.Room.Info.ID)

	sess.Room.Replay(sess.Conn)
	return resp
}

func createRoomActivity(ctx context.Context, req Request) Response {
	// TODO: rewrite to pure crud
	sess := GetSession(ctx)
	if sess == nil {
		return ErrorResponseTo(
			req,
			fmt.Errorf("session not found"),
		)
	}
	if sess.Conn == nil {
		return ErrorResponseTo(
			req,
			fmt.Errorf("socket not found in session"),
		)
	}
	if sess.Room == nil {
		return ErrorResponseTo(
			req,
			fmt.Errorf("room not found in session"),
		)
	}

	// get raw json request from context
	// TODO: remove the need to do json decode here
	// TODO: validate payload format
	var activity models.RoomActivity
	jsonRequest := GetJsonReq(ctx)
	if jsonRequest == nil {
		return ErrorResponseTo(
			req,
			fmt.Errorf("jsonRequest not found in context"),
		)
	}
	jsonRequest.Get("payload").Unmarshal(&activity)
	activity.UserID = sess.User.ID // enforce user session
	if activity.Action == "" {
		activity.Action = "message"
	}
	log.Printf("roomActivity: user-%d %s in room-%d: %s",
		activity.UserID,
		activity.Action,
		sess.Room.Info.ID,
		activity.Message,
	)

	resp := SuccessResponseTo(req, nil)
	sess.Room.Broadcast(activity)
	return resp
}

func joinRoom(ctx context.Context, req Request) Response {

	sess := GetSession(ctx)
	if sess == nil {
		return ErrorResponseTo(
			req,
			fmt.Errorf("session not found"),
		)
	}
	if sess.Conn == nil {
		return ErrorResponseTo(
			req,
			fmt.Errorf("socket not found in session"),
		)
	}

	jsonRequest := GetJsonReq(ctx)
	if jsonRequest == nil {
		return ErrorResponseTo(
			req,
			fmt.Errorf("jsonRequest not found in context"),
		)
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
	if roomToJoin.ID == idToJoin {

		// unregister client from old room
		if sess.Room != nil {
			sess.Room.Unregister(sess.Conn)
		}

		// attach the client to the room
		if _, ok := srv.rooms[uint64(roomToJoin.ID)]; ok {
			log.Printf("%s joinned room %d",
				sess.HttpRequest.RemoteAddr,
				roomToJoin.ID,
			)
			sess.Room = srv.rooms[uint64(roomToJoin.ID)]
		} else {
			log.Printf("%s reactivated and joinned room %d",
				sess.HttpRequest.RemoteAddr,
				roomToJoin.ID,
			)
			sess.Room = NewRoom()
			sess.Room.Info = roomToJoin
			srv.rooms[uint64(roomToJoin.ID)] = sess.Room
		}

		// register client to new room
		sess.Room.Register(sess.Conn)
		return SuccessResponseTo(req, nil)
	} else {
		log.Printf("%s failed to join room %d",
			sess.HttpRequest.RemoteAddr,
			roomToJoin.ID,
		)
		return ErrorResponseTo(
			req,
			fmt.Errorf("room (id=%d) not found", idToJoin),
		)
	}
}
