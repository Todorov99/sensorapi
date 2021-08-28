package repository

import (
	"errors"

	"github.com/Todorov99/server/pkg/database/postgres"
	"github.com/Todorov99/server/pkg/models"
)

type sensorRepository struct{}

const (
	addSensor    string = `INSERT INTO sensor(id,name,description,device_id,sensor_groups_id,unit) VALUES ($1,$2,$3,$4,$5,$6)`
	updateSensor string = `UPDATE sensor set name=$1,description=$2,sensor_groups_id=$3,unit=$4 where id=$5`
	deleteSensor string = `DELETE from sensor where id=$1`
)

func (s *sensorRepository) GetAll() (interface{}, error) {
	return nil, errors.New("Not Implemented")
}

func (s *sensorRepository) GetByID(args ...string) (interface{}, error) {

	var sensors []models.Sensor

	rowsRs, err := postgres.DatabaseConnection.Query("SELECT s.id, s.name, s.description, s.unit, ss.group_name FROM sensor as s join sensor_groups as ss on s.sensor_groups_id = ss.id")

	if len(args) > 0 {
		rowsRs, err = postgres.DatabaseConnection.Query("SELECT s.id, s.name, s.description, s.unit, ss.group_name FROM sensor as s join sensor_groups as ss on s.sensor_groups_id = ss.id where s.device_id = $1", args[0])
	}

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

func (s *sensorRepository) Add(args ...string) error {

	chackForExistingSensor, err := getSensorIDByName(args[0])

	if err != nil {
		return err
	}

	if chackForExistingSensor != "" {
		return errors.New("There is sensor with that name")
	}

	sensorGroupID, sensorGroupError := getSensorGroupByName(args[3])
	sensorID, err := executeSelectQuery("SELECT max(id) + 1 from sensor")

	if err != nil {
		return err
	}

	if sensorGroupError != nil {
		return sensorGroupError
	}

	return executeModifyingQuery(addSensor, sensorID, args[0], args[1], args[2], sensorGroupID, args[4])
}

func (s *sensorRepository) Update(args ...string) error {

	if !checkForExistingSensorByID(args[4]) {
		return errors.New("There is not sensor with this id")
	}

	return updateSensorsByID(args[4], args[0], args[1], args[2], args[3])
}

func (s *sensorRepository) Delete(id string) (interface{}, error) {

	if !checkForExistingSensorByID(id) {
		return nil, errors.New("There is not sensor with this id")
	}

	deletedSensor, err := getSensorByID(id)

	if err != nil {
		return nil, err
	}

	return deletedSensor, executeModifyingQuery(deleteSensor, id)
}

func updateSensorsByID(sensorID string, name string, description string, unit string, sensorGroup string) error {

	sensorGroupID, err := getSensorGroupByName(sensorGroup)

	if sensorGroupID == "" {
		return errors.New("You've entered invalid sensor group")
	}

	if err != nil {
		return err
	}

	return executeModifyingQuery(updateSensor, name, description, sensorGroupID, unit, sensorID)
}

func getSensorIDByName(name string) (string, error) {

	sensorID, err := executeSelectQuery("SELECT id from sensor where name=$1", name)

	if err != nil {
		return "", err
	}

	return sensorID, nil
}

func getSensorGroupByName(sensorGroup string) (string, error) {
	return executeSelectQuery("SELECT id from sensor_groups where group_name=$1", sensorGroup)
}

func checkForExistingSensorsByDeviceID(deviceID string) bool {

	sensor, _ := executeSelectQuery("SELECT s.id from sensor as s where s.device_id=$1", deviceID)

	if sensor != "" {
		return true
	}

	return false
}

func getSensorByID(sensorID string) (models.Sensor, error) {
	var sensor models.Sensor

	rowsRs, err := postgres.DatabaseConnection.Query("SELECT s.id, s.name, s.description, s.unit, ss.group_name FROM sensor as s join sensor_groups as ss on s.sensor_groups_id = ss.id where s.id=$1", sensorID)

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
