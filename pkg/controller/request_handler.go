package controller

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Todorov99/sensorcli/pkg/logger"
	"github.com/gorilla/mux"
)

var controllerLogger = logger.NewLogrus("controller", os.Stdout)

// HandleRequest http requests
func HandleRequest(port string) error {

	fs := http.FileServer(http.Dir("/Users/t.todorov/Develop/server/resources/static"))
	routes := mux.NewRouter().StrictSlash(true)

	routes.PathPrefix("/static").Handler(http.StripPrefix("/static", fs))
	http.Handle("/", routes)

	routes.HandleFunc("/device/{id}", createDeviceController().Get).Methods("GET")
	routes.HandleFunc("/device", getAllDevices).Methods("GET")
	routes.HandleFunc("/device", createDeviceController().Post).Methods("POST")
	routes.HandleFunc("/device/{id}", createDeviceController().Put).Methods("PUT")
	routes.HandleFunc("/device/{id}", createDeviceController().Delete).Methods("DELETE")

	routes.HandleFunc("/sensor", createSensorController().Get).Methods("GET")
	routes.HandleFunc("/sensor/{id}", createSensorController().Get).Methods("GET")
	routes.HandleFunc("/sensor", createSensorController().Post).Methods("POST")
	routes.HandleFunc("/sensor/{id}", createSensorController().Put).Methods("PUT")
	routes.HandleFunc("/sensor/{id}", createSensorController().Delete).Methods("DELETE")

	routes.HandleFunc("/measurement", createMeasurementController().Get).Methods("GET")
	routes.HandleFunc("/monitor", monitor).Methods("GET")
	routes.HandleFunc("/measurement", createMeasurementController().Post).Methods("POST")
	routes.HandleFunc("/measurement", createMeasurementController().Put).Methods("PUT")
	routes.HandleFunc("/measurement", createMeasurementController().Delete).Methods("DELETE")
	routes.HandleFunc("/sensorAverageValue", getSensorAverageValue).Methods("GET")
	routes.HandleFunc("/sensorsCorrelationCoefficient", getSensorsCorrelationCoefficient).Methods("GET")

	return http.ListenAndServe(fmt.Sprintf(":%s", port), routes)
}
