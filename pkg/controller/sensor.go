package controller

import (
	"encoding/json"
	"net/http"

	"github.com/Todorov99/sensorapi/pkg/dto"
	"github.com/Todorov99/sensorapi/pkg/service"
	"github.com/dgrijalva/jwt-go"
)

var sensor = dto.Sensor{}

type sensorController struct {
	sensorService service.SensorService
}

func NewSensorController() IController {
	return &sensorController{
		sensorService: service.NewSensorService(),
	}
}

func (s *sensorController) GetAll(w http.ResponseWriter, r *http.Request, token *jwt.Token) {
	defer func() {
		r.Body.Close()
	}()

	data, err := s.sensorService.GetAll(r.Context())
	response(w, "Sensor GET query execution.", err, data, http.StatusNotFound)
}

func (s *sensorController) GetByID(w http.ResponseWriter, r *http.Request, token *jwt.Token) {
	defer func() {
		r.Body.Close()
	}()

	sensorID := getIDFromPathVariable(r)
	controllerLogger.Debugf("DEBUG path variable: %s", sensorID)

	data, err := s.sensorService.GetById(r.Context(), sensorID)
	response(w, "Sensor GET query execution.", err, data, http.StatusNotImplemented)
}

// Not implemened
func (s *sensorController) Post(w http.ResponseWriter, r *http.Request, token *jwt.Token) {
}

func (s *sensorController) Put(w http.ResponseWriter, r *http.Request, token *jwt.Token) {
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

// Not implemened
func (s *sensorController) Delete(w http.ResponseWriter, r *http.Request, token *jwt.Token) {
}
