package pubsub

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/go-restit/lzjson"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/tomatorpg/tomatorpg/models"
)

func TestProcedure_ping(t *testing.T) {
	resp, err := ping(nil, nil)
	if want, have := "pong", resp; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if err != nil {
		t.Errorf("expected nil, got %#v", err)
	}
}

func TestProcedure_whoami(t *testing.T) {
	ctx := WithSession(context.Background(), &Session{
		User: models.User{
			Model: gorm.Model{
				ID: uint(12345),
			},
			Name: "Dummy User",
		},
	})
	rawResp, err := whoami(ctx, nil)
	if err != nil {
		t.Errorf("expected nil, got %#v", err)
	}
	if resp, ok := rawResp.(whoamiResp); !ok {
		t.Errorf("expected whoamiResp, got %#v", resp)
		return
	} else if want, have := uint(12345), resp.ID; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	} else if want, have := "Dummy User", resp.Name; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}

func TestProcedure_whoami_err(t *testing.T) {
	resp, err := whoami(context.Background(), nil)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if resp != nil {
		t.Errorf("expected nil, got %#v", resp)
	}
	if want, have := "session not found", err.Error(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}

func TestProcedure_createRoom(t *testing.T) {
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	defer db.Close()
	db.AutoMigrate(&models.Room{})

	ctx := WithDB(context.Background(), db)
	resp, err := createRoom(ctx, nil)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if resp == nil {
		t.Errorf("expected resp, got nil")
	}
	if room, ok := resp.(models.Room); !ok {
		t.Errorf("expected models.Room, got %#v", resp)
	} else if want, have := uint(1), room.ID; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}

func TestProcedure_listRoom(t *testing.T) {
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	defer db.Close()
	db.AutoMigrate(&models.Room{})

	db.Create(models.Room{
		ID:        1234,
		Name:      "Hello Room 1",
		CreatedAt: time.Now(),
	})
	db.Create(models.Room{
		ID:        1235,
		Name:      "Hello Room 2",
		CreatedAt: time.Now().Add(1 * time.Second),
	})

	ctx := WithDB(context.Background(), db)
	resp, err := listRooms(ctx, nil)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if resp == nil {
		t.Errorf("expected resp, got nil")
	}
	rooms, ok := resp.([]models.Room)
	if !ok {
		t.Errorf("expected []models.Room, got %#v", resp)
	} else if want, have := 2, len(rooms); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
		return
	}

	if want, have := uint(1235), rooms[0].ID; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "Hello Room 2", rooms[0].Name; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := uint(1234), rooms[1].ID; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "Hello Room 1", rooms[1].Name; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}

func TestProcedure_joinRoom(t *testing.T) {
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	defer db.Close()
	db.AutoMigrate(&models.Room{})
	db.Create(models.Room{
		ID:        1234,
		Name:      "Hello Room 1",
		CreatedAt: time.Now(),
	})
	db.Create(models.Room{
		ID:        1235,
		Name:      "Hello Room 2",
		CreatedAt: time.Now().Add(1 * time.Second),
	})

	// server with dummy components for test
	coll := make(intlDummyChanColl)
	srv := NewServer(
		db,
		coll,
		NewRouter(),
		"abcde",
	)
	conn := &dummyWriter{}
	sess := &Session{
		Conn: conn,
	}

	reqJSON := lzjson.NewNode()
	json.Unmarshal([]byte(`{"room_id": 1234}`), reqJSON)

	ctx := context.Background()
	ctx = WithDB(ctx, db)
	ctx = WithServer(ctx, srv)
	ctx = WithSession(ctx, sess)
	ctx = WithJSONReq(ctx, reqJSON)
	joinRoom(ctx, nil)

	ch, ok := coll[1234]
	if !ok {
		t.Errorf("expect coll[1234] to exist")
	}

	chReal := ch.(*intlDummyChannel)
	if _, ok := chReal.conns[sess.Conn]; !ok {
		t.Errorf("connection not found in channel")
	}
}

