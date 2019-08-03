package main

// web server handlers module
//
// Copyright (c) 2019 - Valentin Kuznetsov <vkuznet@gmail.com>
//
import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func logRequest(r *http.Request) {
	log.Println(r.Method, r.URL, r.RemoteAddr)
}

// UploadHandler uploads predictions to the server (/upload API)
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()
	logRequest(r)
	if _verbose > 0 {
		log.Println("UploadHandler: Header", r.Header)
	}

	// read name of the prediction file
	name := r.FormValue("name")
	if name == "" {
		log.Println("UploadHandler name is not provided")
		http.Error(w, "Please provide name for your submission", http.StatusInternalServerError)
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		log.Println("UploadHandler unable to open provided file", err, "request:", r)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// create destination file where we will put uploaded one
	destination := Config.Destination
	if !strings.HasSuffix(destination, "/") {
		destination += "/"
	}
	dstFileName := destination + fmt.Sprintf("%s-%s", name, header.Filename)
	if strings.HasSuffix(dstFileName, ".gz") {
		dstFileName = strings.Replace(dstFileName, ".gz", "", -1)
	}
	dst, err := os.Create(dstFileName)
	defer dst.Close()
	if err != nil {
		log.Println("UploadHandler unable to create destination file", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// copy file to destination
	if strings.HasSuffix(header.Filename, ".gz") {
		// we got gzip'ed file and will create appropriate reader
		gzipReader, err := gzip.NewReader(file)
		if err != nil {
			log.Println("UploadHandler, unable to read gzip file", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer gzipReader.Close()
		_, err = io.Copy(dst, gzipReader)
	} else {
		_, err = io.Copy(dst, file)
	}
	if err != nil {
		log.Println("UploadHandler unable to copy file to destination", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// check HTTP header and return appropriate response
	var accept string
	if _, ok := r.Header["Accept"]; ok {
		accept = r.Header["Accept"][0]
	}
	if strings.Contains(accept, "json") {
		response := make(map[string]interface{})
		response["status"] = "ok"
		response["File"] = header.Filename
		js, err := json.Marshal(&response)
		if err != nil {
			log.Println("UploadHandler unable to marshal response", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		return
	}
	score, private := getScores(dstFileName)
	InsertScore(name, score, private)
	msg := fmt.Sprintf("Your file %s has been successfully uploaded, score: %f", header.Filename, score)
	var templates ServerTemplates
	tmplData := make(map[string]interface{})
	tmplData["Message"] = msg
	page := templates.Confirm(Config.Templates, tmplData)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(_top + page + _bottom))

}

// HomeHandler handles home page
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	logRequest(r)
	if _verbose > 0 {
		log.Println("RequestHandler", r)
	}
	var templates ServerTemplates
	tmplData := make(map[string]interface{})
	page := templates.Home(Config.Templates, tmplData)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(_top + page + _bottom))
}

// DashboardHandler handlers /status API
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	logRequest(r)
	if _verbose > 0 {
		log.Println("StatusHandler", r)
	}
	err := r.ParseForm()
	if err != nil {
		log.Println("DashboardHandler unable to parse input parameter", err)
	}
	stype := r.FormValue("stype")
	var templates ServerTemplates
	tmplData := make(map[string]interface{})
	tmplData["Type"] = strings.ToTitle(stype)
	tmplData["Records"] = GetScores(stype)
	page := templates.Dashboard(Config.Templates, tmplData)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(_top + page + _bottom))
}
