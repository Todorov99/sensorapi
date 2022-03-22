package controller

import (
	"encoding/json"
	"net/http"

	"github.com/Todorov99/server/pkg/dto"
	"github.com/Todorov99/server/pkg/service"
)

var sensor = dto.Sensor{}

type sensorController struct {
	sensorService service.IService
}

func NewSensorController() IController {
	return &sensorController{
		sensorService: service.NewSensorService(),
	}
}

func (s *sensorController) GetAll(w http.ResponseWriter, r *http.Request) {
	defer func() {
		r.Body.Close()
	}()

	data, err := s.sensorService.GetAll(r.Context())
	response(w, "Sensor GET query execution.", err, data, http.StatusNotFound)
}

func (s *sensorController) GetByID(w http.ResponseWriter, r *http.Request) {
	defer func() {
		r.Body.Close()
	}()

	sensorID := getIDFromPathVariable(r)
	controllerLogger.Debugf("DEBUG path variable: %s", sensorID)

	data, err := s.sensorService.GetById(r.Context(), sensorID)
	response(w, "Sensor GET query execution.", err, data, http.StatusNotImplemented)
}

func (s *sensorController) Post(w http.ResponseWriter, r *http.Request) {
	defer func() {
		r.Body.Close()
	}()

	err := json.NewDecoder(r.Body).Decode(&sensor)
	if err != nil {
		response(w, "Sensor Post query", err, sensor, 500)
		return
	}

	err = s.sensorService.Add(r.Context(), sensor)
	response(w, "Sensor POST query execution.", err, sensor, http.StatusConflict)
}

func (s *sensorController) Put(w http.ResponseWriter, r *http.Request) {
	defer func() {
		r.Body.Close()
	}()

	err := json.NewDecoder(r.Body).Decode(&sensor)
	if err != nil {
		response(w, "Sensor Put query", err, sensor, 500)
		return
	}

	sensorID := getIDFromPathVariable(r)
	sensor.ID = int32(sensorID)
	err = s.sensorService.Update(r.Context(), sensor)
	response(w, "Sensor PUT query execution.", err, sensor, http.StatusConflict)
}

func (s *sensorController) Delete(w http.ResponseWriter, r *http.Request) {
	sensorID := getIDFromPathVariable(r)
	sensor, err := s.sensorService.Delete(r.Context(), sensorID)
	response(w, "Sensor DELETE query execution.", err, sensor, http.StatusConflict)
}
