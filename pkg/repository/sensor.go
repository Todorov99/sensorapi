package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Todorov99/sensorapi/pkg/dto"
	"github.com/Todorov99/sensorapi/pkg/entity"
	"github.com/Todorov99/sensorapi/pkg/repository/query"
	"github.com/Todorov99/sensorapi/pkg/server/config"
)

type SensorRepository interface {
	GetSensorIDByName(ctx context.Context, name string) (string, error)
	GetSensorGroupByName(ctx context.Context, sensorGroup string) (string, error)
	GetSensorByDeviceID(ctx context.Context, deviceID string) ([]entity.Sensor, error)
	IRepository
}

type sensorRepository struct {
	postgreClient *sql.DB
}

func NewSensorRepository() SensorRepository {
	return &sensorRepository{
		postgreClient: config.GetDatabaseCfg().GetPostgreClient(),
	}
}

func (s *sensorRepository) GetAll(ctx context.Context) (interface{}, error) {
	repositoryLogger.Info("Getting all sensors...")
	var sensors []entity.Sensor
	err := executeSelectQuery(ctx, query.GetAllSensors, s.postgreClient, &sensors)
	if err != nil {
		return nil, err
	}

	return sensors, nil
}

func (s *sensorRepository) GetByID(ctx context.Context, id int) (interface{}, error) {
	repositoryLogger.Debugf("Getting sensor by ID: %d", id)
	var sensors entity.Sensor

	err := executeSelectQuery(ctx, query.GetSensorByID, s.postgreClient, &sensors, id)
	if err != nil {
		return nil, err
	}

	return sensors, nil
}

func (s *sensorRepository) Add(ctx context.Context, model interface{}) error {
	sensor := model.(dto.Sensor)
	repositoryLogger.Infof("Adding sensor with name: %s...", sensor.Name)

	sensorGroupID, err := s.GetSensorGroupByName(ctx, sensor.SensorGroups)
	if err != nil {
		return err
	}

	return executeModifyingQuery(ctx, query.AddSensor, s.postgreClient, sensor.Name, sensor.Description, sensor.DeviceId, sensorGroupID, sensor.Unit)
}

func (s *sensorRepository) Update(ctx context.Context, model interface{}) error {
	sensor := model.(dto.Sensor)
	sensorGroupID, err := s.GetSensorGroupByName(ctx, sensor.SensorGroups)
	if err != nil {
		return err
	}

	if sensorGroupID == "" {
		return fmt.Errorf("invalid sensor group")
	}

	return executeModifyingQuery(ctx, query.UpdateSensor, s.postgreClient, sensor.Name, sensor.Description, sensorGroupID, sensor.Unit, sensor.ID)
}

func (s *sensorRepository) Delete(ctx context.Context, id int) error {
	repositoryLogger.Infof("Deleting sensor with id: %s", id)
	return executeModifyingQuery(ctx, query.DeleteSensor, s.postgreClient, id)
}

func (s *sensorRepository) GetSensorIDByName(ctx context.Context, name string) (string, error) {
	sensorID := ""
	err := executeSelectQuery(ctx, query.GetSensorByName, s.postgreClient, &sensorID, name)
	if err != nil {
		return "", err
	}

	return sensorID, nil
}

func (s *sensorRepository) GetSensorGroupByName(ctx context.Context, sensorGroup string) (string, error) {
	sensorGroupID := ""
	err := executeSelectQuery(ctx, query.GetSensorIDByGroupName, s.postgreClient, &sensorGroupID, sensorGroup)
	if err != nil {
		return "", err
	}

	return sensorGroupID, nil
}

// func (s *sensorRepository) GetSensorByID(ctx context.Context, id string) (entity.Sensor, error) {
// 	var sensor entity.Sensor
// 	err := executeSelectQuery(query.GetSensorByID, s.postgreClient, &sensor)
// 	if err != nil {
// 		return entity.Sensor{}, err
// 	}

// 	return sensor, nil
// }

func (s *sensorRepository) GetSensorByDeviceID(ctx context.Context, deviceID string) ([]entity.Sensor, error) {
	var sensors []entity.Sensor
	err := executeSelectQuery(ctx, query.GetAllSensors, s.postgreClient, &sensors)
	if err != nil {
		return nil, err
	}

	return sensors, nil
}
