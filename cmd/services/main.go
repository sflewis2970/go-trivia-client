package main

import (
	"log"
	"net/http"

	"github.com/sflewis2970/go-trivia-client/config"
	"github.com/sflewis2970/go-trivia-client/controllers"
)

func main() {
	// Initialize logging
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Get config data
	cfgData, getCfgDataErr := config.Get().GetData(config.REFRESH_CONFIG_DATA)
	if getCfgDataErr != nil {
		log.Fatal("Error getting config data: ", getCfgDataErr)
	}

	// Create controller
	controllers.New()

	// Server Address info
	addr := cfgData.HostName + cfgData.HostPort
	log.Print("The address used by the service is: ", addr)

	log.Print("Web Server is ready!")
	log.Fatal(http.ListenAndServe(addr, nil))
}