func TestProcedure_joinRoom_withprevchan(t *testing.T) {
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	defer db.Close()
	db.AutoMigrate(&models.Room{})
	db.Create(models.Room{
		ID:        1234,
		Name:      "Hello Room 1",
		CreatedAt: time.Now(),
	})
	db.Create(models.Room{
		ID:        1235,
		Name:      "Hello Room 2",
		CreatedAt: time.Now().Add(1 * time.Second),
	})

	// server with dummy components for test
	coll := make(intlDummyChanColl)
	srv := NewServer(
		db,
		coll,
		NewRouter(),
		"abcde",
	)
	conn := &dummyWriter{}
	sess := &Session{
		Conn: conn,
	}

	reqJSON := lzjson.NewNode()
	json.Unmarshal([]byte(`{"room_id": 1234}`), reqJSON)

	ctx := context.Background()
	ctx = WithDB(ctx, db)
	ctx = WithServer(ctx, srv)
	ctx = WithSession(ctx, sess)
	ctx = WithJSONReq(ctx, reqJSON)
	joinRoom(ctx, nil)

	// record the room result for 1234
	sess.RoomChan = coll[1234]
	chReal1234 := sess.RoomChan.(*intlDummyChannel)

	// join room 1235 and get new status
	json.Unmarshal([]byte(`{"room_id": 1235}`), reqJSON)
	ctx = WithJSONReq(ctx, reqJSON)
	resp, err := joinRoom(ctx, nil)
	if err != nil {
		t.Errorf("unexpected error: %#v", err.Error())
	}
	ch, ok := coll[1235]
	if !ok {
		t.Errorf("expect coll[1235] to exist")
	}

	if want, have := uint(1235), sess.RoomInfo.ID; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "Hello Room 2", sess.RoomInfo.Name; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	// conn should be in channel 1235
	chReal1235 := ch.(*intlDummyChannel)
	if _, ok := chReal1235.conns[sess.Conn]; !ok {
		t.Errorf("connection not found in channel 1235, unexpected")
	}

	// conn should not be in channel 1234
	if _, ok := chReal1234.conns[sess.Conn]; ok {
		t.Errorf("connection found in channel 1234, unexpected")
	}

	// resp test
	realResp, ok := resp.(models.Room)
	if !ok {
		t.Errorf("expected models.Room, got %T", resp)
		return
	}
	if want, have := uint(1235), realResp.ID; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "Hello Room 2", realResp.Name; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

}

func TestProcedure_joinRoom_strid(t *testing.T) {
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	defer db.Close()
	db.AutoMigrate(&models.Room{})
	db.Create(models.Room{
		ID:        1234,
		Name:      "Hello Room 1",
		CreatedAt: time.Now(),
	})
	db.Create(models.Room{
		ID:        1235,
		Name:      "Hello Room 2",
		CreatedAt: time.Now().Add(1 * time.Second),
	})

	// server with dummy components for test
	coll := make(intlDummyChanColl)
	srv := NewServer(
		db,
		coll,
		NewRouter(),
		"abcde",
	)
	conn := &dummyWriter{}
	sess := &Session{
		Conn: conn,
	}

	reqJSON := lzjson.NewNode()
	json.Unmarshal([]byte(`{"room_id": "1234"}`), reqJSON)

	ctx := context.Background()
	ctx = WithDB(ctx, db)
	ctx = WithServer(ctx, srv)
	ctx = WithSession(ctx, sess)
	ctx = WithJSONReq(ctx, reqJSON)
	joinRoom(ctx, nil)

	ch, ok := coll[1234]
	if !ok {
		t.Errorf("expect coll[1234] to exist")
	}

	chReal := ch.(*intlDummyChannel)
	if _, ok := chReal.conns[sess.Conn]; !ok {
		t.Errorf("connection not found in channel")
	}
}

