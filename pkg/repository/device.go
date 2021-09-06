package repository

import (
	"encoding/json"
	"fmt"

	"github.com/Todorov99/server/pkg/database/postgres"
	"github.com/Todorov99/server/pkg/models"
	"github.com/Todorov99/server/pkg/repository/query"
)

type deviceRepository struct{}

func (d *deviceRepository) GetAll() (interface{}, error) {
	devices, err := d.GetByID()
	if err != nil {
		return nil, err
	}

	return devices, nil
}

func (d *deviceRepository) GetByID(args ...string) (interface{}, error) {

	var devices []models.Device
	deviceSensors := []models.Sensor{}

	rowsRs, err := postgres.DatabaseConnection.Query(query.GetAllDevices)

	if len(args) != 0 {
		rowsRs, err = postgres.DatabaseConnection.Query(query.GetDeviceByID, args[0])
	}

	if err != nil {
		return nil, err
	}

	for rowsRs.Next() {

		currentDevice := models.Device{}

		err := rowsRs.Scan(&currentDevice.ID, &currentDevice.Name, &currentDevice.Description)
		if err != nil {
			return nil, err
		}

		sensors, _ := getSensorByDeviceID(currentDevice.ID)

		sensorBytes, err := json.Marshal(sensors)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(sensorBytes, &deviceSensors)
		if err != nil {
			return nil, err
		}

		currentDevice.Sensors = deviceSensors
		devices = append(devices, currentDevice)

	}

	if devices == nil {
		return nil, fmt.Errorf("failed to get devices")
	}

	return devices, nil
}

func (d *deviceRepository) Add(args ...string) error {
	repositoryLogger.Info("Adding device...")
	deviceID, err := executeSelectQuery(query.GetHighestDeviceID)
	if err != nil {
		return err
	}

	checkForExistingDevice, err := getDeviceIDByName(args[0])
	if err != nil {
		return err
	}

	if checkForExistingDevice != "" {
		return fmt.Errorf("device with name: %q already exists", args[0])
	}

	return executeModifyingQuery(query.InsertDevice, deviceID, args[0], args[1])
}

func (d *deviceRepository) Update(args ...string) error {
	repositoryLogger.Info("Updating device with id: %s", args[2])
	if !checkForExistingDeviceByID(args[2]) {
		return fmt.Errorf("invalid device id: %q", args[2])
	}

	return executeModifyingQuery(query.UpdateDevice, args[0], args[1], args[2])
}

func (d *deviceRepository) Delete(id string) (interface{}, error) {
	repositoryLogger.Infof("Deleting device with id: %q", id)
	if !checkForExistingDeviceByID(id) {
		return nil, fmt.Errorf("invalid device id: %q", id)
	}

	checkForAvailabeSensorsByDeviceID := checkForExistingSensorsByDeviceID(id)

	if checkForAvailabeSensorsByDeviceID {
		return nil, fmt.Errorf("failed to delete device with ID: %q. First you have to delete the sensors that belongs to it", id)
	}

	deletedDevice, err := d.GetByID(id)
	if err != nil {
		return nil, err
	}

	return deletedDevice, executeModifyingQuery(query.DeleteDevice, id)
}

func getDeviceIDByName(deviceName string) (string, error) {
	repositoryLogger.Infof("Getting device ID by name: %q", deviceName)
	id, err := executeSelectQuery(query.GetDeviceIDByName, deviceName)
	if err != nil {
		return "", fmt.Errorf("failed getting device with name: %q", deviceName)
	}

	return id, nil
}
