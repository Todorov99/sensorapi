package repository

import (
	"errors"
	"fmt"

	"strconv"
	"time"

	"github.com/Todorov99/server/pkg/database/influx"
	"github.com/Todorov99/server/pkg/models"

	"github.com/influxdata/influxdb/client/v2"
)

func createBatchPoint() client.BatchPoints {

	batchPoint, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  influx.DbName,
		Precision: "s",
	})

	if err != nil {
		//	logger.ErrorLogger.Fatalln(err)
	}

	return batchPoint
}

func createPoint(data models.Measurement) client.Point {

	time, err := time.Parse(time.RFC3339, data.MeasuredAt)
	value, _ := strconv.ParseFloat(data.Value, 64)

	tags := map[string]string{"deviceID": data.DeviceID, "sensorID": data.SensorID}
	fields := map[string]interface{}{
		"value": value,
	}

	point, err := client.NewPoint("sensor", tags, fields, time)

	if err != nil {
		//	logger.ErrorLogger.Fatalln(err)
	}

	return *point
}

func writePointToBatch(measurementData models.Measurement) {

	batchPoint := createBatchPoint()

	point := createPoint(measurementData)

	batchPoint.AddPoint(&point)

	if err := influx.InfluxdbClient.Write(batchPoint); err != nil {
		//	logger.ErrorLogger.Fatalln(err)
	}

	if err := influx.InfluxdbClient.Close(); err != nil {
		//logger.ErrorLogger.Fatalln(err)
	}

}

func executeSelectQueryInflux(querry string, args ...interface{}) ([][]interface{}, error) {

	influxQuery := client.Query{
		Command:  fmt.Sprintf(querry, args...),
		Database: influx.DbName,
	}

	response, err := influx.InfluxdbClient.Query(influxQuery)

	if err != nil {
		return nil, err
	}

	if response.Results == nil {
		return nil, errors.New("There is not available measurements")
	}

	if response.Results[0].Series == nil {
		return nil, errors.New("There is not available measurements")
	}

	return response.Results[0].Series[0].Values, nil
}
