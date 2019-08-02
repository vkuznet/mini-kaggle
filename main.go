package main

// mini-kaggle dashboard app based on gonum metrics (AUC)
//
// Copyright (c) 2019 - Valentin Kuznetsov <vkuznet@gmail.com>
//

import (
	"flag"
)

func main() {
	var config string
	flag.StringVar(&config, "config", "server.json", "server config JSON file")
	flag.Parse()
	Server(config)
}
