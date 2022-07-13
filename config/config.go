package config

import (
	"log"
	"os"
)

const REFRESH_CONFIG_DATA string = "refresh"

// Config variable keys
const (
	// Host system info
	HOST              string = "HOST"
	PORT              string = "PORT"
	TRIVIASERVICENAME string = "TRIVIA_SERVICE_NAME"
	TRIVIASERVICEPORT string = "TRIVIA_SERVICE_PORT"
)

type ConfigData struct {
	Host              string `json:"host"`
	Port              string `json:"port"`
	TriviaServiceName string `json:"triviaservicename"`
	TriviaServicePort string `json:"triviaserviceport"`
}

type config struct {
	cfgData *ConfigData
}

var cfg *config

// Unexported type functions
func (c *config) getConfigEnv() error {
	// Loading config environment variables
	log.Print("loading config environment variables...")

	// Update config data
	c.cfgData.Host = os.Getenv(HOST)
	c.cfgData.Port = os.Getenv(PORT)
	c.cfgData.TriviaServiceName = os.Getenv(TRIVIASERVICENAME)
	c.cfgData.TriviaServicePort = os.Getenv(TRIVIASERVICEPORT)

	return nil
}

func (c *config) GetData(args ...string) (*ConfigData, error) {
	if len(args) > 0 && args[0] == REFRESH_CONFIG_DATA {
		log.Print("Using config environment to load config")

		getErr := cfg.getConfigEnv()
		if getErr != nil {
			log.Print("Error getting config environment data: ", getErr)
			return nil, getErr
		}
	}

	return c.cfgData, nil
}

func Get() *config {
	if cfg == nil {
		log.Print("creating config object")

		// Initialize config
		cfg = new(config)

		// Initialize config data
		cfg.cfgData = new(ConfigData)
	}

	log.Print("returning config object")
	return cfg
}
