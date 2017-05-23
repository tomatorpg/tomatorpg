package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/tomatorpg/tomatorpg/assets"
)

var tplIndex *template.Template

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
