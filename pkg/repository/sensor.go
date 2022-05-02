package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/Todorov99/sensorapi/pkg/dto"
	"github.com/Todorov99/sensorapi/pkg/entity"
	"github.com/Todorov99/sensorapi/pkg/repository/query"
	"github.com/Todorov99/sensorapi/pkg/server/config"
	"github.com/Todorov99/sensorcli/pkg/logger"
	"github.com/sirupsen/logrus"
)

type SensorRepository interface {
	GetSensorIDByName(ctx context.Context, name string) (string, error)
	GetSensorGroupByName(ctx context.Context, sensorGroup string) (string, error)
	GetSensorByDeviceID(ctx context.Context, deviceID string) ([]entity.Sensor, error)
	GetAll(ctx context.Context) (interface{}, error)
	GetByID(ctx context.Context, sensorID int) (interface{}, error)
	Update(ctx context.Context, entity interface{}) error
}

type sensorRepository struct {
	logger        *logrus.Entry
	postgreClient *sql.DB
}

func NewSensorRepository() SensorRepository {
	return &sensorRepository{
		logger:        logger.NewLogrus("deviceRepository", os.Stdout),
		postgreClient: config.GetDatabaseCfg().GetPostgreClient(),
	}
}

func (s *sensorRepository) GetAll(ctx context.Context) (interface{}, error) {
	s.logger.Info("Getting all sensors...")
	var sensors []entity.Sensor
	err := executeSelectQuery(ctx, query.GetAllSensors, s.postgreClient, &sensors)
	if err != nil {
		return nil, err
	}

	return sensors, nil
}

func (s *sensorRepository) GetByID(ctx context.Context, sensorID int) (interface{}, error) {
	s.logger.Debugf("Getting sensor by ID: %d", sensorID)
	var sensors entity.Sensor

	err := executeSelectQuery(ctx, query.GetSensorByID, s.postgreClient, &sensors, sensorID)
	if err != nil {
		return nil, err
	}

	return sensors, nil
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

func (s *sensorRepository) GetSensorByDeviceID(ctx context.Context, deviceID string) ([]entity.Sensor, error) {
	var sensors []entity.Sensor
	err := executeSelectQuery(ctx, query.GetAllSensors, s.postgreClient, &sensors)
	if err != nil {
		return nil, err
	}

	return sensors, nil
}
