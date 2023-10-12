package main

import (
	"github.com/Jero075/gotables/fs"
	"github.com/Jero075/gotables/server"
	"log"
)

func main() {
	config, err := fs.Config()
	if err != nil {
		return
	}
	go server.Run(config)
	if config.HTTPSMode {
		log.Println("Started server at " + "https://127.0.0.1" + config.Port)
	} else {
		log.Println("Started server at " + "http://127.0.0.1" + config.Port)
	}
	// log.Println("Logs are stored in " + config.LogDir)
	log.Println("Press 'ctrl' + 'C' to stop this program")
	for {

	}
}
