package utils

import (
	"context"
	"os"

	kitlog "github.com/go-kit/kit/log"
)

type contextKey int

const (
	logCtxKey contextKey = iota
)

// WithLogger stores a go-kit log *Context to a context.Context
func WithLogger(parent context.Context, logCtx kitlog.Logger) context.Context {
	return context.WithValue(parent, logCtxKey, logCtx)
}

// GetLogger get the go-kit log *Context from the context.Context
func GetLogger(ctx context.Context) (logCtx kitlog.Logger) {
	logCtx, _ = ctx.Value(logCtxKey).(kitlog.Logger)
	if logCtx == nil {
		logCtx = kitlog.NewLogfmtLogger(os.Stdout)
	}
	return
}
