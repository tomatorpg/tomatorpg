package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/tools/godoc/vfs/httpfs"

	kitlog "github.com/go-kit/kit/log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/joho/godotenv"
	"github.com/tomatorpg/tomatorpg/assets"
	"github.com/tomatorpg/tomatorpg/protocol/pubsub"
	"github.com/tomatorpg/tomatorpg/userauth"
)

var port uint64
var webpackDevHost string
var publicURL string

func init() {

	var err error

	// load dot env file, if exists
	if _, err = os.Stat(".env"); err == nil {
		if err = godotenv.Load(); err != nil {
			log.Fatalf("Unable to load .env, %#v", err)
		}
	}

	// load port
	if port, err = strconv.ParseUint(os.Getenv("PORT"), 10, 16); os.Getenv("PORT") != "" && err != nil {
		log.Fatalf("Unable to parse PORT: %s", err.Error())
		return
	}
	if port == 0 {
		port = 8080
	}

	// check if in development mode
	// if so, try to load webpack dev server host
	if os.Getenv("NODE_ENV") == "development" {
		if webpackDevHost = os.Getenv("WEBPACK_DEV_SERVER_HOST"); webpackDevHost == "" {
			webpackDevHost = "http://localhost:8081" // default, if not set
		}
	}

	// load public url for OAuth2 redirect
	if publicURL = os.Getenv("PUBLIC_URL"); publicURL == "" {
		publicURL = "http://localhost:8080"
	}
}

/**

TODO:
1. Session to store user information
2. Session ID to resume on disconnect
3. Session to be able to store any JSON object payload to use in JS
4. Create room
5. Join room by room id
6. Room listing
7. Room history save and load (only show limited row backward)
8. Room status snapshot to prevent need to read whole history to build current status

Advanced:
1. Operational Transformation?
*/

func main() {

	// connect to database
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	// initialize database / migration
	// TODO: make this optional on start up
	initDB(db)

	// websocket pubsub server
	pubsubServer := pubsub.NewServer(
		db,
		make(pubsub.WebsocketChanColl),
		pubsub.RPCs(),
	)

	// Create a simple file server
	fs := http.FileServer(httpfs.New(assets.FileSystem()))
	http.Handle("/assets/js/", http.StripPrefix("/assets", fs))
	http.HandleFunc("/", handlePage(webpackDevHost))
	http.HandleFunc("/oauth2/google", func(w http.ResponseWriter, r *http.Request) {
		url := userauth.GoogleConfig(publicURL).AuthCodeURL("state", oauth2.AccessTypeOffline)
		http.Redirect(w, r, url, http.StatusFound)
	})
	http.HandleFunc("/oauth2/google/callback", userauth.GoogleCallback(
		userauth.GoogleConfig(publicURL),
		db,
	))
	http.HandleFunc("/oauth2/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:    "tomatorpg-token",
			Path:    "/",
			Expires: time.Now().Add(-1 * time.Hour),
		})
		http.Redirect(w, r, "/", http.StatusFound)
	})
	applyMiddlewares := pubsub.Chain(
		pubsub.ApplyRequestID,
		pubsub.ApplyContextLog(kitlog.NewLogfmtLogger(os.Stdout)),
	)
	http.Handle("/api.v1", applyMiddlewares(pubsubServer))

	log.Printf("listen to port %d", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
