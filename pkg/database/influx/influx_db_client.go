package influx

import (
	"os"

	"github.com/Todorov99/sensorcli/pkg/logger"
	influxClient "github.com/influxdata/influxdb-client-go/v2"
)

var influxLogger = logger.NewLogrus("influx", os.Stdout)

const (
	address string = "http://localhost:8086/"
	// DbName is influxdb name.
	DbName           string = "sensorCLI"
	username         string = "todor"
	influxDbpassword string = "Abcd1234!"
)

var InfluxdbClient influxClient.Client

//TODO property file should be implemented
const (
	token  = `EOrpjNFDzHgKh-9Mpimun7SaWohfSUTXRbyJGTCQCrLM8vXvR10QXMvi4VPg8JgnVp6nIyC2VVK82PMAW08EkQ==`
	Org    = "my-org"
	Bucket = "my-bucket"
)

func init() {
	influxLogger.Info("Initializing influx db 2.0 client")
	dbClient := influxClient.NewClient(address, token)
	InfluxdbClient = dbClient
}
