package main

// configuration module
//
// Copyright (c) 2019 - Valentin Kuznetsov <vkuznet@gmail.com>
//
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// Configuration stores server configuration parameters
type Configuration struct {
	Port             int    `json:"port"`             // server port number
	Uri              string `json:"uri"`              // server scoresdb URI
	Templates        string `json:"templates"`        // location of server templates
	Jscripts         string `json:"jscripts"`         // location of server JavaScript files
	Images           string `json:"images"`           // location of server images
	Styles           string `json:"styles"`           // location of server CSS styles
	ScoreFile        string `json:"scoreFile"`        // score file
	PrivateScoreFile string `json:"privateScoreFile"` // private score file
	Destination      string `json:"destination"`      // location of submissions
	Verbose          int    `json:"verbose"`          // verbosity level
}

// Config variable represents configuration object
var Config Configuration

// String returns string representation of server Config
func (c *Configuration) String() string {
	return fmt.Sprintf("<Config port=%d templates=%s js=%s images=%s css=%s scores=%s privateScores=%s dst=%s>", c.Port, c.Templates, c.Jscripts, c.Images, c.Styles, c.ScoreFile, c.PrivateScoreFile, c.Destination)
}

// helper function to return full path of given file
func path(fname string) string {
	if !strings.HasPrefix(fname, "/") {
		wdir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		fname = fmt.Sprintf("%s/%s", wdir, fname)
	}
	if _, err := os.Stat(fname); err != nil {
		log.Fatal(err)
	}
	return fname
}

// ParseConfig parse given config file
func ParseConfig(configFile string) error {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Println("config", configFile, "error", err)
		return err
	}
	err = json.Unmarshal(data, &Config)
	if err != nil {
		log.Println("config", configFile, "error", err)
		return err
	}
	// make sure that all paths exists
	Config.Templates = path(Config.Templates)
	Config.Jscripts = path(Config.Jscripts)
	Config.Styles = path(Config.Styles)
	Config.Images = path(Config.Images)
	Config.ScoreFile = path(Config.ScoreFile)
	Config.PrivateScoreFile = path(Config.PrivateScoreFile)
	Config.Destination = path(Config.Destination)
	return nil
}
