package main

import (
	"html/template"
	"net/http"

	"github.com/tomatorpg/tomatorpg/assets"
)

var tpls map[string]*template.Template

func init() {
	var err error
	var tplBin []byte

	tpls = make(map[string]*template.Template)
	files, err := assets.FileSystem().ReadDir("/html")
	if err != nil {
		logger.Fatalf("error reading directory %s in asset: %s",
			"/html", err.Error())
	}

	for _, file := range files {
		if !file.IsDir() {
			tplBin, _ = assets.Asset("html/" + file.Name())
			tpls[file.Name()], err = template.New(file.Name()).Parse(string(tplBin))
			if err != nil {
				logger.Fatalf("error parsing template %s: %s",
					file.Name(), err.Error())
			}
		}
		// TODO: do this recursively for directories, maybe
	}
}

func handlePage(tplPath string, data interface{}) http.HandlerFunc {
	tpl, ok := tpls[tplPath]
	if !ok {
		logger.Fatalf("template %s not found", tplPath)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html; charset=utf8")
		w.WriteHeader(http.StatusOK)
		tpl.Execute(w, data)
	}
}
