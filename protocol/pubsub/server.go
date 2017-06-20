package pubsub

import (
	"encoding/json"
	"net/http"

	"github.com/go-restit/lzjson"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	"github.com/tomatorpg/tomatorpg/models"
)

// Server implements pubsub websocket server
type Server struct {
	db       *gorm.DB
	chans    ChanColl
	upgrader websocket.Upgrader
	router   *Router
}

// NewServer create pubsub http handler
func NewServer(db *gorm.DB, coll ChanColl, router *Router) *Server {
	return &Server{
		db:    db,
		chans: coll,
		upgrader: websocket.Upgrader{
			Subprotocols: []string{
				"tomatorpc-v1",
			},
		},
		router: router,
	}
}

// ServeHTTP implements http.Handler interface
func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// inherit logger from server
	logger := GetLogContext(r.Context())

	// Upgrade initial GET request to a websocket
	ws, err := srv.upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Log(
			"at", "error",
			"message", "unable to upgrade connection",
			"error", err.Error(),
		)
		respEnc := json.NewEncoder(w)
		w.WriteHeader(http.StatusBadRequest)
		respEnc.Encode(map[string]interface{}{
			"status":       "error",
			"error":        "unable to upgrade connection",
			"errorDetails": err.Error(),
		})
		return
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	logger.Log(
		"at", "info",
		"message", "connection upgrade success",
	)

	// context variables
	user := models.User{
		Name: "Visitor",
	}

	// load user from token
	if c, err := r.Cookie("tomatorpg-token"); err != nil {
		// TODO: detect error not found and ignore
		logger.Log(
			"at", "error",
			"message", "error reading token from cookie",
			"error", err.Error(),
		)
	} else if token, err := ParseToken("abcdef", c.Value); err != nil {
		logger.Log(
			"at", "error",
			"message", "error parsing / validating token",
			"error", err.Error(),
		)
	} else {
		// get user of the id
		srv.db.Find(&user, token.Claims().Get("id"))
		if user.ID != 0 {
			logger.Log(
				"at", "info",
				"message", "user loaded",
				"user_id", user.ID,
				"user_name", user.Name,
			)
		}
	}

	// session to be used and modified by procedures
	sess := &Session{
		HTTPRequest: r,
		User:        user,
		Conn:        ws,
		Logger:      logger,
	}

	// build common procedure context
	ctx := r.Context()
	ctx = WithDB(ctx, srv.db)
	ctx = WithServer(ctx, srv)
	ctx = WithSession(ctx, sess)

	for {

		// parse as JSON request for flexibility
		jsonRequest := lzjson.NewNode()
		err := ws.ReadJSON(&jsonRequest)
		if err != nil {
			switch terr := err.(type) {
			case *websocket.CloseError:
				sess.Logger.Log(
					"at", "info",
					"message", "websocket disconnected",
					"errCode", terr.Code,
					"error", terr.Text,
				)
			default:
				sess.Logger.Log(
					"at", "error",
					"message", "error reading JSON",
					"error", err.Error(),
				)
			}
			if sess.RoomChan != nil {
				sess.RoomChan.Unsubscribe(sess.Conn)
			}
			break
		}

		// parse and execute the RPC
		var req Request
		jsonRequest.Unmarshal(&req)
		reqCtx := WithJSONReq(ctx, jsonRequest)

		// handle all routes similarly
		resp, err := srv.router.ServeRequest(reqCtx, req)
		if err != nil {
			sess.Logger.Log(
				"at", "error",
				"message", "server request error",
				"error", err.Error(),
			)
			ws.WriteJSON(ErrorResponseTo(req, err))
			return
		}
		ws.WriteJSON(SuccessResponseTo(req, resp))
	}
}
