package controller

import (
	"encoding/json"
	"net/http"

	"github.com/Todorov99/server/pkg/models"
	"github.com/Todorov99/server/pkg/repository"
)

var measurementRepository = repository.CreateMeasurementRepository()

var measurements = models.Measurement{}

type measurementController struct{}

func createMeasurementController() IController {
	return &measurementController{}
}

func (s *measurementController) Get(w http.ResponseWriter, r *http.Request) {

	err := json.NewDecoder(r.Body).Decode(&measurements)

	if err != nil {
		respond(w, "", "Measurement Get query", err, measurements, 500)
		return
	}

	measurements, err := measurementRepository.GetByID(measurements.DeviceID, measurements.SensorID)
	respond(w, "", "Measurement GET query execution.", err, measurements, 404)
}

func (s *measurementController) Post(w http.ResponseWriter, r *http.Request) {

	//controllerLogger.Info().Println("Measurement POST query execution.")

	err := json.NewDecoder(r.Body).Decode(&measurements)

	if err != nil {
		respond(w, "", "Measurement Get query", err, measurements, 500)
		return
	}

	err = measurementRepository.Add(measurements.MeasuredAt, measurements.Value,
		measurements.SensorID, measurements.DeviceID)

	errStatusCode := http.StatusNotFound

	if err.Error() == "Invalid timestamp" {
		errStatusCode = http.StatusConflict
	}

	respond(w, "Measurement added.", "Measurement POST query execution.", err, measurements, errStatusCode)
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
	value, err := repository.GetAverageValueOfMeasurements(urlParams[0], urlParams[1], urlParams[2]+"+03:00", urlParams[3]+"+03:00")
	averageValue["averageValue"] = value

	respond(w, "", "Getting sensor average values.", err, averageValue, 404)

}

func getSensorsCorrelationCoefficient(w http.ResponseWriter, r *http.Request) {

	var correlationCoefficient = make(map[string]float64)

	urlParams := getURLQueryParams(r, "deviceId1", "deviceId2", "sensorId1", "sensorId2", "startTime", "endTime")
	value, err := repository.GetSensorsCorrelationCoefficient(urlParams[0], urlParams[1], urlParams[2], urlParams[3], urlParams[4]+"+03:00", urlParams[5]+"+03:00")

	correlationCoefficient["correlationCoefficient"] = value
	respond(w, "", "Getting Correlation Coefficient.", err, correlationCoefficient, 404)

}

func getURLQueryParams(r *http.Request, params ...string) []string {

	keys := r.URL.Query()

	if len(params) == 6 {
		return []string{keys.Get(params[0]), keys.Get(params[1]), keys.Get(params[2]), keys.Get(params[3]),
			keys.Get(params[4]), keys.Get(params[5])}
	}

	return []string{keys.Get(params[0]), keys.Get(params[1]), keys.Get(params[2]), keys.Get(params[3])}
}
