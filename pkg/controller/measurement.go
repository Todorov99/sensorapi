package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Todorov99/server/pkg/models"
	"github.com/Todorov99/server/pkg/repository"
)

var measurementRepository = repository.CreateMeasurementRepository()

var (
	measurements                = models.Measurement{}
	measurementsBetweeTimestamp = models.MeasurementBetweenTimestamp{}
)

type measurementController struct{}

func createMeasurementController() IController {
	return &measurementController{}
}

//Get gets measurement for current device and sensor ID between concrete timestamp
func (s *measurementController) Get(w http.ResponseWriter, r *http.Request) {
	controllerLogger.Info("Measurement GET query execution.")
	err := json.NewDecoder(r.Body).Decode(&measurementsBetweeTimestamp)
	if err != nil {
		respond(w, "", "Measurement Get query", err, measurementsBetweeTimestamp, 500)
		return
	}

	timestampMeasurements, err := measurementRepository.GetByID(measurementsBetweeTimestamp.StartTime, measurementsBetweeTimestamp.EndTime, measurementsBetweeTimestamp.DeviceID, measurementsBetweeTimestamp.SensorID)
	respond(w, "", "Measurement GET query execution.", err, timestampMeasurements, 404)
}

//Post executes post request to influx 2.0 db
func (s *measurementController) Post(w http.ResponseWriter, r *http.Request) {
	controllerLogger.Info("Measurement POST query execution.")

	err := json.NewDecoder(r.Body).Decode(&measurements)
	if err != nil {
		respond(w, "", "Measurement Get query", err, measurements, 500)
		return
	}

	err = measurementRepository.Add(measurements.MeasuredAt, measurements.Value,
		measurements.SensorID, measurements.DeviceID)
	if err != nil {
		respond(w, "Measurement added.", "Measurement POST query execution.", err, measurements, http.StatusConflict)
	}

	respond(w, "Measurement added.", "Measurement POST query execution.", err, measurements, http.StatusOK)
}

func (s *measurementController) Put(w http.ResponseWriter, r *http.Request) {
	err := measurementRepository.Update()
	respond(w, "", "Measurement PUT query execution.", err, measurements, 501)
}

func (s *measurementController) Delete(w http.ResponseWriter, r *http.Request) {
	measurements, err := measurementRepository.Delete("")
	respond(w, "", "Measurement DELETE query execution.", err, measurements, 501)
}

func getSensorAverageValue(w http.ResponseWriter, r *http.Request) {
	var averageValue = make(map[string]string)

	urlParams := getURLQueryParams(r, "deviceId", "sensorId", "startTime", "endTime")
	fmt.Println(urlParams)
	value, err := repository.GetAverageValueOfMeasurements(urlParams[0], urlParams[1], urlParams[2], urlParams[3])
	if err != nil {
		respond(w, "", "Failed getting average value", err, 0, http.StatusNotFound)
		return
	}

	averageValue["averageValue"] = value

	respond(w, "", "Getting sensor average values.", err, averageValue, http.StatusNotFound)
}

func getSensorsCorrelationCoefficient(w http.ResponseWriter, r *http.Request) {

	var correlationCoefficient = make(map[string]float64)

	urlParams := getURLQueryParams(r, "deviceId1", "deviceId2", "sensorId1", "sensorId2", "startTime", "endTime")
	fmt.Println(urlParams)
	value, err := repository.GetSensorsCorrelationCoefficient(urlParams[0], urlParams[1], urlParams[2], urlParams[3], urlParams[4], urlParams[5])

	correlationCoefficient["correlationCoefficient"] = value
	respond(w, "", "Getting Correlation Coefficient.", err, correlationCoefficient, http.StatusNotFound)
}
