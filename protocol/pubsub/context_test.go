package pubsub_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/go-restit/lzjson"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/tomatorpg/tomatorpg/protocol/pubsub"
)

func TestContext(t *testing.T) {

	// declare variables in context
	getContext := func() (ctx context.Context, err error) {
		httpReq, _ := http.NewRequest("GET", "http://foobar.com/hello/world", nil)

		sess := &pubsub.Session{
			HTTPRequest: httpReq,
		}

		jsonReq := lzjson.NewNode()
		json.Unmarshal([]byte(`{"hello": "world"}`), &jsonReq)

		db, err := gorm.Open("sqlite3", ":memory:")
		if err != nil {
			return
		}

		srv := pubsub.NewServer(
			db,
			make(pubsub.WebsocketChanColl),
			pubsub.RPCs(),
			"abcde",
		)

		ctx = context.Background()
		ctx = pubsub.WithSession(ctx, sess)
		ctx = pubsub.WithJSONReq(ctx, jsonReq)
		ctx = pubsub.WithDB(ctx, db)
		ctx = pubsub.WithServer(ctx, srv)
		return
	}

	ctx, err := getContext()
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
		return
	}

	if sess := pubsub.GetSession(ctx); sess == nil {
		t.Errorf("expected session, got nil")
	} else if want, have := "http://foobar.com/hello/world", sess.HTTPRequest.URL.String(); want != have {
		t.Errorf("expected %s, got %s", want, have)
	}

	if jsonReq := pubsub.GetJSONReq(ctx); jsonReq == nil {
		t.Errorf("expected jsonReq, got nil")
	} else if want, have := "world", jsonReq.Get("hello").String(); want != have {
		t.Errorf("expected %s, got %s", want, have)
	}

	if db := pubsub.GetDB(ctx); db == nil {
		t.Errorf("expected *gorm.DB, got nil")
	} else {
		db.Close()
	}

	if srv := pubsub.GetServer(ctx); srv == nil {
		t.Errorf("expected pubsub.Server, got nil")
	}

}
