package pubsub_test

import (
	"context"
	"fmt"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-restit/lzjson"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/tomatorpg/tomatorpg/protocol/pubsub"
)

func TestServer_ServeHTTP(t *testing.T) {

	// serial generator
	serial := func(init int) func() int {
		out := make(chan int)
		go func() {
			for i := init; true; i++ {
				out <- i
			}
		}()
		return func() int {
			return <-out
		}
	}(1)

	// dummy database
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		return
	}
	defer db.Close()

	mustConnect := func(url string) (conn *websocket.Conn) {
		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			t.Fatalf("cannot make websocket connection: %v", err)
		}
		return
	}

	rtr := pubsub.NewRouter()
	rtr.Add("dummy", "world", "hello", func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := pubsub.Request{}
		pubsub.GetJSONReq(ctx).Unmarshal(&req)
		response = fmt.Sprintf("%s %s", req.Method, req.Entity)
		return
	})
	rtr.Add("dummy", "bar", "foo", func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := pubsub.Request{}
		pubsub.GetJSONReq(ctx).Unmarshal(&req)
		err = fmt.Errorf("%s %s", req.Method, req.Entity)
		return
	})

	// server with dummy components for test
	srv := httptest.NewServer(pubsub.NewServer(
		db,
		make(pubsub.WebsocketChanColl), // TODO: use dummy chan coll
		rtr,
		"abcde",
	))
	defer srv.Close()

	// make 2 separated ws connection to the dummy room server
	u, _ := url.Parse(srv.URL)
	u.Scheme = "ws"
	conn := mustConnect(u.String())
	t.Logf("connection success")

	toTests := []struct {
		desc string
		req  pubsub.Request
		resp pubsub.Response
	}{
		{
			desc: "hello world 1",
			req: pubsub.Request{
				ID:     fmt.Sprintf("%d", serial()),
				Group:  "dummy",
				Entity: "world",
				Method: "hello",
			},
			resp: pubsub.Response{
				ID:     "1",
				Entity: "world",
				Method: "hello",
				Data:   "hello world",
			},
		},
		{
			desc: "hello world 1",
			req: pubsub.Request{
				ID:     fmt.Sprintf("%d", serial()),
				Group:  "dummy",
				Entity: "world",
				Method: "hello",
			},
			resp: pubsub.Response{
				ID:     "2",
				Entity: "world",
				Method: "hello",
				Data:   "hello world",
			},
		},
		{
			desc: "hello world 1",
			req: pubsub.Request{
				ID:     fmt.Sprintf("%d", serial()),
				Group:  "dummy",
				Entity: "bar",
				Method: "foo",
			},
			resp: pubsub.Response{
				ID:     "3",
				Entity: "bar",
				Method: "foo",
				Err:    "foo bar",
			},
		},
	}

	for i, toTest := range toTests {
		resp := make(map[string]interface{})
		rawResp := lzjson.NewNode()
		go func() {
			conn.WriteJSON(toTest.req)
		}()
		conn.ReadJSON(rawResp)
		rawResp.Unmarshal(&resp)
		t.Logf("-- run test %d --", i)
		if want, have := toTest.resp.ID, resp["id"]; want != have {
			t.Errorf("expected %#v, got %#v", want, have)
		}
		if want, have := toTest.resp.Entity, resp["entity"]; want != have {
			t.Errorf("expected %#v, got %#v", want, have)
		}
		if want, have := toTest.resp.Method, resp["method"]; want != have {
			t.Errorf("expected %#v, got %#v", want, have)
		}
		if want, have := toTest.resp.Data, resp["data"]; want != have {
			t.Errorf("expected %#v, got %#v", want, have)
		}
		if toTest.resp.Err != "" {
			if want, have := toTest.resp.Err, resp["error"]; want != have {
				t.Errorf("expected %#v, got %#v", want, have)
			}
		} else if _, ok := resp["error"]; ok {
			t.Errorf("unexpected error: %s", resp["error"])
			t.Logf("raw json: %s", rawResp.Get("error").Raw())
		}
	}

}
