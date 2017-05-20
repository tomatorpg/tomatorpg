package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"gopkg.in/jose.v1/crypto"
	"gopkg.in/jose.v1/jws"

	"golang.org/x/oauth2"
	"golang.org/x/tools/godoc/vfs/httpfs"

	"github.com/go-restit/lzjson"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/joho/godotenv"
	"github.com/tomatorpg/tomatorpg/assets"
	"github.com/tomatorpg/tomatorpg/models"
	"github.com/tomatorpg/tomatorpg/userauth"
)

var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan Action)            // broadcast channel

// Configure the upgrader
var upgrader = websocket.Upgrader{}

// Action object
type Action struct {
	Entity    string    `json:"entity"`
	Action    string    `json:"action"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
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

var port uint64

func handleConnections(w http.ResponseWriter, r *http.Request) {

	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	// Register our new client
	clients[ws] = true

	for {
		var msg Action
		// Read in a new Action as JSON and map it to a Action object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}

		switch msg.Action {
		case "":
			log.Printf("message: %s", msg.Message)
		case "sign_in":
			log.Printf("sign in")
		}
		// Send the newly received message to the broadcast channel
		broadcast <- msg
	}
}

func handleActions() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
		// Send it out to every client that is currently connected
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func handlePage(scriptPath string) http.HandlerFunc {
	tplBin, err := assets.Asset("html/index.html")
	if err != nil {
		log.Fatalf("cannot find index.html in assets")
	}

	t, err := template.New("index").Parse(string(tplBin))
	if err != nil {
		log.Fatal(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			ScriptPath string
		}{
			ScriptPath: scriptPath,
		}
		t.Execute(w, data)
	}
}

func initDB(db *gorm.DB) {

	log.Printf("initDB")

	// Migrate the schema
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.UserEmail{})

	// Create
	//db.Create(&models.User{
	//	Email: "hello+" + time.Now().Format("20060102-150405") + "@world.com",
	//})
}

func main() {

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
	webpackDevHost := ""
	if os.Getenv("NODE_ENV") == "development" {
		if webpackDevHost = os.Getenv("WEBPACK_DEV_SERVER_HOST"); webpackDevHost == "" {
			webpackDevHost = "http://localhost:8081" // default, if not set
		}
	}

	// connect to database
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	// initialize database / migration
	// TODO: make this optional on start up
	initDB(db)

	// Create a simple file server
	fs := http.FileServer(httpfs.New(assets.FileSystem()))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.HandleFunc("/", handlePage(webpackDevHost))
	http.HandleFunc("/oauth2/google", func(w http.ResponseWriter, r *http.Request) {
		url := userauth.ConfigGoogle("http://localhost:8080").AuthCodeURL("state", oauth2.AccessTypeOffline)
		http.Redirect(w, r, url, http.StatusFound)
	})
	http.HandleFunc("/oauth2/google/callback", func(w http.ResponseWriter, r *http.Request) {
		conf := userauth.ConfigGoogle("http://localhost:8080")
		code := r.URL.Query().Get("code")
		token, err := conf.Exchange(oauth2.NoContext, code)
		if err != nil {
			log.Printf("Code exchange failed with '%s'\n", err.Error())
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		client := conf.Client(context.Background(), token)
		resp, err := client.Get("https://www.googleapis.com/oauth2/v1/userinfo")
		if err != nil {
			return
		}

		// read into
		/*
			// NOTE: JSON structure of normal response body
			{
			 "id": "some-id-in-google-account",
			 "email": "email-for-the-account",
			 "verified_email": true,
			 "name": "Some Name",
			 "given_name": "Some",
			 "family_name": "Name",
			 "link": "https://plus.google.com/+SomeUserOnGPlus",
			 "picture": "url-to-some-picture",
			 "gender": "female",
			 "locale": "zh-HK"
			}
		*/

		result := lzjson.Decode(resp.Body)
		// TODO: detect read  / decode error
		// TODO: check if the email has been verified or not
		authUser := models.User{
			Name:         result.Get("name").String(),
			PrimaryEmail: result.Get("email").String(),
		}

		// search existing user with the email
		var userEmail models.UserEmail
		var prevUser models.User

		if db.First(&prevUser, "primary_email = ?", authUser.PrimaryEmail); prevUser.PrimaryEmail != "" {
			// TODO: log this?
			authUser = prevUser
		} else if db.First(&userEmail, "email = ?", authUser.PrimaryEmail); userEmail.Email != "" {
			// TODO: log this?
			db.First(&authUser, "id = ?", userEmail.UserID)
		} else {

			tx := db.Begin()

			// create user
			if res := tx.Create(&authUser); res.Error != nil {
				// TODO: log and provide error to user
				tx.Rollback()
				return
			}

			// create user-email relation
			newUserEmail := models.UserEmail{
				UserID: authUser.ID,
				Email:  authUser.PrimaryEmail,
			}
			if res := tx.Create(&newUserEmail); res.Error != nil {
				tx.Rollback()
				return
			}

			tx.Commit()
		}
		log.Printf("user found or created: %#v", authUser)

		// Create JWS claims with the user info
		claims := jws.Claims{}
		claims.Set("id", authUser.ID)
		claims.Set("name", authUser.Name)
		claims.SetAudience("localhost")

		jwtToken := jws.NewJWT(claims, crypto.SigningMethodHS256)
		serializedToken, _ := jwtToken.Serialize([]byte("abcdef"))

		http.Redirect(w, r, "http://localhost:8080?token="+string(serializedToken), http.StatusFound)
	})

	// Configure websocket route
	http.HandleFunc("/api.v1", handleConnections)

	// Start listening for incoming chat messages
	go handleActions()

	log.Printf("listen to port %d", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
