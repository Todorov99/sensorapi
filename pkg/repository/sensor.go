package repository

import (
	"fmt"

	"github.com/Todorov99/server/pkg/database/postgres"
	"github.com/Todorov99/server/pkg/models"
	"github.com/Todorov99/server/pkg/repository/query"
)

type sensorRepository struct{}

func (s *sensorRepository) GetAll() (interface{}, error) {
	repositoryLogger.Info("Getting all sensors...")
	var sensors []models.Sensor

	rowsRs, err := postgres.DatabaseConnection.Query(query.GetAllSensors)
	if err != nil {
		return nil, err
	}

	defer rowsRs.Close()

	for rowsRs.Next() {

		currentSensor := models.Sensor{}

		err := rowsRs.Scan(&currentSensor.ID, &currentSensor.Name, &currentSensor.Description,
			&currentSensor.Unit, &currentSensor.SensorGroups)
		if err != nil {
			return nil, err
		}

		sensors = append(sensors, currentSensor)
	}

	return sensors, nil
}

func (s *sensorRepository) GetByID(args ...string) (interface{}, error) {
	repositoryLogger.Debugf("Getting sensor by ID: %s", args[0])
	var sensors []models.Sensor

	rowsRs, err := postgres.DatabaseConnection.Query(query.GetAllSensorsBySensorID, args[0])
	if err != nil {
		return nil, fmt.Errorf("failed executing query %s: %w", fmt.Sprintf(query.GetAllSensorsBySensorID, args[0]), err)
	}

	defer rowsRs.Close()

	for rowsRs.Next() {

		currentSensor := models.Sensor{}

		err := rowsRs.Scan(&currentSensor.ID, &currentSensor.Name, &currentSensor.Description,
			&currentSensor.Unit, &currentSensor.SensorGroups)
		if err != nil {
			return nil, err
		}

		sensors = append(sensors, currentSensor)
	}

	return sensors, nil
}

func (s *sensorRepository) Add(args ...string) error {
	repositoryLogger.Infof("Adding sensor with name: %s...", args[0])
	chackForExistingSensor, err := getSensorIDByName(args[0])

	if err != nil {
		return err
	}

	if chackForExistingSensor != "" {
		return fmt.Errorf("sensor with name: %s already exists", args[0])
	}

	sensorGroupID, sensorGroupError := getSensorGroupByName(args[3])
	sensorID, err := executeSelectQuery("SELECT max(id) + 1 from sensor")

	if err != nil {
		return err
	}

	if sensorGroupError != nil {
		return sensorGroupError
	}

	return executeModifyingQuery(query.AddSensor, sensorID, args[0], args[1], args[2], sensorGroupID, args[4])
}

func (s *sensorRepository) Update(args ...string) error {
	repositoryLogger.Infof("Updating sensor with id: %s", args[4])
	if !checkForExistingSensorByID(args[4]) {
		return fmt.Errorf("sensor with id: %s does not exist", args[4])
	}

	return updateSensorsByID(args[4], args[0], args[1], args[2], args[3])
}

func (s *sensorRepository) Delete(id string) (interface{}, error) {
	repositoryLogger.Infof("Deleting sensot with id: %s", id)
	if !checkForExistingSensorByID(id) {
		return nil, fmt.Errorf("sensor with id: %s does not exist", id)
	}

	deletedSensor, err := getSensorByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed deleting sensor with id: %s: err: %w", id, err)
	}

	return deletedSensor, executeModifyingQuery(query.DeleteSensor, id)
}

func updateSensorsByID(sensorID string, name string, description string, unit string, sensorGroup string) error {

	sensorGroupID, err := getSensorGroupByName(sensorGroup)

	if sensorGroupID == "" {
		return fmt.Errorf("invalid sensor group")
	}

	if err != nil {
		return err
	}

	return executeModifyingQuery(query.UpdateSensor, name, description, sensorGroupID, unit, sensorID)
}

func getSensorIDByName(name string) (string, error) {

	sensorID, err := executeSelectQuery(query.GetSensorByName, name)

	if err != nil {
		return "", err
	}

	return sensorID, nil
}

func getSensorGroupByName(sensorGroup string) (string, error) {
	return executeSelectQuery(query.GetSensorIDByGroupName, sensorGroup)
}

func checkForExistingSensorsByDeviceID(deviceID string) bool {
	sensor, _ := executeSelectQuery(query.GetSensorIDByDeviceID, deviceID)
	return sensor != ""
}

func getSensorByID(sensorID string) (models.Sensor, error) {
	var sensor models.Sensor

	rowsRs, err := postgres.DatabaseConnection.Query(query.GetSensorByID, sensorID)

	if err != nil {
		return sensor, err
	}

	defer rowsRs.Close()

	for rowsRs.Next() {

		err := rowsRs.Scan(&sensor.ID, &sensor.Name, &sensor.Description,
			&sensor.Unit, &sensor.SensorGroups)

		if err != nil {
			return sensor, err
		}

	}

	return sensor, nil
}

func getSensorByDeviceID(deviceID string) ([]models.Sensor, error) {
	var sensors []models.Sensor

	rowsRs, err := postgres.DatabaseConnection.Query(query.GetAllSensorsByDeviceID, deviceID)

	if err != nil {
		return nil, err
	}

	defer rowsRs.Close()

	for rowsRs.Next() {

		currentSensor := models.Sensor{}

		err := rowsRs.Scan(&currentSensor.ID, &currentSensor.Name, &currentSensor.Description,
			&currentSensor.Unit, &currentSensor.SensorGroups)
		if err != nil {
			return nil, err
		}

		sensors = append(sensors, currentSensor)

	}

	return sensors, nil
}