func TestProcedure_joinRoom_nosession(t *testing.T) {
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	defer db.Close()
	db.AutoMigrate(&models.Room{})
	db.Create(models.Room{
		ID:        1234,
		Name:      "Hello Room 1",
		CreatedAt: time.Now(),
	})
	db.Create(models.Room{
		ID:        1235,
		Name:      "Hello Room 2",
		CreatedAt: time.Now().Add(1 * time.Second),
	})

	// server with dummy components for test
	coll := make(intlDummyChanColl)
	srv := NewServer(
		db,
		coll,
		NewRouter(),
		"abcde",
	)

	reqJSON := lzjson.NewNode()
	json.Unmarshal([]byte(`{"room_id": 1234}`), reqJSON)

	ctx := context.Background()
	ctx = WithDB(ctx, db)
	ctx = WithServer(ctx, srv)
	ctx = WithJSONReq(ctx, reqJSON)

	_, err = joinRoom(ctx, nil)
	if err == nil {
		t.Errorf("expected to have error, got nil")
		return
	}
	if want, have := "session not found", err.Error(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}

func TestProcedure_joinRoom_nosocket(t *testing.T) {
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	defer db.Close()
	db.AutoMigrate(&models.Room{})
	db.Create(models.Room{
		ID:        1234,
		Name:      "Hello Room 1",
		CreatedAt: time.Now(),
	})
	db.Create(models.Room{
		ID:        1235,
		Name:      "Hello Room 2",
		CreatedAt: time.Now().Add(1 * time.Second),
	})

	// server with dummy components for test
	coll := make(intlDummyChanColl)
	srv := NewServer(
		db,
		coll,
		NewRouter(),
		"abcde",
	)
	sess := &Session{
		HTTPRequest: &http.Request{},
	}

	reqJSON := lzjson.NewNode()
	json.Unmarshal([]byte(`{"room_id": 1234}`), reqJSON)

	ctx := context.Background()
	ctx = WithDB(ctx, db)
	ctx = WithServer(ctx, srv)
	ctx = WithSession(ctx, sess)
	ctx = WithJSONReq(ctx, reqJSON)

	_, err = joinRoom(ctx, nil)
	if err == nil {
		t.Errorf("expected to have error, got nil")
		return
	}
	if want, have := "socket not found in session", err.Error(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}

func TestProcedure_joinRoom_norequest(t *testing.T) {
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	defer db.Close()
	db.AutoMigrate(&models.Room{})
	db.Create(models.Room{
		ID:        1234,
		Name:      "Hello Room 1",
		CreatedAt: time.Now(),
	})
	db.Create(models.Room{
		ID:        1235,
		Name:      "Hello Room 2",
		CreatedAt: time.Now().Add(1 * time.Second),
	})

	// server with dummy components for test
	coll := make(intlDummyChanColl)
	srv := NewServer(
		db,
		coll,
		NewRouter(),
		"abcde",
	)
	conn := &dummyWriter{}
	sess := &Session{
		Conn:        conn,
		HTTPRequest: &http.Request{},
	}

	ctx := context.Background()
	ctx = WithDB(ctx, db)
	ctx = WithServer(ctx, srv)
	ctx = WithSession(ctx, sess)

	_, err = joinRoom(ctx, nil)
	if err == nil {
		t.Errorf("expected to have error, got nil")
		return
	}
	if want, have := "jsonRequest not found in context", err.Error(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}
func TestProcedure_joinRoom_notfound(t *testing.T) {
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	defer db.Close()
	db.AutoMigrate(&models.Room{})
	db.Create(models.Room{
		ID:        1234,
		Name:      "Hello Room 1",
		CreatedAt: time.Now(),
	})
	db.Create(models.Room{
		ID:        1235,
		Name:      "Hello Room 2",
		CreatedAt: time.Now().Add(1 * time.Second),
	})

	// server with dummy components for test
	coll := make(intlDummyChanColl)
	srv := NewServer(
		db,
		coll,
		NewRouter(),
		"abcde",
	)
	conn := &dummyWriter{}
	sess := &Session{
		Conn:        conn,
		HTTPRequest: &http.Request{},
	}

	reqJSON := lzjson.NewNode()
	json.Unmarshal([]byte(`{"room_id": 1230}`), reqJSON)

	ctx := context.Background()
	ctx = WithDB(ctx, db)
	ctx = WithServer(ctx, srv)
	ctx = WithSession(ctx, sess)
	ctx = WithJSONReq(ctx, reqJSON)

	_, err = joinRoom(ctx, nil)
	if err == nil {
		t.Errorf("expected to have error, got nil")
		return
	}
	if want, have := "room (id=1230) not found", err.Error(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}
