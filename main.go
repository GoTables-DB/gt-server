package main

import (
	"git.jereileu.ch/gotables/server/gt-server/fs"
	"git.jereileu.ch/gotables/server/gt-server/server"
	"log"
	"os"
)

const (
	ConfigEnvvar   = "GT_CONFIG"
	Version        = "0.2.2"
	CopyrightYear  = "2024"
	CopyrightName  = "Jeroen Leuenberger"
	CopyrightEmail = "jereileu@proton.me"
)

func main() {
	// Load config.json
	location := os.Getenv(ConfigEnvvar)
	config, err := fs.Config(location)
	if err != nil {
		log.Fatal(err)
	}
	// Start server
	log.Println("==================== GoTables server " + Version + " ====================")
	log.Println("Copyright © " + CopyrightYear + " " + CopyrightName + " <" + CopyrightEmail + ">")
	if config.HTTPSMode {
		log.Println("Started server at " + "https://127.0.0.1" + config.Port)
	} else {
		log.Println("Started server at " + "http://127.0.0.1" + config.Port)
	}
	log.Println("Logs are stored in " + config.LogDir)
	log.Println("Press 'Ctrl' + 'C' to stop this program")
	end := ""
	for i := 0; i < 58+len(Version); i++ {
		end += "="
	}
	log.Println(end)
	log.Println("")
	server.Run(config)
}
