package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Todorov99/server/pkg/controller"
	"github.com/Todorov99/server/pkg/dto"
	"github.com/Todorov99/server/pkg/server/config"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

// HandleRequest handles the supported REST request of the Web Server
func HandleRequest(port string) error {
	routes := mux.NewRouter().StrictSlash(true)

	measurementController := controller.NewMeasurementController()
	measurementAnalizer := controller.NewMeasurementAnalizer()
	deviceController := controller.NewDeviceController()
	sensorController := controller.NewSensorController()
	userController := controller.NewUserController()

	routes.HandleFunc("/api/users/login", userController.Get).Methods("Get")
	routes.HandleFunc("/api/users/register", userController.Post).Methods("POST")

	routes.Handle("/device/{id}", isAuthorized(deviceController.Get)).Methods("GET")
	routes.Handle("/device", isAuthorized(deviceController.GetAll)).Methods("GET")
	routes.Handle("/device", isAuthorized(deviceController.Post)).Methods("POST")
	routes.Handle("/device/{id}", isAuthorized(deviceController.Put)).Methods("PUT")
	routes.Handle("/device/{id}", isAuthorized(deviceController.Delete)).Methods("DELETE")

	routes.Handle("/sensor", isAuthorized(sensorController.GetAll)).Methods("GET")
	routes.Handle("/sensor/{id}", isAuthorized(sensorController.Get)).Methods("GET")
	routes.Handle("/sensor", isAuthorized(sensorController.Post)).Methods("POST")
	routes.Handle("/sensor/{id}", isAuthorized(sensorController.Put)).Methods("PUT")
	routes.Handle("/sensor/{id}", isAuthorized(sensorController.Delete)).Methods("DELETE")

	routes.Handle("/measurement", isAuthorized(measurementController.GetAll)).Methods("GET")
	routes.Handle("/measurement", isAuthorized(measurementController.Post)).Methods("POST")

	routes.Handle("/collectMeasurements", isAuthorized(measurementAnalizer.Monitor)).Methods("POST")
	routes.Handle("/sensorAverageValue", isAuthorized(measurementAnalizer.GetSensorAverageValue)).Methods("GET")
	routes.Handle("/sensorsCorrelationCoefficient", isAuthorized(measurementAnalizer.GetSensorsCorrelationCoefficient)).Methods("GET")

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
