package controller

import (
	"encoding/json"
	"net/http"

	"github.com/Todorov99/server/pkg/models"
	"github.com/Todorov99/server/pkg/service"
)

var sensor = models.Sensor{}

type sensorController struct {
	sensorService service.IService
}

func NewSensorController() IController {
	return &sensorController{
		sensorService: service.NewSensorService(),
	}
}

func (s *sensorController) GetAll(w http.ResponseWriter, r *http.Request) {
	data, err := s.sensorService.GetAll()
	response(w, "", "Sensor GET query execution.", err, data, http.StatusNotFound)
}

func (s *sensorController) Get(w http.ResponseWriter, r *http.Request) {
	sensorID := getIDFromPathVariable(r)
	controllerLogger.Debugf("DEBUG path variable: %s", sensorID)

	data, err := s.sensorService.GetById(sensorID)
	response(w, "", "Sensor GET query execution.", err, data, http.StatusNotImplemented)
}

func (s *sensorController) Post(w http.ResponseWriter, r *http.Request) {
	err := json.NewDecoder(r.Body).Decode(&sensor)
	if err != nil {
		response(w, "", "Sensor Post query", err, sensor, 500)
		return
	}

	err = s.sensorService.Add(sensor)
	response(w, "You successfully add your sensor: ", "Sensor POST query execution.", err, sensor, http.StatusConflict)
}

func (s *sensorController) Put(w http.ResponseWriter, r *http.Request) {
	err := json.NewDecoder(r.Body).Decode(&sensor)
	if err != nil {
		response(w, "", "Sensor Put query", err, sensor, 500)
		return
	}

	sensorID := getIDFromPathVariable(r)
	sensor.ID = sensorID
	err = s.sensorService.Update(sensor)
	response(w, "You successfully uppdate sensor device: ", "Sensor PUT query execution.", err, sensor, http.StatusConflict)
}

func (s *sensorController) Delete(w http.ResponseWriter, r *http.Request) {
	sensorID := getIDFromPathVariable(r)
	sensor, err := s.sensorService.Delete(sensorID)
	response(w, "You successfully delete your sensor: ", "Sensor DELETE query execution.", err, sensor, http.StatusConflict)
}
