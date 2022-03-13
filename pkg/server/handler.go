package server

import (
	"fmt"
	"net/http"

	"github.com/Todorov99/server/pkg/controller"
	"github.com/gorilla/mux"
)

// HandleRequest handles the supported REST request of the Web Server
func HandleRequest(port string) error {
	routes := mux.NewRouter().StrictSlash(true)

	measurementController := controller.NewMeasurementController()
	measurementAnalizer := controller.NewMeasurementAnalizer()
	deviceController := controller.NewDeviceController()
	sensorController := controller.NewSensorController()

	routes.HandleFunc("/device/{id}", deviceController.Get).Methods("GET")
	routes.HandleFunc("/device", deviceController.GetAll).Methods("GET")
	routes.HandleFunc("/device", deviceController.Post).Methods("POST")
	routes.HandleFunc("/device/{id}", deviceController.Put).Methods("PUT")
	routes.HandleFunc("/device/{id}", deviceController.Delete).Methods("DELETE")

	routes.HandleFunc("/sensor", sensorController.GetAll).Methods("GET")
	routes.HandleFunc("/sensor/{id}", sensorController.Get).Methods("GET")
	routes.HandleFunc("/sensor", sensorController.Post).Methods("POST")
	routes.HandleFunc("/sensor/{id}", sensorController.Put).Methods("PUT")
	routes.HandleFunc("/sensor/{id}", sensorController.Delete).Methods("DELETE")

	routes.HandleFunc("/measurement", measurementController.GetAll).Methods("GET")
	routes.HandleFunc("/measurement", measurementController.Post).Methods("POST")

	routes.HandleFunc("/collectMeasurements", measurementAnalizer.Monitor).Methods("POST")
	routes.HandleFunc("/sensorAverageValue", measurementAnalizer.GetSensorAverageValue).Methods("GET")
	routes.HandleFunc("/sensorsCorrelationCoefficient", measurementAnalizer.GetSensorsCorrelationCoefficient).Methods("GET")

	return http.ListenAndServe(fmt.Sprintf(":%s", port), routes)
}
