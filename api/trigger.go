package main

import (
	"net"
	"syscall"

	"github.com/jafarlihi/geolocation-service/dataservice"
)

// listenForTriggerSocketConnections listens on the configured socket and triggers data import process when a connection initiates
func listenForTriggerSocketConnections(ds dataservice.IDataService) {
	syscall.Unlink(Config.Server.TriggerImportSocket)

	l, err := net.Listen("unix", Config.Server.TriggerImportSocket)
	if err != nil {
		Log.Errorf("Failed to listen on import trigger socket: %v", err)
		return
	}
	defer l.Close()

	for {
		_, err := l.Accept()
		if err != nil {
			Log.Errorf("Failed to accept connection on import trigger socket: %v", err)
			continue
		}

		Log.Info("Running data import because a trigger connection was received")
		importStatistics, err := ds.ImportData(Config.Mongo.DropOnUpdate)
		if err != nil {
			Log.Errorf("Could not complete the data import: %v", err)
			continue
		}
		Log.Infof("Data import complete: %+v", importStatistics)
	}
}
