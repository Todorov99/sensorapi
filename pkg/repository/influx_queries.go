package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Todorov99/server/pkg/database/influx"
	"github.com/Todorov99/server/pkg/models"
	"github.com/Todorov99/server/pkg/repository/query"

	//influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	//"github.com/influxdata/influxdb/client/v2"
)

func createPoint(data models.Measurement) *write.Point {

	_, err := time.Parse(time.RFC3339, data.MeasuredAt)
	if err != nil {
		repositoryLogger.Panic(err)
	}

	value, _ := strconv.ParseFloat(data.Value, 64)

	tags := map[string]string{"deviceID": data.DeviceID, "sensorID": data.SensorID}
	fields := map[string]interface{}{
		"value": value,
	}

	point := influxdb2.NewPoint("sensor", tags, fields, time.Now())

	return point
}

func writePointToBatch(measurementData models.Measurement) {
	defer func() {
		influx.InfluxdbClient.Close()
	}()

	writeAPI := influx.InfluxdbClient.WriteAPIBlocking("my-org", "my-bucket")

	err := writeAPI.WritePoint(context.Background(), createPoint(measurementData))
	if err != nil {
		repositoryLogger.Error(err)
		repositoryLogger.Panic(err)
	}

}

func executeSelectQueryInflux(querry string, args ...interface{}) ([]interface{}, error) {
	var measurement []interface{}

	queryAPI := influx.InfluxdbClient.QueryAPI("my-org")

	startTimestamp := args[0]
	endTimestamp := args[1]
	deviceID := args[2]
	sensorID := args[3]

	influxQuery := fmt.Sprintf(query.GetSensorAndDeviceBeetweenTimestampQuery, startTimestamp, endTimestamp, deviceID, sensorID)

	queryResult, err := queryAPI.Query(context.Background(), influxQuery)
	if err != nil {
		return nil, fmt.Errorf("failed executing query: %q, err: %w", influxQuery, err)
	}

	for queryResult.Next() {
		measurement = append(measurement, models.Measurement{
			MeasuredAt: queryResult.Record().Time().String(),
			Value:      strconv.FormatFloat(queryResult.Record().ValueByKey("_value").(float64), 'f', -1, 64),
			SensorID:   queryResult.Record().ValueByKey("sensorID").(string),
			DeviceID:   queryResult.Record().ValueByKey("deviceID").(string),
		})
	}

	if queryResult.Err() != nil {
		repositoryLogger.Errorf("query error: %w", queryResult.Err())
		return measurement, queryResult.Err()
	}

	return measurement, nil
}
