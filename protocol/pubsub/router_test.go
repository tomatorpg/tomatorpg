package pubsub_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/tomatorpg/tomatorpg/protocol/pubsub"
)

func TestRouter(t *testing.T) {
	r := pubsub.NewRouter()
	r.Add("test", "world", "hello", func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		resp = "hello world"
		return
	})
	r.Add("test", "foo", "bar", func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		err = fmt.Errorf("foo bar error")
		return
	})
	r.NotFound(func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		err = fmt.Errorf("foobar invalid request")
		return
	})

	resp, err := r.ServeRequest(context.Background(), pubsub.Request{
		ID:     "hello-id",
		Group:  "test",
		Entity: "world",
		Method: "hello",
	})
	if want, have := "hello world", resp; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	resp, err = r.ServeRequest(context.Background(), pubsub.Request{
		ID:     "hello-id",
		Group:  "test",
		Entity: "foo",
		Method: "bar",
	})
	if resp != nil {
		t.Errorf("expected nil, got %#v", resp)
	}
	if err == nil {
		t.Errorf("expected error, got nil")
	} else if want, have := "foo bar error", err.Error(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	resp, err = r.ServeRequest(context.Background(), pubsub.Request{
		ID:     "hello-id",
		Group:  "test",
		Entity: "ass",
		Method: "kick",
	})
	if resp != nil {
		t.Errorf("expected nil, got %#v", resp)
	}
	if err == nil {
		t.Errorf("expected error, got nil")
	} else if want, have := "foobar invalid request", err.Error(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}

func TestRouter_notFound(t *testing.T) {
	r := pubsub.NewRouter()
	resp, err := r.ServeRequest(context.Background(), pubsub.Request{
		ID:     "hello-id",
		Group:  "test",
		Entity: "ass",
		Method: "kick",
	})
	if resp != nil {
		t.Errorf("expected nil, got %#v", resp)
	}
	if err == nil {
		t.Errorf("expected error, got nil")
	} else if want, have := "invalid request", err.Error(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}
