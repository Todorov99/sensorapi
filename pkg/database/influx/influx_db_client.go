package influx

import (
	"github.com/Todorov99/sensorcli/pkg/logger"
	influxClient "github.com/influxdata/influxdb-client-go/v2"
)

var influxLogger = logger.NewLogger("./influx")

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
	bucket = "my-bucket"
)

func init() {
	influxLogger.Info("Initializing influx db 2.0 client")
	dbClient := influxClient.NewClient(address, token)
	InfluxdbClient = dbClient
}
