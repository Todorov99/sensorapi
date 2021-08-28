package controller

import (
	"fmt"
	"net/http"

	"github.com/Todorov99/sensorcli/pkg/logger"
	"github.com/gorilla/mux"
)

var controllerLogger = logger.NewLogger("controller")

// HandleRequest http requests
func HandleRequest() {

	routes := mux.NewRouter().StrictSlash(true)

	routes.HandleFunc("/test", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(rw, "Testing")
	})
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
	routes.HandleFunc("/measurement", createMeasurementController().Post).Methods("POST")
	routes.HandleFunc("/measurement", createMeasurementController().Put).Methods("PUT")
	routes.HandleFunc("/measurement", createMeasurementController().Delete).Methods("DELETE")
	routes.HandleFunc("/sensorAverageValue", getSensorAverageValue).Methods("GET")
	routes.HandleFunc("/sensorsCorrelationCoefficient", getSensorsCorrelationCoefficient).Methods("GET")

	controllerLogger.Panic(http.ListenAndServe(":8081", routes))
}
