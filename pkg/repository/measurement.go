package repository

import (
	"encoding/json"
	"errors"
	"math"
	"time"

	"github.com/Todorov99/server/pkg/models"
)

type measurementRepository struct{}

func (m *measurementRepository) GetAll() (interface{}, error) {
	return nil, errors.New("Not Implemented")
}

func (m *measurementRepository) GetByID(args ...string) (interface{}, error) {
	err := checkForExistingDevicesAndSensors(args[0], args[1])

	if err != nil {
		return nil, err
	}

	var measurements []models.Measurement
	var currentMeasurement models.Measurement

	querry := "select * from sensor where deviceID = '%s' and sensorID = '%s'"

	response, err := executeSelectQueryInflux(querry, args[0], args[1])

	if err != nil {
		return 0, err
	}

	for _, i := range response {

		currentMeasurement.MeasuredAt = i[0].(string)
		currentMeasurement.DeviceID = i[1].(string)
		currentMeasurement.SensorID = i[2].(string)
		currentMeasurement.Value = i[3].(json.Number).String()

		measurements = append(measurements, currentMeasurement)

	}

	return measurements, nil
}

func (m *measurementRepository) Add(args ...string) error {
	err := checkForExistingDevicesAndSensors(args[3], args[2])

	if err != nil {
		return err
	}

	_, err = time.Parse(time.RFC3339, args[0])

	if err != nil {
		return errors.New("Invalid timestamp")
	}

	addMeasurementBindingModel := models.Measurement{args[0], args[1], args[2], args[3]}

	writePointToBatch(addMeasurementBindingModel)
	return nil
}

func (m *measurementRepository) Update(args ...string) error {
	return errors.New("Not Implemented")
}

func (m *measurementRepository) Delete(name string) (interface{}, error) {
	return nil, errors.New("Not Implemented")
}

// GetAverageValueOfMeasurements gets average values between two timestamps.
func GetAverageValueOfMeasurements(deviceID string, sensorID string, startTime string, endTime string) (string, error) {

	err := checkForExistingDevicesAndSensors(deviceID, sensorID)

	if err != nil {
		return "", err
	}

	querry := "select MEAN(value) from sensor where time > '%s' and time < '%s' and deviceID = '%s' and sensorID='%s'"

	response, err := executeSelectQueryInflux(querry, startTime, endTime, deviceID, sensorID)

	if err != nil {
		return "", err
	}

	return response[0][1].(json.Number).String(), nil
}

// GetSensorsCorrelationCoefficient gets Pearson's correlation coefficient between two sensors.
func GetSensorsCorrelationCoefficient(deviceID1 string, deviceID2 string, sensorID1 string, sensorID2 string, startTime string, endTime string) (float64, error) {

	err := checkForExistingDevicesAndSensors(deviceID1, sensorID1)

	if err != nil {
		return 0, err
	}

	sensorValues := "select value from sensor where deviceID='%s' and sensorID='%s' and time > '%s' and time < '%s'"
	countQuery := "select count(value) from sensor where deviceID='%s' and sensorID='%s' and time > '%s' and time < '%s'"

	firstSensorValues, err := executeSelectQueryInflux(sensorValues, deviceID1, sensorID1, startTime, endTime)

	if err != nil {
		return 0, err
	}

	secondSensorValues, err := executeSelectQueryInflux(sensorValues, deviceID2, sensorID2, startTime, endTime)

	if err != nil {
		return 0, err
	}

	valueCount, err := executeSelectQueryInflux(countQuery, deviceID1, sensorID1, startTime, endTime)

	if err != nil {
		return 0, err
	}

	sensorValuesCount, _ := valueCount[0][1].(json.Number).Float64()

	firstSensorValuesSum := appendDataToSlice(firstSensorValues)
	secondSensorValuesSum := appendDataToSlice(secondSensorValues)

	coef := correlationCoefficient(firstSensorValuesSum, secondSensorValuesSum, sensorValuesCount)

	return coef, nil
}

func correlationCoefficient(firstSensorValues []float64, secondSensorValues []float64, valueCount float64) float64 {

	sumFirstSensor := 0.0
	sumSecondSensor := 0.0
	sumBothSensorValues := 0.0
	squareSumFirstSensor := 0.0
	squareSumSecondSensor := 0.0

	for i := 0; i < int(valueCount)-1; i++ {

		if i == len(firstSensorValues) || i == len(secondSensorValues) {
			break
		}

		sumFirstSensor = sumFirstSensor + firstSensorValues[i]

		sumSecondSensor = sumSecondSensor + secondSensorValues[i]

		sumBothSensorValues = sumBothSensorValues + firstSensorValues[i]*secondSensorValues[i]

		squareSumFirstSensor = squareSumFirstSensor + firstSensorValues[i]*firstSensorValues[i]
		squareSumSecondSensor = squareSumSecondSensor + secondSensorValues[i]*secondSensorValues[i]
	}

	return float64((valueCount*sumBothSensorValues - sumFirstSensor*sumSecondSensor)) /
		(math.Sqrt(float64((valueCount*squareSumFirstSensor - sumFirstSensor*sumFirstSensor) *
			(valueCount*squareSumSecondSensor - sumSecondSensor*sumSecondSensor))))
}

func appendDataToSlice(data [][]interface{}) []float64 {

	var dataSlice []float64

	for _, i := range data {
		num, _ := i[1].(json.Number).Float64()
		dataSlice = append(dataSlice, num)
	}

	return dataSlice
}
