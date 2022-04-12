package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Todorov99/serverapi/pkg/controller"
	"github.com/Todorov99/serverapi/pkg/dto"
	"github.com/Todorov99/serverapi/pkg/server/config"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

// HandleRequest handles the supported REST request of the Web Server
func HandleRequest(port string) error {
	routes := mux.NewRouter().StrictSlash(true)

	measurementController := controller.NewMeasurementController()
	deviceController := controller.NewDeviceController()
	sensorController := controller.NewSensorController()
	userController := controller.NewUserController()

	routes.HandleFunc("/api/users/login", userController.Login).Methods("Get")
	routes.HandleFunc("/api/users/register", userController.Register).Methods("POST")

	routes.Handle("/api/device/{id}", isAuthorized(deviceController.GetByID)).Methods("GET")
	routes.Handle("/api/devices/all", isAuthorized(deviceController.GetAll)).Methods("GET")
	routes.Handle("/api/device/add", isAuthorized(deviceController.Post)).Methods("POST")
	routes.Handle("/api/device/{id}", isAuthorized(deviceController.Put)).Methods("PUT")
	routes.Handle("/api/device/{id}", isAuthorized(deviceController.Delete)).Methods("DELETE")

	routes.Handle("/api/sensor/all", isAuthorized(sensorController.GetAll)).Methods("GET")
	routes.Handle("/api/sensor/{id}", isAuthorized(sensorController.GetByID)).Methods("GET")
	routes.Handle("/api/sensor/add", isAuthorized(sensorController.Post)).Methods("POST")
	routes.Handle("/api/sensor/{id}", isAuthorized(sensorController.Put)).Methods("PUT")
	routes.Handle("/api/sensor/{id}", isAuthorized(sensorController.Delete)).Methods("DELETE")

	routes.Handle("/api/measurement", isAuthorized(measurementController.GetAllMeasurementsForSensorAndDeviceIDBetweenTimestamp)).Methods("GET")
	routes.Handle("/api/measurement", isAuthorized(measurementController.AddMeasurement)).Methods("POST")

	routes.Handle("/api/measurement/collect", isAuthorized(measurementController.Monitor)).Methods("POST")
	routes.Handle("/api/measurement/average", isAuthorized(measurementController.GetSensorAverageValue)).Methods("GET")
	routes.Handle("/api/measurement/correlation", isAuthorized(measurementController.GetSensorsCorrelationCoefficient)).Methods("GET")
	routes.Handle("/api/measurement/monitor/status", isAuthorized(measurementController.MonitorStatus)).Methods("GET")
	routes.Handle("/api/measurement/monitor/report", isAuthorized(measurementController.GetReportFile)).Methods("GET")

	return http.ListenAndServe(fmt.Sprintf(":%s", port), routes)
}

func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bareerToken := r.Header.Get("Authorization")
		if bareerToken == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			responseError := dto.ResponseError{
				ErrMessage: "No Authorization toke provided",
			}
			json.NewEncoder(w).Encode(responseError)
			return
		}

		t := strings.Split(bareerToken, " ")[1]

		token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}

			jwtCfg := config.GetJWTCfg()
			tokenClaims := token.Claims.(jwt.MapClaims)
			if !tokenClaims.VerifyAudience(jwtCfg.GetJWTAudience(), false) {
				return nil, errors.New("invalid aud of the JWT")
			}

			if !tokenClaims.VerifyIssuer(jwtCfg.GetJWTIssuer(), false) {
				return nil, errors.New("invalid iss of the JWT")
			}

			return jwtCfg.GetJWTSigningKey(), nil
		})

		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			responseError := dto.ResponseError{
				ErrMessage: err.Error(),
			}
			json.NewEncoder(w).Encode(responseError)
		}

		if token.Valid {
			endpoint(w, r)
		}
	})
}
