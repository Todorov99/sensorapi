package influx

import (
	"os"

	"github.com/Todorov99/sensorcli/pkg/logger"
	influxClient "github.com/influxdata/influxdb-client-go/v2"
)

var influxLogger = logger.NewLogrus("influx", os.Stdout)

const (
	address string = "http://influxdb:8086/"
	// DbName is influxdb name.
	DbName           string = "sensorCLI"
	username         string = "todor"
	influxDbpassword string = "Abcd1234!"
)

var InfluxdbClient influxClient.Client

//TODO property file should be implemented
const (
	token  = "testToken"
	Org    = "org"
	Bucket = "bucket"
)

func init() {
	influxLogger.Info("Initializing influx db 2.0 client")
	dbClient := influxClient.NewClient(address, token)
	InfluxdbClient = dbClient
}
