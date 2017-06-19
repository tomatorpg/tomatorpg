package pubsub

import (
	"fmt"
	"log"
	"testing"
)

type dummyChanColl map[uint]Channel

func (coll dummyChanColl) LoadOrOpen(id uint) Channel {
	if _, ok := coll[id]; !ok {
		coll[id] = newDummyChannel()
	}
	return coll[id]
}

func (coll dummyChanColl) Close(id uint) {

}

type dummyChannel struct {
	broadcast chan interface{}
	conns     map[MessageWriteCloser]bool
}

func newDummyChannel() Channel {
	ch := &dummyChannel{
		broadcast: make(chan interface{}),
		conns:     make(map[MessageWriteCloser]bool),
	}
	go ch.run()
	return ch
}

func (ch *dummyChannel) Subscribe(conn MessageWriteCloser) {
	ch.conns[conn] = true
}

func (ch *dummyChannel) Unsubscribe(conn MessageWriteCloser) {
	delete(ch.conns, conn)
}

func (ch *dummyChannel) BroadcastJSON(v interface{}) {
	for conn := range ch.conns {
		conn.WriteJSON(v)
	}
}

func (ch *dummyChannel) run() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-ch.broadcast

		// Send it out to every client that is currently connected
		for client := range ch.conns {
			err := client.WriteJSON(msg)
			if err != nil {
				client.Close()
				ch.Unsubscribe(client)
				log.Printf("error: %v", err)
			}
		}
	}
}

type dummyWriter struct {
	lastMsg      interface{}
	lastMsgType  int
	lastMsgBytes []byte
}

func (w *dummyWriter) WriteMessage(messageType int, p []byte) error {
	w.lastMsg = nil
	w.lastMsgType = messageType
	w.lastMsgBytes = make([]byte, len(p))
	copy(w.lastMsgBytes, p)
	return nil
}

func (w *dummyWriter) WriteJSON(v interface{}) error {
	w.lastMsg = v
	w.lastMsgType = -1
	w.lastMsgBytes = make([]byte, 0)
	return nil
}

func (w *dummyWriter) Close() error {
	return nil
}

type errMsgWriter int

func (w errMsgWriter) WriteMessage(messageType int, p []byte) error {
	return fmt.Errorf("dummy error, %#v, %#v", messageType, p)
}

func (w errMsgWriter) WriteJSON(v interface{}) error {
	return fmt.Errorf("dummy error, %#v", v)
}

func (w errMsgWriter) Close() error {
	return nil
}

func TestMessageTo(t *testing.T) {
	ch := &dummyChannel{
		conns: make(map[MessageWriteCloser]bool),
	}
	w1 := errMsgWriter(0)
	ch.Subscribe(w1)
	err := messageTo(ch, w1, "hello message")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if want, have := `dummy error, "hello message"`, err.Error(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}
