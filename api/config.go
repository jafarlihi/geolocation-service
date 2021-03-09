package main

import (
	"encoding/json"
	"os"
)

type serverConfig struct {
	HTTPPort            int    `json:"httpPort"`
	TriggerImportSocket string `json:"triggerImportSocket"`
}

type mongoConfig struct {
	URL          string `json:"url"`
	Database     string `json:"database"`
	Collection   string `json:"collection"`
	DropOnUpdate bool   `json:"dropOnUpdate"`
}

type importFileConfig struct {
	Path string `json:"path"`
}

type configuration struct {
	Server     serverConfig     `json:"server"`
	Mongo      mongoConfig      `json:"mongo"`
	ImportFile importFileConfig `json:"importFile"`
}

// Config is a global variable holding application configuration
var Config configuration

// initConfig reads and parses the configuration file and makes the information available through the global Config variable
func initConfig() {
	configFile, err := os.Open("./config.json")
	if err != nil {
		Log.Errorf("Failed to open the config file: %v", err)
		os.Exit(1)
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&Config)
	if err != nil {
		Log.Errorf("Failed to decode the config file: %v", err)
		os.Exit(1)
	}
}
