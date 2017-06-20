package pubsub_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	kitlog "github.com/go-kit/kit/log"
	"github.com/tomatorpg/tomatorpg/protocol/pubsub"
)

func TestChain(t *testing.T) {
	i, pm1, pm2 := 0, 0, 0
	m1 := func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			i++
			pm1 = i
			inner.ServeHTTP(w, r)
		})
	}
	m2 := func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			i++
			pm2 = i
			inner.ServeHTTP(w, r)
		})
	}
	srv := pubsub.Chain(m1, m2)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// do nothing
	}))

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "http://foobar.com/hello/world", nil)
	srv.ServeHTTP(w, r)
	if want, have := 1, pm1; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := 2, pm2; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	srv.ServeHTTP(w, r)
	if want, have := 3, pm1; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := 4, pm2; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}

func TestApplyRequestID(t *testing.T) {
	catchID := false
	srv := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if reqID := r.Header.Get("X-Request-ID"); reqID != "" {
			catchID = true
			t.Logf("got request id: %s", reqID)
			fmt.Fprintf(w, "success")
			return
		}
		t.Errorf("got no request id")
		fmt.Fprintf(w, "failed")
	}))
	srv = pubsub.ApplyRequestID(srv)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "http://foobar.com/hello/world", nil)
	srv.ServeHTTP(w, r)
}

func TestApplyContextLog(t *testing.T) {
	srv := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := pubsub.GetLogContext(r.Context())
		if logger != nil {
			logger.Log("hello", "world")
			fmt.Fprintf(w, "success")
			return
		}
		t.Logf("got no log context")
		fmt.Fprintf(w, "failed")
	}))
	buf := bytes.NewBuffer(make([]byte, 256))
	logger := kitlog.NewLogfmtLogger(buf)
	srv = pubsub.ApplyContextLog(logger)(srv)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "http://foobar.com/hello/world", nil)
	r.RemoteAddr = "http://somewhere.com:1234"
	r.Header.Add("X-Request-ID", "helloid")
	srv.ServeHTTP(w, r)

	if want, have := `request_id=helloid method=GET path=/hello/world protocol=http remote_addr=http://somewhere.com:1234`+"\nrequest_id=helloid hello=world\n", strings.Trim(string(buf.Bytes()), "\x00"); want != have {
		t.Errorf("\nexpected %#v\n     got %#v", want, have)
	}
}
