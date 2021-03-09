package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/jafarlihi/geolocation-service/dataservice"
)

func main() {
	initLogger() // Initialize the logging backend

	Log.Info("Processing config file")
	initConfig() // Initialize the config
	Log.Infof("Config file processed: %+v", Config)

	ds, err := dataservice.NewDataService(dataservice.Settings{
		ImportFilePath:  Config.ImportFile.Path,
		MongoURL:        Config.Mongo.URL,
		MongoDB:         Config.Mongo.Database,
		MongoCollection: Config.Mongo.Collection,
	}) // Create a new dataService instance
	if err != nil {
		Log.Errorf("Could not initialize the dataService: %v", err)
		os.Exit(1)
	}

	if len(os.Args) > 1 && os.Args[1] == "-noInitialImport" {
		Log.Info("Skipping initial data import")
	} else {
		Log.Info("Running initial data import")
		importStatistics, err := ds.ImportData(Config.Mongo.DropOnUpdate) // Trigger the initial data import
		if err != nil {
			Log.Errorf("Could not complete the initial data import: %v", err)
			os.Exit(1)
		}
		Log.Infof("Initial data import complete: %+v", importStatistics)
	}

	go listenForTriggerSocketConnections(ds) // Spawn a goroutine to listen on a socket for data import trigger requests

	dsw := &dataServiceWrapper{ds: ds}                      // Wrap the dataService so it can be accessed in the handlers
	http.HandleFunc("/location", dsw.serveRetrievalRequest) // Register the handler for fetching locations
	Log.Infof("Bringing HTTP server up at port %d", Config.Server.HTTPPort)
	Log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", Config.Server.HTTPPort), nil)) // Bring up the HTTP server
}
