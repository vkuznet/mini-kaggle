package main

// web server module
//
// Copyright (c) 2019 - Valentin Kuznetsov <vkuznet@gmail.com>
//
import (
	"fmt"
	"log"
	"net/http"
)

var VERBOSE int
var ScoreFile string
var _top, _bottom, _search string

func Server(configFile string) {
	err := ParseConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}
	VERBOSE = Config.Verbose
	ScoreFile = Config.ScoreFile
	var templates ServerTemplates
	tmplData := make(map[string]interface{})
	tmplData["Bottom"] = "&#169;, Valentin Kuznetsov, 2019"
	_top = templates.Top(Config.Templates, tmplData)
	_bottom = templates.Bottom(Config.Templates, tmplData)

	// Initialize ScoresDB
	ScoresDB, err = InitScoresDB(Config.Uri)

	// http handlers
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir(Config.Styles))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir(Config.Jscripts))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir(Config.Images))))
	http.HandleFunc("/dashboard", DashboardHandler)
	http.HandleFunc("/upload", UploadHandler)
	http.HandleFunc("/", HomeHandler)
	http.ListenAndServe(fmt.Sprintf(":%d", Config.Port), nil)
}
