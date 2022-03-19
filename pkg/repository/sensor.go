package repository

import (
	"database/sql"
	"fmt"

	"github.com/Todorov99/server/pkg/dto"
	"github.com/Todorov99/server/pkg/repository/query"
)

type sensorRepository struct {
	postgreClient *sql.DB
}

func (s *sensorRepository) GetAll() (interface{}, error) {
	repositoryLogger.Info("Getting all sensors...")
	var sensors []dto.Sensor

	rowsRs, err := s.postgreClient.Query(query.GetAllSensors)
	if err != nil {
		return nil, err
	}

	defer rowsRs.Close()

	for rowsRs.Next() {

		currentSensor := dto.Sensor{}

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
	var sensors []dto.Sensor

	rowsRs, err := s.postgreClient.Query(query.GetAllSensorsBySensorID, args[0])
	if err != nil {
		return nil, fmt.Errorf("failed executing query %s: %w", fmt.Sprintf(query.GetAllSensorsBySensorID, args[0]), err)
	}

	defer rowsRs.Close()

	for rowsRs.Next() {

		currentSensor := dto.Sensor{}

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
	chackForExistingSensor, err := getSensorIDByName(args[0], s.postgreClient)
	if err != nil {
		return err
	}

	if chackForExistingSensor != "" {
		return fmt.Errorf("sensor with name: %s already exists", args[0])
	}

	sensorGroupID, err := getSensorGroupByName(args[3], s.postgreClient)
	if err != nil {
		return err
	}

	fmt.Println(sensorGroupID)
	return executeModifyingQuery(query.AddSensor, s.postgreClient, args[0], args[1], args[2], sensorGroupID, args[4])
}

func (s *sensorRepository) Update(args ...string) error {
	repositoryLogger.Infof("Updating sensor with id: %s", args[4])
	if !checkForExistingSensorByID(args[4], s.postgreClient) {
		return fmt.Errorf("sensor with id: %s does not exist", args[4])
	}

	return updateSensorsByID(args[4], args[0], args[1], args[2], args[3], s.postgreClient)
}

func (s *sensorRepository) Delete(id string) (interface{}, error) {
	repositoryLogger.Infof("Deleting sensot with id: %s", id)
	if !checkForExistingSensorByID(id, s.postgreClient) {
		return nil, fmt.Errorf("sensor with id: %s does not exist", id)
	}

	deletedSensor, err := getSensorByID(id, s.postgreClient)
	if err != nil {
		return nil, fmt.Errorf("failed deleting sensor with id: %s: err: %w", id, err)
	}

	return deletedSensor, executeModifyingQuery(query.DeleteSensor, s.postgreClient, id)
}

func updateSensorsByID(sensorID string, name string, description string, unit string, sensorGroup string, postgreClient *sql.DB) error {

	sensorGroupID, err := getSensorGroupByName(sensorGroup, postgreClient)

	if sensorGroupID == "" {
		return fmt.Errorf("invalid sensor group")
	}

	if err != nil {
		return err
	}

	return executeModifyingQuery(query.UpdateSensor, postgreClient, name, description, sensorGroupID, unit, sensorID)
}

func getSensorIDByName(name string, postgreClient *sql.DB) (string, error) {
	sensorID := ""
	err := executeSelectQuery(query.GetSensorByName, postgreClient, &sensorID, name)
	if err != nil {
		return "", err
	}

	return sensorID, nil
}

func getSensorGroupByName(sensorGroup string, postgreClient *sql.DB) (string, error) {
	sensorGroupID := ""
	err := executeSelectQuery(query.GetSensorIDByGroupName, postgreClient, &sensorGroupID, sensorGroup)
	if err != nil {
		return "", err
	}

	return sensorGroupID, nil
}

func checkForExistingSensorsByDeviceID(deviceID string, postgreClient *sql.DB) bool {
	sensor := ""
	_ = executeSelectQuery(query.GetSensorIDByDeviceID, postgreClient, &sensor, deviceID)
	return sensor != ""
}

func getSensorByID(sensorID string, posgreClient *sql.DB) (dto.Sensor, error) {
	var sensor dto.Sensor

	rowsRs, err := posgreClient.Query(query.GetSensorByID, sensorID)

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

func getSensorByDeviceID(deviceID string, posgreClient *sql.DB) ([]dto.Sensor, error) {
	var sensors []dto.Sensor

	rowsRs, err := posgreClient.Query(query.GetAllSensorsByDeviceID, deviceID)

	if err != nil {
		return nil, err
	}

	defer rowsRs.Close()

	for rowsRs.Next() {

		currentSensor := dto.Sensor{}

		err := rowsRs.Scan(&currentSensor.ID, &currentSensor.Name, &currentSensor.Description,
			&currentSensor.Unit, &currentSensor.SensorGroups)
		if err != nil {
			return nil, err
		}

		sensors = append(sensors, currentSensor)

	}

	return sensors, nil
}
