package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Todorov99/server/pkg/database/influx"
	"github.com/Todorov99/server/pkg/models"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

func createPoint(data models.Measurement) (*write.Point, error) {
	_, err := time.Parse(time.RFC3339, data.MeasuredAt)
	if err != nil {
		return nil, err
	}

	value, err := strconv.ParseFloat(data.Value, 64)
	if err != nil {
		return nil, err
	}

	tags := map[string]string{"deviceID": data.DeviceID, "sensorID": data.SensorID}
	fields := map[string]interface{}{
		"value": value,
	}

	point := influxdb2.NewPoint("sensor", tags, fields, time.Now())

	return point, nil
}

func writePointToBatch(measurementData models.Measurement) {
	defer func() {
		influx.InfluxdbClient.Close()
	}()

	writeAPI := influx.InfluxdbClient.WriteAPIBlocking(influx.Org, influx.Bucket)

	influxDbPoint, err := createPoint(measurementData)
	if err != nil {
		repositoryLogger.Panic(fmt.Errorf("failed creating a influx DB point: %w", err))
	}

	err = writeAPI.WritePoint(context.Background(), influxDbPoint)
	if err != nil {
		repositoryLogger.Panic(err)
	}

}

func executeSelectQueryInflux(querry string, isType bool) ([]interface{}, error) {
	var measurement []interface{}

	queryAPI := influx.InfluxdbClient.QueryAPI(influx.Org)

	queryResult, err := queryAPI.Query(context.Background(), querry)
	if err != nil {
		return nil, fmt.Errorf("failed executing query: %s, err: %w", querry, err)
	}

	for queryResult.Next() {
		if !isType {
			measurement = append(measurement, queryResult.Record().ValueByKey("_value"))
		} else {
			measurement = append(measurement, models.Measurement{
				MeasuredAt: queryResult.Record().Time().String(),
				Value:      strconv.FormatFloat(queryResult.Record().ValueByKey("_value").(float64), 'f', -1, 64),
				SensorID:   queryResult.Record().ValueByKey("sensorID").(string),
				DeviceID:   queryResult.Record().ValueByKey("deviceID").(string),
			})
		}

	}

	if queryResult.Err() != nil {
		repositoryLogger.Errorf("query error: %w", queryResult.Err())
		return measurement, queryResult.Err()
	}

	return measurement, nil
}
