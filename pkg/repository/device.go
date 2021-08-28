package repository

import (
	"encoding/json"
	"errors"

	"github.com/Todorov99/server/pkg/database/postgres"
	"github.com/Todorov99/server/pkg/models"
)

type deviceRepository struct{}

const (
	insertDevice string = `INSERT INTO device(id,name,description) VALUES ($1,$2,$3)`
	updateDevice string = `UPDATE device set name=$1,description=$2 where id=$3`
	deleteDevice string = `DELETE from device where id=$1`

	selectDeviceQuery = `SELECT d.id, d.name, d.description from device as d where d.id=$1`
	getAllDevices     = `SELECT d.id, d.name, d.description from device as d`
)

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

	rowsRs, err := postgres.DatabaseConnection.Query(getAllDevices)

	if len(args) != 0 {
		rowsRs, err = postgres.DatabaseConnection.Query(selectDeviceQuery, args[0])
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

		sensors, _ := CreateSensorRepository().GetByID(currentDevice.ID)

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
		return nil, errors.New("There is not devices with this id")
	}

	return devices, nil
}

func (d *deviceRepository) Add(args ...string) error {

	deviceID, err := executeSelectQuery("SELECT max(id) + 1 from device")

	if err != nil {
		return err
	}

	checkForExistingDevice, err := getDeviceIDByName(args[0])

	if err != nil {
		return err
	}

	if checkForExistingDevice != "" {
		return errors.New("There is device with that name")
	}

	return executeModifyingQuery(insertDevice, deviceID, args[0], args[1])
}

func (d *deviceRepository) Update(args ...string) error {

	if !checkForExistingDeviceByID(args[2]) {
		return errors.New("You've entered incorrect device id")
	}

	return executeModifyingQuery(updateDevice, args[0], args[1], args[2])
}

func (d *deviceRepository) Delete(id string) (interface{}, error) {

	if !checkForExistingDeviceByID(id) {
		return nil, errors.New("You've entered incorrect device id")
	}

	checkForAvailabeSensorsByDeviceID := checkForExistingSensorsByDeviceID(id)

	if checkForAvailabeSensorsByDeviceID == true {
		return nil, errors.New("You must not delete this device. It has sensors who belongs to it")
	}

	deletedDevice, err := d.GetByID(id)

	if err != nil {
		return nil, errors.New("There is not device with this name")
	}

	return deletedDevice, executeModifyingQuery(deleteDevice, id)
}

func getDeviceIDByName(deviceName string) (string, error) {

	query := "SELECT id FROM device where name=$1"

	id, err := executeSelectQuery(query, deviceName)

	if err != nil {
		return "", err
	}

	return id, nil
}
