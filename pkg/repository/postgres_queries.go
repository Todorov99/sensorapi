package repository

import (
	"fmt"

	"github.com/Todorov99/server/pkg/database/postgres"
	"github.com/Todorov99/server/pkg/repository/query"
)

func executeSelectQuery(query string, args ...interface{}) (string, error) {

	var value string

	rowsRs, err := postgres.DatabaseConnection.Query(query, args...)

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

func executeModifyingQuery(query string, args ...interface{}) error {
	_, err := postgres.DatabaseConnection.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed executing query %q with arguments %q: %w", query, args, err)
	}

	return nil
}

func checkForExistingSensorByID(id string) bool {
	existingSensor, _ := executeSelectQuery(query.GetSensorNameByID, id)
	return existingSensor != ""
}

func checkForExistingDeviceByID(id string) bool {
	existingDevice, _ := executeSelectQuery(query.GetDeviceNameByID, id)
	return existingDevice != ""
}

func checkForExistingDevicesAndSensors(deviceID string, sensorID string) error {
	repositoryLogger.Info("Checking for existing device and sensor...")
	if !checkForExistingDeviceByID(deviceID) {
		return fmt.Errorf("failed getting device with %s ID", deviceID)
	}

	if !checkForExistingSensorByID(sensorID) {
		return fmt.Errorf("failed getting sensor with %s ID", sensorID)
	}

	return nil
}
