package pubsub

import (
	"context"
	"net/http"

	"github.com/go-restit/lzjson"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	"github.com/tomatorpg/tomatorpg/models"
)

type contextKey int

const (
	dbKey contextKey = iota
	jsonReqKey
	sessionKey
	srvKey
)

// Session store connection session information
// for a pubsub websocket sesions
type Session struct {
	HTTPRequest *http.Request
	Room        *RoomChannel
	User        models.User
	Conn        *websocket.Conn
}

// WithSession stores a websocket connection reference into context
func WithSession(parent context.Context, sess *Session) context.Context {
	return context.WithValue(parent, sessionKey, sess)
}

// GetSession get a websocket connection reference from context, if any
func GetSession(ctx context.Context) (sess *Session) {
	sess, _ = ctx.Value(sessionKey).(*Session)
	return
}

// WithDB stores a *gorm.DB into context
func WithDB(parent context.Context, db *gorm.DB) context.Context {
	return context.WithValue(parent, dbKey, db)
}

// GetDB retrieve *gorm.DB from context
func GetDB(ctx context.Context) (db *gorm.DB) {
	db, _ = ctx.Value(dbKey).(*gorm.DB)
	return
}

// WithJSONReq stores a websocket connection reference into context
func WithJSONReq(parent context.Context, node lzjson.Node) context.Context {
	return context.WithValue(parent, jsonReqKey, node)
}

// GetJSONReq get a websocket connection reference from context, if any
func GetJSONReq(ctx context.Context) (node lzjson.Node) {
	node, _ = ctx.Value(jsonReqKey).(lzjson.Node)
	return
}

// WithServer stores pointer to the server struct
func WithServer(parent context.Context, srv *Server) context.Context {
	return context.WithValue(parent, srvKey, srv)
}

// GetServer get the server struct reference from context
func GetServer(ctx context.Context) (srv *Server) {
	srv, _ = ctx.Value(srvKey).(*Server)
	return
}
