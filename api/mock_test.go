package main

import (
	"github.com/jafarlihi/geolocation-service/dataservice"
)

type dataServiceMock struct {
	importDataCalled    bool
	fetchLocationCalled bool
}

func (ds *dataServiceMock) ImportData(dropExistingRecords bool) (*dataservice.ImportStatistics, error) {
	ds.importDataCalled = true
	return nil, nil
}

func (ds *dataServiceMock) FetchLocation(ipAddress string) (*dataservice.Location, error) {
	ds.fetchLocationCalled = true
	location := dataservice.Location{IPAddress: "8.8.8.8", Latitude: 50.50}
	return &location, nil
}
