package pubsub

import (
	"context"
	"testing"
	"time"

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
