package repository

import (
	"fmt"
	"math"
	"time"

	"github.com/Todorov99/server/pkg/models"
	"github.com/Todorov99/server/pkg/repository/query"
)

type measurementRepository struct{}

func (m *measurementRepository) GetAll() (interface{}, error) {
	return nil, nil
}

// GetByID gets measurements for current sensor ID and device ID
// between current timestamp
func (m *measurementRepository) GetByID(args ...string) (interface{}, error) {
	repositoryLogger.Infof("Getting measurements between %s - %s for device ID: %s and sensor ID: %s")
	err := checkForExistingDevicesAndSensors(args[2], args[3])
	if err != nil {
		msg := "failed checking existing device %s and sensor %s"
		repositoryLogger.Errorf(msg, args[2], args[3])
		return nil, fmt.Errorf(msg, args[2], args[3])
	}

	startTimestamp := args[0]
	endTimestamp := args[1]
	deviceID := args[2]
	sensorID := args[3]
	influxQuery := fmt.Sprintf(query.GetMeasurementsBeetweenTimestampByDeviceIdAndSensorId, startTimestamp, endTimestamp, deviceID, sensorID)

	measurements, err := executeSelectQueryInflux(influxQuery, true)
	if err != nil {
		return 0, err
	}

	return measurements, nil
}

//Add adds measurement into influx 2.0 db
func (m *measurementRepository) Add(args ...string) error {
	err := checkForExistingDevicesAndSensors(args[3], args[2])
	if err != nil {
		return err
	}

	_, err = time.Parse(time.RFC3339, args[0])
	if err != nil {
		return fmt.Errorf("invalid timestamp")
	}

	addMeasurementBindingModel := models.Measurement{
		MeasuredAt: args[0],
		Value:      args[1],
		SensorID:   args[2],
		DeviceID:   args[3],
	}

	writePointToBatch(addMeasurementBindingModel)
	return nil
}

func (m *measurementRepository) Update(args ...string) error {
	return nil
}

func (m *measurementRepository) Delete(name string) (interface{}, error) {
	return nil, nil
}

// GetAverageValueOfMeasurements gets average values between two timestamps.
func GetAverageValueOfMeasurements(deviceID string, sensorID string, startTime string, endTime string) (string, error) {
	repositoryLogger.Infof("Getting average value of measurements between %s - %s", startTime, endTime)
	err := checkForExistingDevicesAndSensors(deviceID, sensorID)
	if err != nil {
		return "", err
	}

	influxQuery := fmt.Sprintf(query.GetAverageValueOfMeasurementsBetweenTimeStampByDeviceIdAndSensorId, startTime, endTime, deviceID, sensorID)

	response, err := executeSelectQueryInflux(influxQuery, true)
	if err != nil {
		return "", err
	}

	return response[0].(models.Measurement).Value, nil
}

// GetSensorsCorrelationCoefficient gets Pearson's correlation coefficient between two sensors.
func GetSensorsCorrelationCoefficient(deviceID1 string, deviceID2 string, sensorID1 string, sensorID2 string, startTime string, endTime string) (float64, error) {
	repositoryLogger.Info("Getting correlation coficient...")
	err := checkForExistingDevicesAndSensors(deviceID1, sensorID1)
	if err != nil {
		return 0, err
	}

	repositoryLogger.Infof("Getting values for deviceID: %s and sensorID %s...", deviceID1, sensorID1)
	firstSensorValues, err := executeSelectQueryInflux(fmt.Sprintf(query.GetMeasurementValuesByDeviceAndSensorIdBeetweenTimestamp, startTime, endTime, deviceID1, sensorID1), false)
	if err != nil {
		return 0, err
	}

	repositoryLogger.Infof("Getting values for deviceID: %s and sensorID %s...", deviceID2, sensorID2)
	secondSensorValues, err := executeSelectQueryInflux(fmt.Sprintf(query.GetMeasurementValuesByDeviceAndSensorIdBeetweenTimestamp, startTime, endTime, deviceID2, sensorID2), false)
	if err != nil {
		return 0, err
	}

	repositoryLogger.Info("Getting the count of values...")
	valueCount, err := executeSelectQueryInflux(fmt.Sprintf(query.CountMeasurementValues, startTime, endTime, deviceID1, sensorID1), false)
	if err != nil {
		return 0, err
	}

	return correlationCoefficient(firstSensorValues, secondSensorValues, parseFloat(valueCount[0])), nil
}

func correlationCoefficient(firstSensorValues []interface{}, secondSensorValues []interface{}, valueCount float64) float64 {

	sumFirstSensor := 0.0
	sumSecondSensor := 0.0
	sumBothSensorValues := 0.0
	squareSumFirstSensor := 0.0
	squareSumSecondSensor := 0.0

	for i := 0; i < int(valueCount)-1; i++ {

		if i == len(firstSensorValues) || i == len(secondSensorValues) {
			break
		}

		sumFirstSensor = sumFirstSensor + firstSensorValues[i].(float64)

		sumSecondSensor = sumSecondSensor + secondSensorValues[i].(float64)

		sumBothSensorValues = sumBothSensorValues + firstSensorValues[i].(float64)*secondSensorValues[i].(float64)

		squareSumFirstSensor = squareSumFirstSensor + firstSensorValues[i].(float64)*firstSensorValues[i].(float64)
		squareSumSecondSensor = squareSumSecondSensor + secondSensorValues[i].(float64)*secondSensorValues[i].(float64)
	}

	return float64((valueCount*sumBothSensorValues - sumFirstSensor*sumSecondSensor)) /
		(math.Sqrt(float64((valueCount*squareSumFirstSensor - sumFirstSensor*sumFirstSensor) *
			(valueCount*squareSumSecondSensor - sumSecondSensor*sumSecondSensor))))
}

func parseFloat(v interface{}) float64 {

	switch v.(type) {
	case int64:
		return float64(v.(int64))
	case int32:
		return float64(v.(int32))
	case float64:
		return float64(v.(float64))
	case float32:
		return float64(v.(float64))
	}

	return 0
}
