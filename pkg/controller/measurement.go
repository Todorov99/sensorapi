package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Todorov99/sensorapi/pkg/dto"
	"github.com/Todorov99/sensorapi/pkg/server/config"
	"github.com/Todorov99/sensorapi/pkg/service"
	"github.com/dgrijalva/jwt-go"
)

var (
	measurements                = dto.Measurement{}
	measurementsBetweeTimestamp = dto.MeasurementBetweenTimestamp{}
)

type measurementController struct {
	measurementService service.MeasurementService
}

func NewMeasurementController() *measurementController {
	return &measurementController{
		measurementService: service.NewMeasurementService(),
	}
}

func (m *measurementController) GetAllMeasurementsForSensorAndDeviceIDBetweenTimestamp(w http.ResponseWriter, r *http.Request, token *jwt.Token) {
	controllerLogger.Info("Measurement GET query execution.")
	defer func() {
		r.Body.Close()
	}()

	err := json.NewDecoder(r.Body).Decode(&measurementsBetweeTimestamp)
	if err != nil {
		response(w, "Measurement Get query", err, measurementsBetweeTimestamp, 500)
		return
	}

	timestampMeasurements, err := m.measurementService.GetMeasurementsBetweenTimestamp(r.Context(), measurementsBetweeTimestamp)
	if err != nil {
		response(w, "Get all measurements between timestamp finished with error", err, nil, http.StatusBadRequest)
		return
	}
	response(w, "Measurement GET query execution.", err, timestampMeasurements, http.StatusOK)
}

func (s *measurementController) AddMeasurement(w http.ResponseWriter, r *http.Request, token *jwt.Token) {
	defer func() {
		r.Body.Close()
	}()

	controllerLogger.Info("Measurement POST query execution.")

	err := json.NewDecoder(r.Body).Decode(&measurements)
	if err != nil {
		response(w, "Measurement Get query", err, measurements, 500)
		return
	}

	err = s.measurementService.AddMeasurements(r.Context(), measurements)
	if err != nil {
		response(w, "Measurement POST query execution.", err, measurements, http.StatusConflict)
	}

	response(w, "Measurement POST query execution.", err, measurements, http.StatusOK)
}

func (m *measurementController) GetSensorAverageValue(w http.ResponseWriter, r *http.Request, token *jwt.Token) {
	defer func() {
		r.Body.Close()
	}()

	var averageValue = make(map[string]string)

	keys := r.URL.Query()
	deviceID := keys["deviceId"][0]
	sensorID := keys["sensorId"][0]
	startTime := keys["startTime"][0]
	endTime := keys["endTime"][0]

	averageVal, err := m.measurementService.GetAverageValueOfMeasurements(r.Context(), deviceID, sensorID, startTime, endTime)
	if err != nil {
		response(w, "Failed getting average value", err, 0, http.StatusNotFound)
		return
	}

	averageValue["averageValue"] = averageVal
	response(w, "Getting sensor average values.", err, averageValue, http.StatusNotFound)
}

func (m *measurementController) GetSensorsCorrelationCoefficient(w http.ResponseWriter, r *http.Request, token *jwt.Token) {
	defer func() {
		r.Body.Close()
	}()

	var correlationCoefficient = make(map[string]float64)

	keys := r.URL.Query()
	deviceId1 := keys["deviceId1"][0]
	deviceId2 := keys["deviceId2"][0]
	sensorId1 := keys["sensorId1"][0]
	sensorId2 := keys["sensorId2"][0]
	startTime := keys["startTime"][0]
	endTime := keys["endTime"][0]

	coefficient, err := m.measurementService.GetSensorsCorrelationCoefficient(r.Context(), deviceId1, deviceId2, sensorId1, sensorId2, startTime, endTime)
	correlationCoefficient["correlationCoefficient"] = coefficient
	response(w, "Getting Correlation Coefficient.", err, correlationCoefficient, http.StatusNotFound)
}

func (m *measurementController) Monitor(w http.ResponseWriter, r *http.Request, token *jwt.Token) {
	defer func() {
		r.Body.Close()
	}()

	monitorDto := dto.MonitorDto{}
	decodeErr := json.NewDecoder(r.Body).Decode(&monitorDto)
	if decodeErr != nil {
		controllerLogger.Error(decodeErr)
		return
	}

	if m.measurementService.GetMonitorStatus(monitorDto.DeviceID).Status == service.StateInProgress {
		response(w, "There is running measurement", fmt.Errorf("running process"), nil, http.StatusBadRequest)
		return
	}

	done, err := m.measurementService.Monitor(r.Context(), config.GetJWTEmailClaim(token), config.GetJWTUserIDClaim(token), monitorDto)
	if err != nil {
		response(w, "Starting the monitoring process failed", err, nil, http.StatusBadRequest)
		return
	}

	for {
		select {
		case <-done:
			response(w, "Monitoring finished", nil, "Monitoring finished successfully", http.StatusOK)
			return
		case <-r.Context().Done():
			response(w, "Monitoring finished with error", r.Context().Err(), nil, http.StatusConflict)
			return
		}
	}
}

func (m *measurementController) MonitorStatus(w http.ResponseWriter, r *http.Request, token *jwt.Token) {
	deviceID := getIDFromPathVariable(r)
	monitorStatus := m.measurementService.GetMonitorStatus(deviceID)
	response(w, "Getting monitor status...", nil, monitorStatus, http.StatusAccepted)
}

func (m *measurementController) GetReportFile(w http.ResponseWriter, r *http.Request, token *jwt.Token) {
	keys := r.URL.Query()
	filename := keys.Get("reportFilename")
	serverFile(w, r, filename)
}
