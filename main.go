package main

import (
	"git.jereileu.ch/gotables/server/gt-server/fs"
	"git.jereileu.ch/gotables/server/gt-server/server"
	"log"
)

func main() {
	config, err := fs.Config()
	if err != nil {
		log.Fatal(err)
	}
	go server.Run(config)
	if config.HTTPSMode {
		log.Println("Started server at " + "https://127.0.0.1" + config.Port)
	} else {
		log.Println("Started server at " + "http://127.0.0.1" + config.Port)
	}
	log.Println("Logs are stored in " + config.LogDir)
	log.Println("Press 'Ctrl' + 'C' to stop this program")
	for {

	}
}
