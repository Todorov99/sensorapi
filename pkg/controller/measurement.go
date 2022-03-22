package controller

import (
	"encoding/json"
	"net/http"

	"github.com/Todorov99/server/pkg/dto"
	"github.com/Todorov99/server/pkg/global"
	"github.com/Todorov99/server/pkg/service"
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

func (m *measurementController) GetAllMeasurementsForSensorAndDeviceIDBetweenTimestamp(w http.ResponseWriter, r *http.Request) {
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

func (s *measurementController) AddMeasurement(w http.ResponseWriter, r *http.Request) {
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

func (m *measurementController) GetSensorAverageValue(w http.ResponseWriter, r *http.Request) {
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

func (m *measurementController) GetSensorsCorrelationCoefficient(w http.ResponseWriter, r *http.Request) {
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

func (m *measurementController) Monitor(w http.ResponseWriter, r *http.Request) {
	defer func() {
		r.Body.Close()
	}()

	keys := r.URL.Query()

	valueCfg := dto.ValueCfg{}
	decodeErr := json.NewDecoder(r.Body).Decode(&valueCfg)
	if decodeErr != nil {
		controllerLogger.Error(decodeErr)
		return
	}

	metricChan := make(chan interface{})
	err := make(chan error)
	done := make(chan bool)

	sensorGroupsWithSysFiles := map[string]string{
		global.CpuTempGroup:  keys.Get("tempSysFile"),
		global.CpuUsageGroup: "",
		global.MemoryGroup:   "",
	}

	go m.measurementService.Monitor(r.Context(), keys.Get("duration"), sensorGroupsWithSysFiles, valueCfg, err, metricChan, done)

	for {
		select {
		case err := <-err:
			if err != nil {
				metric := <-metricChan
				response(w, "Monitoring finished", err, metric, http.StatusOK)
				return
			}
		case <-done:
			response(w, "Monitoring finished", nil, nil, http.StatusOK)
			return
		case rs := <-metricChan:
			response(w, "Monitoring...", nil, rs, http.StatusOK)
		case <-r.Context().Done():
			response(w, "Monitoring finished with error", r.Context().Err(), nil, http.StatusConflict)
		}
	}
}
