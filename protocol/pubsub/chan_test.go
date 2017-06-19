package pubsub_test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tomatorpg/tomatorpg/models"
	"github.com/tomatorpg/tomatorpg/protocol/pubsub"
)

func TestWebsocketChan_Broadcast(t *testing.T) {

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	// dummy chan to test
	wsChan := pubsub.NewWebsocketChan()

	// serial generator
	serial := func() func() int {
		out := make(chan int)
		go func() {
			for i := 0; true; i++ {
				out <- i
			}
		}()
		return func() int {
			return <-out
		}
	}()

	testRoomServer := func(w http.ResponseWriter, req *http.Request) {

		reqID := serial()

		// upgrade to websocket connection
		conn, err := upgrader.Upgrade(w, req, nil)
		if err != nil {
			t.Logf("[req: %d] cannot upgrade: %v", reqID, err)
			errResp := pubsub.Response{
				Status: "error",
				Err:    err,
			}
			if b, err := json.Marshal(errResp); err != nil {
				http.Error(w, fmt.Sprintf("%s", b), http.StatusInternalServerError)
			} else {
				http.Error(w, fmt.Sprintf("[req: %d] cannot encode response: %v", reqID, err), http.StatusInternalServerError)
			}
		}

		go func() {
			// register connection to chan
			wsChan.Subscribe(conn)
			defer conn.Close()
			defer wsChan.Unsubscribe(conn)

			// dummy loop for connection handle
			for {
				v := make(map[string]interface{})
				conn.ReadJSON(&v)
				t.Logf("[req: %d] server received: %#v", reqID, v)
			}
		}()
	}

	mustConnect := func(url string) (conn *websocket.Conn) {
		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			log.Fatalf("cannot make websocket connection: %v", err)
		}
		return
	}

	var err error

	srv := httptest.NewServer(http.HandlerFunc(testRoomServer))
	defer srv.Close()

	// make 2 separated ws connection to the dummy room server
	u, _ := url.Parse(srv.URL)
	u.Scheme = "ws"
	conn1 := mustConnect(u.String())
	conn2 := mustConnect(u.String())
	t.Logf("connection success")

	go func() {
		pubsub.BroadcastActivity(wsChan, models.RoomActivity{
			Action:  "say",
			Message: "hello",
		})
		t.Logf("broadcast sent")
	}()

	bc1 := pubsub.Broadcast{}
	err = conn1.ReadJSON(&bc1)
	if err != nil {
		t.Fatalf("cannot read message: %v", err)
	}
	if want, have := "say", bc1.Data.Action; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "hello", bc1.Data.Message; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	bc2 := pubsub.Broadcast{}
	err = conn2.ReadJSON(&bc2)
	if err != nil {
		t.Fatalf("cannot read message: %v", err)
	}
	if want, have := "say", bc2.Data.Action; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "hello", bc2.Data.Message; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

}

func TestWebsocketChan_Unsubscribe(t *testing.T) {

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	// dummy chan to test
	wsChan := pubsub.NewWebsocketChan()

	// serial generator
	serial := func() func() int {
		out := make(chan int)
		go func() {
			for i := 0; true; i++ {
				out <- i
			}
		}()
		return func() int {
			return <-out
		}
	}()

	testRoomServer := func(w http.ResponseWriter, req *http.Request) {

		reqID := serial()

		// upgrade to websocket connection
		conn, err := upgrader.Upgrade(w, req, nil)
		if err != nil {
			t.Logf("[req: %d] cannot upgrade: %v", reqID, err)
			errResp := pubsub.Response{
				Status: "error",
				Err:    err,
			}
			if b, err := json.Marshal(errResp); err != nil {
				http.Error(w, fmt.Sprintf("%s", b), http.StatusInternalServerError)
			} else {
				http.Error(w, fmt.Sprintf("[req: %d] cannot encode response: %v", reqID, err), http.StatusInternalServerError)
			}
		}

		go func() {
			// register connection to room
			wsChan.Subscribe(conn)
			defer conn.Close()
			defer wsChan.Unsubscribe(conn)

			// dummy loop for connection handle
			for {
				v := make(map[string]interface{})
				conn.ReadJSON(&v)
				t.Logf("[req: %d] server received: %#v", reqID, v)
			}
		}()
	}

	mustConnect := func(url string) (conn *websocket.Conn) {
		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			log.Fatalf("cannot make websocket connection: %v", err)
		}
		return
	}

	var err error

	srv := httptest.NewServer(http.HandlerFunc(testRoomServer))
	defer srv.Close()

	// make 2 separated ws connection to the dummy room server
	u, _ := url.Parse(srv.URL)
	u.Scheme = "ws"
	conn1 := mustConnect(u.String())
	conn2 := mustConnect(u.String())
	t.Logf("connection success")

	go func() {
		pubsub.BroadcastActivity(wsChan, models.RoomActivity{
			Action:  "say",
			Message: "hello 1",
		})
		t.Logf("broadcast sent: say hello 1")
	}()

	bc1 := pubsub.Broadcast{}
	err = conn1.ReadJSON(&bc1)
	if err != nil {
		t.Fatalf("cannot read message: %v", err)
	}
	if want, have := "say", bc1.Data.Action; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "hello 1", bc1.Data.Message; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	bc2 := pubsub.Broadcast{}
	err = conn2.ReadJSON(&bc2)
	if err != nil {
		t.Fatalf("cannot read message: %v", err)
	}
	if want, have := "say", bc2.Data.Action; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "hello 1", bc2.Data.Message; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	wsChan.Unsubscribe(conn2)

	go func() {
		pubsub.BroadcastActivity(wsChan, models.RoomActivity{
			Action:  "say",
			Message: "hello 2",
		})
		t.Logf("broadcast sent: say hello 2")
	}()

	err = conn1.ReadJSON(&bc1)
	if err != nil {
		t.Fatalf("cannot read message: %v", err)
	}
	if want, have := "say", bc1.Data.Action; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "hello 2", bc1.Data.Message; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	// TODO: test if bc2 received anything before timeout

	select {
	case <-time.After(20 * time.Millisecond):
		t.Logf("conn2 timeout as expected")
	case recieved := <-func() <-chan interface{} {
		out := make(chan interface{})
		go func() {
			conn2.ReadJSON(&bc2)
			out <- bc2
		}()
		return out
	}():
		t.Logf("conn2 received unexpected message: %#v", recieved)
	}
}

type errMsgWriter int

func (w errMsgWriter) WriteMessage(messageType int, p []byte) error {
	return nil
}

func (w errMsgWriter) WriteJSON(v interface{}) error {
	return nil
}

func (w errMsgWriter) Close() error {
	return nil
}

func TestWebsocketChan_Subscribe_err(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("expected error, got nil")
			return
		}
		if prefix, have := "*WebsocketChannel only allow registering *websocket.Conn", r.(string); !strings.HasPrefix(have, prefix) {
			t.Errorf("expected error message to have prefix %#v, got %#v",
				prefix, have)
		}
	}()
	wsChan := pubsub.NewWebsocketChan()
	wsChan.Subscribe(errMsgWriter(0))
}

func TestWebsocketChan_Unsubscribe_err(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("expected error, got nil")
			return
		}
		if prefix, have := "*WebsocketChannel only allow unregistering *websocket.Conn", r.(string); !strings.HasPrefix(have, prefix) {
			t.Errorf("expected error message to have prefix %#v, got %#v",
				prefix, have)
		}
	}()
	wsChan := pubsub.NewWebsocketChan()
	wsChan.Unsubscribe(errMsgWriter(0))
}
