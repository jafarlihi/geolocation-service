# geolocation-service

`geolocation-service` is made up of two parts: `dataservice` library and `api` application.

`dataservice` library is resposible for importing the data from CSV to the database and providing an interface for fetching stored location information of a given IP address.

`api` application makes use of the `dataservice` library and provides a REST endpoint for fetching location information of a given IP address and also listens on a Unix socket and triggers the data import process when connections are made to the socket. Running `api` also triggers initial data import process but this can be avoided with `-noInitialImport` flag.

Records are stored in MongoDB. An unique index is created on `ipAddress` field to avoid inserting more than one location information for any given IP address.

## Running
### Locally
1. Edit ./api/config.json
2. Run `cd ./api && go build && ./api`

### Docker
(Data import process is a lot slower in Docker)
1. Edit ./config.json
2. Run `sudo docker-compose up`

## dataservice
### Interface
`dataService` struct instances are the entry point to the library.
To create one, call `NewDataService` function. Example:
```go
ds, err := dataservice.NewDataService(dataservice.Settings{
    ImportFilePath:  "./data_dump.csv",
    MongoURL:        "mongodb://localhost:27017",
    MongoDB:         "geolocation",
    MongoCollection: "location",
})
```
Then you'll have two methods available to call on your `dataService` instance: `ImportData` and `FetchLocation`.
- `FetchLocation` takes in a single string argument indicating the IP address to fetch the location information for. It returns an instance of `Location` struct when an entry for the given IP address is found, error otherwise. `Location` struct is defined as:
```go
type Location struct {
	IPAddress    string  `bson:"ipAddress"`
	CountryCode  string  `bson:"countryCode"`
	Country      string  `bson:"country"`
	City         string  `bson:"city"`
	Latitude     float64 `bson:"latitude"`
	Longitude    float64 `bson:"longitude"`
	MysteryValue int64   `bson:"mysteryValue"`
}
```
- `ImportData` reads the CSV file and inserts the records into the database, it takes in a single boolean argument called `dropExistingRecords`. When `true` is passed the existing records in the database are dropped each time before data import process runs. If `false` is passed then old records are preserved and new records from the initiated data import are inserted alongside them. `ImportData` returns an instance of `ImportStatistics` struct, which is defined as:
```go
type ImportStatistics struct {
	Duration            time.Duration
	AcceptedRecordCount int
	RejectedRecordCount int
}
```

## api
### Configuration
Here's a sample configuration file, `config.json`:
```json
{
    "server": {
        "httpPort": 8080,
        "triggerImportSocket": "/tmp/geo.sock"
    },
    "mongo": {
        "url": "mongodb://mongo:27017",
        "database": "geolocation",
        "collection": "location",
        "dropOnUpdate": true
    },
    "importFile": {
        "path": "./data_dump.csv"
    }
}
```
- `triggerImportSocket` is the location of Unix socket that is being listened for connections, when a connection is initiated the data import process runs. This provides an easy interface to trigger the data import process using cron jobs. A connection can be made to the Unix socket using `socat` or `netcat-openbsd`:
```bash
socat - UNIX-CONNECT:/tmp/geo.sock
nc -U /tmp/geo.sock
```
- `dropOnUpdate` passes the value to `dropExistingRecords` argument of `ImportData` method of `dataservice` library.

### Command-line flags
- `-noInitialImport` flag can be passed to `api` to prevent the application from running the initial data import. If this flag is not present then the data import process will run each time the application is started.

### REST API
#### GET /location
A single query parameter indicating the IP address is expected, under key "ipAddress".
Example: `GET /location?ipAddress=8.8.8.8`
