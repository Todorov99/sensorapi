package repository

import (
	"database/sql"
	"fmt"

	"github.com/Todorov99/server/pkg/repository/query"
)

func executeSelectQuery(query string, postgreClient *sql.DB, args ...interface{}) (string, error) {

	var value string

	rowsRs, err := postgreClient.Query(query, args...)

	for rowsRs.Next() {

		err := rowsRs.Scan(&value)

		if err != nil {
			return "", err
		}
	}

	if err != nil {
		return "", err
	}

	return value, nil
}

func executeModifyingQuery(query string, postgreClient *sql.DB, args ...interface{}) error {
	_, err := postgreClient.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed executing query %q with arguments %q: %w", query, args, err)
	}

	return nil
}

func checkForExistingSensorByID(id string, postgreClient *sql.DB) bool {
	existingSensor, _ := executeSelectQuery(query.GetSensorNameByID, postgreClient, id)
	return existingSensor != ""
}

func checkForExistingDeviceByID(id string, postgreClient *sql.DB) bool {
	existingDevice, _ := executeSelectQuery(query.GetDeviceNameByID, postgreClient, id)
	return existingDevice != ""
}

func checkForExistingDevicesAndSensors(deviceID string, sensorID string, postgreClient *sql.DB) error {
	repositoryLogger.Info("Checking for existing device and sensor...")
	if !checkForExistingDeviceByID(deviceID, postgreClient) {
		return fmt.Errorf("failed getting device with %s ID", deviceID)
	}

	if !checkForExistingSensorByID(sensorID, postgreClient) {
		return fmt.Errorf("failed getting sensor with %s ID", sensorID)
	}

	return nil
}
