package main

import (
	"os"

	"github.com/op/go-logging"
)

// Log is a global variable through which all logging is done
var Log = logging.MustGetLogger("geolocation-service")

// initLogger initializes the logging backend
func initLogger() {
	backend := logging.NewLogBackend(os.Stdout, "", 0)
	formatter := logging.NewBackendFormatter(backend, logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	))
	logging.SetBackend(formatter)
}
