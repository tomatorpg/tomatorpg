package pubsub_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/tomatorpg/tomatorpg/protocol/pubsub"
)

func TestRouter(t *testing.T) {
	r := pubsub.NewRouter()
	r.Add("test", "world", "hello", func(ctx context.Context, req pubsub.Request) pubsub.Response {
		return pubsub.SuccessResponseTo(req, "hello world")
	})
	r.Add("test", "foo", "bar", func(ctx context.Context, req pubsub.Request) pubsub.Response {
		return pubsub.ErrorResponseTo(req, fmt.Errorf("foo bar error"))
	})
	r.NotFound(func(ctx context.Context, req pubsub.Request) pubsub.Response {
		return pubsub.ErrorResponseTo(req, fmt.Errorf("foobar invalid request"))
	})

	resp := r.ServeRequest(context.Background(), pubsub.Request{
		ID:     "hello-id",
		Group:  "test",
		Entity: "world",
		Method: "hello",
	})
	if want, have := "hello-id", resp.ID; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "response", resp.Type; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "world", resp.Entity; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "hello", resp.Method; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "success", resp.Status; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "hello world", resp.Data; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	resp = r.ServeRequest(context.Background(), pubsub.Request{
		ID:     "hello-id",
		Group:  "test",
		Entity: "foo",
		Method: "bar",
	})
	if want, have := "hello-id", resp.ID; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "response", resp.Type; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "foo", resp.Entity; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "bar", resp.Method; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "error", resp.Status; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if resp.Err == nil {
		t.Errorf("expected error, got nil")
	} else if want, have := "foo bar error", resp.Err.Error(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	resp = r.ServeRequest(context.Background(), pubsub.Request{
		ID:     "hello-id",
		Group:  "test",
		Entity: "ass",
		Method: "kick",
	})
	if want, have := "hello-id", resp.ID; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "response", resp.Type; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "ass", resp.Entity; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "kick", resp.Method; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "error", resp.Status; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if resp.Err == nil {
		t.Errorf("expected error, got nil")
	} else if want, have := "foobar invalid request", resp.Err.Error(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}

func TestRouter_notFound(t *testing.T) {
	r := pubsub.NewRouter()
	resp := r.ServeRequest(context.Background(), pubsub.Request{
		ID:     "hello-id",
		Group:  "test",
		Entity: "ass",
		Method: "kick",
	})
	if want, have := "hello-id", resp.ID; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "response", resp.Type; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "ass", resp.Entity; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "kick", resp.Method; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "error", resp.Status; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if resp.Err == nil {
		t.Errorf("expected error, got nil")
	} else if want, have := "invalid request", resp.Err.Error(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}
