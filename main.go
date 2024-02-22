package main

import (
	"git.jereileu.ch/gotables/server/gt-server/fs"
	"git.jereileu.ch/gotables/server/gt-server/server"
	"log"
)

const (
	CONFIG  = "" // Change this const to point to the config.json file
	VERSION = "0.2.0-beta2"
)

func main() {
	// Load config.json
	config, err := fs.Config(CONFIG)
	if err != nil {
		log.Fatal(err)
	}
	// Start server
	go server.Run(config)
	log.Println("==================== GoTables server " + VERSION + " ====================")
	if config.HTTPSMode {
		log.Println("Started server at " + "https://127.0.0.1" + config.Port)
	} else {
		log.Println("Started server at " + "http://127.0.0.1" + config.Port)
	}
	log.Println("Logs are stored in " + config.LogDir)
	log.Println("Press 'Ctrl' + 'C' to stop this program")
	end := ""
	for i := 0; i < 58+len(VERSION); i++ {
		end += "="
	}
	log.Println(end)
	log.Println("")
	for {

	}
}
