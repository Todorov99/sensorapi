package controller

import (
	"encoding/json"
	"net/http"

	"github.com/Todorov99/server/pkg/models"
	"github.com/Todorov99/server/pkg/repository"
)

var sensor = models.Sensor{}

var sensorRepository = repository.CreateSensorRepository()

type sensorController struct{}

func createSensorController() IController {
	return &sensorController{}
}

func (s *sensorController) Get(w http.ResponseWriter, r *http.Request) {
	pathVariable := getIDFromPathVariable(r)

	if pathVariable != "" {
		data, err := sensorRepository.GetAll()
		respond(w, "", "Sensor GET query execution.", err, data, http.StatusNotImplemented)
		return
	}

	data, err := sensorRepository.GetByID()
	respond(w, "", "Sensor GET query execution.", err, data, http.StatusNotFound)
}

func (s *sensorController) Post(w http.ResponseWriter, r *http.Request) {

	err := json.NewDecoder(r.Body).Decode(&sensor)

	if err != nil {
		respond(w, "", "Sensor Post query", err, sensor, 500)
		return
	}

	err = sensorRepository.Add(sensor.Name, sensor.Description, sensor.DeviceId, sensor.SensorGroups, sensor.Unit)
	respond(w, "You successfully add your sensor: ", "Sensor POST query execution.", err, sensor, http.StatusConflict)

}

func (s *sensorController) Put(w http.ResponseWriter, r *http.Request) {

	err := json.NewDecoder(r.Body).Decode(&sensor)

	if err != nil {
		respond(w, "", "Sensor Put query", err, sensor, 500)
		return
	}

	sensorID := getIDFromPathVariable(r)

	err = sensorRepository.Update(sensor.Name, sensor.Description, sensor.Unit, sensor.SensorGroups, sensorID)
	respond(w, "You successfully uppdate sensor device: ", "Sensor PUT query execution.", err, sensor, http.StatusConflict)
}

func (s *sensorController) Delete(w http.ResponseWriter, r *http.Request) {

	sensorID := getIDFromPathVariable(r)

	sensor, err := sensorRepository.Delete(sensorID)
	respond(w, "You successfully delete your sensor: ", "Sensor DELETE query execution.", err, sensor, http.StatusConflict)

}
