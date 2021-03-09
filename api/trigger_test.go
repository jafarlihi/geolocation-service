package main

import (
	"net"
	"os"
	"testing"
	"time"
)

func TestTrigger(t *testing.T) {
	socketPath := "/tmp/geo.sock"

	os.Create(socketPath)
	Config = configuration{Server: serverConfig{TriggerImportSocket: socketPath}}
	dataServiceMock := &dataServiceMock{}

	go listenForTriggerSocketConnections(dataServiceMock)
	time.Sleep(100 * time.Millisecond) // Wait for listening on socket to be initialized

	c, err := net.Dial("unix", socketPath)
	if err != nil {
		t.Errorf("Failed to connect to the Unix socket: %v", err)
		return
	}
	defer c.Close()
	time.Sleep(100 * time.Millisecond) // Wait for trigger to initiate ImportData call

	if !dataServiceMock.importDataCalled {
		t.Errorf("Expected ImportData call was not found")
		return
	}
}
