package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/tools/godoc/vfs/httpfs"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/tomatorpg/tomatorpg/assets"
	"github.com/tomatorpg/tomatorpg/userauth"
)

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

	// Create a simple file server
	fs := http.FileServer(httpfs.New(assets.FileSystem()))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.HandleFunc("/", handlePage(webpackDevHost))
	http.HandleFunc("/oauth2/google", func(w http.ResponseWriter, r *http.Request) {
		url := userauth.GoogleConfig("http://localhost:8080").AuthCodeURL("state", oauth2.AccessTypeOffline)
		http.Redirect(w, r, url, http.StatusFound)
	})
	http.HandleFunc("/oauth2/google/callback", userauth.GoogleCallback(
		userauth.GoogleConfig("http://localhost:8080"),
		db,
	))

	// Configure websocket route
	http.HandleFunc("/api.v1", handleConnections)

	log.Printf("listen to port %d", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
