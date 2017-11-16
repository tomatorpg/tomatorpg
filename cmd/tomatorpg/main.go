package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"golang.org/x/tools/godoc/vfs/httpfs"

	kitlog "github.com/go-kit/kit/log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/joho/godotenv"
	"github.com/tomatorpg/tomatorpg/assets"
	"github.com/tomatorpg/tomatorpg/models"
	"github.com/tomatorpg/tomatorpg/protocol/pubsub"
	"github.com/tomatorpg/tomatorpg/userauth"
	"github.com/tomatorpg/tomatorpg/utils"
)

var port uint64
var webpackDevHost string
var publicURL string
var jwtSecret string

var logger *log.Logger

func init() {
	// TODO; detect if is in heroku, skip timestamp
	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds)
}

func init() {

	var err error

	// load dot env file, if exists
	if _, err = os.Stat(".env"); err == nil {
		if err = godotenv.Load(); err != nil {
			logger.Fatalf("Unable to load .env, %#v", err)
		}
	}

	// load port
	if port, err = strconv.ParseUint(os.Getenv("PORT"), 10, 16); os.Getenv("PORT") != "" && err != nil {
		logger.Fatalf("Unable to parse PORT: %s", err.Error())
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

	// load JWT secret
	if jwtSecret = os.Getenv("JWT_SECRET"); jwtSecret == "" {
		jwtSecret = "abcdef"
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
	if err := models.AutoMigrate(db); err != nil {
		logger.Printf("error: %s", err.Error())
		panic("failed to migrate database")
	} else {
		logger.Printf("database migration done")
	}

	// login cookies
	genLoginCookie := func(r *http.Request) *http.Cookie {
		return &http.Cookie{
			Name:     "tomatorpg-token",
			Path:     "/",
			HttpOnly: true,
			Expires:  time.Now().Add(7 * 24 * time.Hour), // 7 days later
		}
	}

	// websocket pubsub server
	pubsubServer := pubsub.NewServer(
		db,
		make(pubsub.WebsocketChanColl),
		pubsub.RPCs(),
		jwtSecret,
	)

	// get auth providers from os environment
	authProviders := userauth.EnvProviders(os.Getenv, "/oauth2")
	if len(authProviders) == 0 {
		logger.Print("warning: No authentication provider is properly setup. Please setup at least one.")
	}

	// Create a simple file server
	fs := http.FileServer(httpfs.New(assets.FileSystem()))
	mainServer := http.NewServeMux()
	mainServer.Handle("/assets/js/", http.StripPrefix("/assets", fs))
	mainServer.Handle("/assets/images/", http.StripPrefix("/assets", fs))
	mainServer.Handle("/assets/css/", http.StripPrefix("/assets", fs))
	mainServer.Handle("/", handlePage(
		"index.html",
		struct{ ScriptPath string }{ScriptPath: webpackDevHost},
	))
	mainServer.Handle("/oauth2/", userauth.LoginHandler(
		db,
		genLoginCookie,
		authProviders,
		jwtSecret,
		publicURL,
		"/oauth2",
		"/",
		"/oauth2/error",
	))
	mainServer.Handle("/oauth2/login", handlePage(
		"login.html",
		struct {
			PageTitle       string
			PageHeaderTitle string
			BasePath        string
			Actions         []userauth.AuthProvider
		}{
			PageTitle:       "TomatoRPG | Login",
			PageHeaderTitle: "Login TomatoRPG",
			BasePath:        "/oauth2",
			Actions:         authProviders,
		},
	))
	mainServer.Handle("/oauth2/logout",
		userauth.LogoutHandler("/", genLoginCookie))
	mainServer.Handle("/api.v1", pubsubServer)

	// some custom reroutes
	reroutes := func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/favicon.ico" {
				r.URL.Path = "/assets/images/favicon.ico"
			}
			inner.ServeHTTP(w, r)
		})
	}

	applyMiddlewares := utils.Chain(
		utils.ApplyRequestID,
		utils.ApplyLogger(func() kitlog.Logger {
			return kitlog.NewLogfmtLogger(utils.LogWriter(logger))
		}),
		reroutes,
	)

	logger.Printf("listen to port %d", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), applyMiddlewares(mainServer))
	if err != nil {
		logger.Fatal("ListenAndServe: ", err)
	}
}
