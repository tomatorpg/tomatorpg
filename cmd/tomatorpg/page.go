package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/tomatorpg/tomatorpg/assets"
)

var port uint64
var tplIndex *template.Template
var webpackDevHost string

// Configure the upgrader
var upgrader = websocket.Upgrader{}

func init() {

	var err error

	tplBin, err := assets.Asset("html/index.html")
	if err != nil {
		log.Fatalf("cannot find index.html in assets")
	}

	tplIndex, err = template.New("index").Parse(string(tplBin))
	if err != nil {
		log.Fatal(err)
	}

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
}

func handlePage(scriptPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			ScriptPath string
		}{
			ScriptPath: scriptPath,
		}
		tplIndex.Execute(w, data)
	}
}
