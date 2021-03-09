package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/jafarlihi/geolocation-service/dataservice"
)

// dataServiceWrapper wraps the dataService and makes it available to handlers by acting as a receiver the handlers are methods of
type dataServiceWrapper struct {
	ds dataservice.IDataService
}

// serveRetrievalRequest services GET requests and returns location information about the IP address indicated in the query parameters
func (dsw *dataServiceWrapper) serveRetrievalRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	ipAddress := r.URL.Query().Get("ipAddress")
	if ipAddress == "" {
		http.Error(w, "ipAddress query parameter is required", http.StatusBadRequest)
		return
	}
	Log.Infof("Received a retrieval request for IP address %s", ipAddress)

	location, err := dsw.ds.FetchLocation(ipAddress)
	if err != nil {
		if strings.Contains(err.Error(), "no documents in result") {
			Log.Infof("No location information was found for IP address %s", ipAddress)
			http.Error(w, err.Error(), http.StatusNoContent)
		} else {
			Log.Errorf("Fetching location failed for IP address %s: %v", ipAddress, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	Log.Infof("Fetched location for IP address %s: %+v", ipAddress, location)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(location)
}
