package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/Todorov99/sensorapi/pkg/entity"
	"github.com/Todorov99/sensorapi/pkg/global"
	"github.com/Todorov99/sensorapi/pkg/repository/query"
	"github.com/Todorov99/sensorapi/pkg/server/config"
	"github.com/Todorov99/sensorcli/pkg/logger"
	"github.com/sirupsen/logrus"
)

type DeviceRepository interface {
	GetDeviceNameByID(ctx context.Context, id int) (string, error)
	GetDeviceIDByName(ctx context.Context, deviceName string) (int32, error)
	AddDeviceSensors(ctx context.Context, deviceID, sensorID int32) error
	GetAll(ctx context.Context, userID int) (interface{}, error)
	GetByID(ctx context.Context, id, userID int) (interface{}, error)
	Add(ctx context.Context, entity interface{}, userID int) error
	Update(ctx context.Context, entity interface{}, userID int) error
	Delete(ctx context.Context, id, userID int) error
}

type deviceRepository struct {
	logger        *logrus.Entry
	postgreClient *sql.DB
}

func NewDeviceRepository() DeviceRepository {
	return &deviceRepository{
		logger:        logger.NewLogrus("deviceRepository", os.Stdout),
		postgreClient: config.GetDatabaseCfg().GetPostgreClient(),
	}
}

func (d *deviceRepository) GetAll(ctx context.Context, userID int) (interface{}, error) {
	d.logger.Debug("Getting all devicess...")
	devices := []*entity.Device{}

	err := executeSelectQuery(ctx, query.GetAllDevices, d.postgreClient, &devices, userID)
	if err != nil {
		return nil, err
	}

	for _, device := range devices {
		sensors := []entity.Sensor{}
		err = executeSelectQuery(ctx, query.GetAllSensorsByDeviceID, d.postgreClient, &sensors, device.ID)
		if err != nil {
			return nil, err
		}

		device.Sensors = append(device.Sensors, sensors...)
	}
	d.logger.Debug("Devices successfully retrieved")
	return devices, nil
}

func (d *deviceRepository) Add(ctx context.Context, model interface{}, userID int) error {
	d.logger.Info("Adding device with all predifined sensors...")
	device := model.(entity.Device)

	err := executeModifyingQuery(ctx, query.InsertDevice, d.postgreClient, device.Name, device.Description, userID)
	if err != nil {
		return err
	}

	d.logger.Info("Device sucessfully added")
	return nil
}

func (d *deviceRepository) AddDeviceSensors(ctx context.Context, deviceID, sensorID int32) error {
	return executeModifyingQuery(ctx, query.InsertDeviceSensors, d.postgreClient, deviceID, sensorID)
}

func (d *deviceRepository) Update(ctx context.Context, model interface{}, userID int) error {
	device := model.(entity.Device)
	return executeModifyingQuery(ctx, query.UpdateDevice, d.postgreClient, device.Name, device.Description, device.ID, userID)
}

func (d *deviceRepository) GetByID(ctx context.Context, id, userID int) (interface{}, error) {
	d.logger.Debugf("Getting device by ID: %d", id)
	device := &entity.Device{}

	err := executeSelectQuery(ctx, query.GetDeviceByID, d.postgreClient, device, id, userID)
	if err != nil {
		return nil, err
	}

	sensors := []entity.Sensor{}
	err = executeSelectQuery(ctx, query.GetAllSensorsByDeviceID, d.postgreClient, &sensors, device.ID)
	if err != nil {
		return nil, err
	}

	device.Sensors = append(device.Sensors, sensors...)
	d.logger.Debug("Devices successfully retrieved")

	return device, nil
}

func (d *deviceRepository) Delete(ctx context.Context, id, userID int) error {
	d.logger.Infof("Deleting device with id: %q", id)
	return executeModifyingQuery(ctx, query.DeleteDevice, d.postgreClient, id, userID)
}

func (d *deviceRepository) GetDeviceIDByName(ctx context.Context, deviceName string) (int32, error) {
	d.logger.Infof("Getting device ID by name: %q", deviceName)
	var id int32
	err := executeSelectQuery(ctx, query.GetDeviceIDByName, d.postgreClient, &id, deviceName)
	if err != nil {
		return 0, fmt.Errorf("failed getting device with name: %q: %w", deviceName, err)
	}

	if id == 0 {
		return 0, global.ErrorObjectNotFound
	}

	return id, nil
}

func (d *deviceRepository) GetDeviceNameByID(ctx context.Context, id int) (string, error) {
	var deviceName string
	err := executeSelectQuery(ctx, query.GetDeviceNameByID, d.postgreClient, &deviceName, id)
	if err != nil {
		return "", err
	}

	if deviceName == "" {
		return "", global.ErrorObjectNotFound
	}

	return deviceName, nil
}
