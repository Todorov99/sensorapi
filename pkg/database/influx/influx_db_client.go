package influx

import (
	"github.com/Todorov99/sensorcli/pkg/logger"
	"github.com/influxdata/influxdb/client/v2"
)

var influxLogger = logger.NewLogger("./influx")

const (
	address string = "http://influx:8086/"
	// DbName is influxdb name.
	DbName           string = "sensorCLI"
	username         string = "todor"
	influxDbpassword string = "1234"
)

// InfluxdbClient opens influx connection.
var InfluxdbClient client.Client

func init() {

	dbClient, err := client.NewHTTPClient(client.HTTPConfig{
		Password: influxDbpassword,
		Addr:     address,
		Username: username,
	})

	if err != nil {
		influxLogger.Panic(err)
	}

	InfluxdbClient = dbClient

}
