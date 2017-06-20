package pubsub

import (
	"math/rand"
	"net/http"
	"time"

	kitlog "github.com/go-kit/kit/log"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// HTTPMiddleware is signature to a plain http middleware
type HTTPMiddleware func(inner http.Handler) http.Handler

// Chain chains HTTPMiddleware to form a single middleware
func Chain(mwares ...HTTPMiddleware) HTTPMiddleware {
	return func(inner http.Handler) http.Handler {
		for i := len(mwares) - 1; i >= 0; i-- {
			inner = mwares[i](inner)
		}
		return inner
	}
}

// ApplyRequestID is a middleware to apply X-Request-ID to request header
// if it is not set
func ApplyRequestID(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if id := r.Header.Get("X-Request-ID"); id == "" {
			r.Header.Set("X-Request-ID", randStringRunes(20))
		}
		inner.ServeHTTP(w, r)
	})
}

// ApplyContextLog logs access and also provide the kitlog context to inner
// http handler
func ApplyContextLog(newlogger func() kitlog.Logger) HTTPMiddleware {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			reqID := r.Header.Get("X-Request-ID")
			logger := newlogger()
			logger = kitlog.With(
				logger,
				"request_id", reqID,
			)

			// access log
			logger.Log(
				"at", "info",
				"method", r.Method,
				"path", r.URL.Path,
				"protocol", r.URL.Scheme,
				"remote_addr", r.RemoteAddr,
			)

			inner.ServeHTTP(w, r.WithContext(WithLogContext(r.Context(), logger)))
		})
	}
}
