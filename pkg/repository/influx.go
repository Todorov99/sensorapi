package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Todorov99/sensorapi/pkg/entity"
	"github.com/Todorov99/sensorapi/pkg/global"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

func createPoint(data entity.Measurement) (*write.Point, error) {
	measurementTime, err := time.Parse(time.RFC3339, data.MeasuredAt)
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

	fmt.Println(measurementTime)
	point := influxdb2.NewPoint("sensor", tags, fields, measurementTime)
	return point, nil
}

func writePointToBatch(measurementData entity.Measurement, influxClient influxdb2.Client, org, bucket string) {
	defer func() {
		influxClient.Close()
	}()

	writeAPI := influxClient.WriteAPIBlocking(org, bucket)

	influxDbPoint, err := createPoint(measurementData)
	if err != nil {
		repositoryLogger.Panic(fmt.Errorf("failed creating a influx DB point: %w", err))
	}

	err = writeAPI.WritePoint(context.Background(), influxDbPoint)
	if err != nil {
		repositoryLogger.Panic(err)
	}

}

func executeSelectQueryInflux(ctx context.Context, querry string, isType bool, influxClient influxdb2.Client, org, bucket string) ([]interface{}, error) {
	var measurement []interface{}

	queryAPI := influxClient.QueryAPI(org)

	queryResult, err := queryAPI.Query(ctx, querry)
	if err != nil {
		return nil, fmt.Errorf("failed executing query: %s, err: %w", querry, err)
	}

	for queryResult.Next() {
		if !isType {
			measurement = append(measurement, queryResult.Record().ValueByKey("_value"))
		} else {
			measurement = append(measurement, entity.Measurement{
				MeasuredAt: queryResult.Record().Time().Format(global.TimeFormat),
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
