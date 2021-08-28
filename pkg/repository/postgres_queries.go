package repository

import (
	"errors"

	"github.com/Todorov99/server/pkg/database/postgres"
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
		//	logger.ErrorLogger.Println(err)
		err = errors.New("You've entered incorrect device data")
	}

	return err
}

func checkForExistingSensorByID(id string) bool {

	existingSensor, _ := executeSelectQuery("Select name from sensor where id=$1", id)
	return existingSensor != ""
}

func checkForExistingDeviceByID(id string) bool {

	existingDevice, _ := executeSelectQuery("Select name from device where id=$1", id)

	return existingDevice != ""
}

func checkForExistingDevicesAndSensors(deviceID string, sensorID string) error {

	if !checkForExistingDeviceByID(deviceID) {
		return errors.New("There is not device with such id")
	}

	if !checkForExistingSensorByID(sensorID) {
		return errors.New("There is not sensor with such id")
	}

	return nil
}
