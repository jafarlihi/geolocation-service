module github.com/jafarlihi/geolocation-service/api

go 1.16

replace github.com/jafarlihi/geolocation-service/dataservice => ../dataservice

require (
	github.com/jafarlihi/geolocation-service/dataservice v0.0.0-00010101000000-000000000000
	github.com/op/go-logging v0.0.0-20160315200505-970db520ece7
)
