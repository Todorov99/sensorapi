package service

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/Todorov99/sensorapi/pkg/dto"
	"github.com/Todorov99/sensorapi/pkg/entity"
	"github.com/Todorov99/sensorapi/pkg/global"
	"github.com/Todorov99/sensorapi/pkg/repository"
	"github.com/Todorov99/sensorcli/pkg/logger"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type DeviceService interface {
	GetAll(ctx context.Context, userID int) (interface{}, error)
	GetById(ctx context.Context, ID, userID int) (interface{}, error)
	Add(ctx context.Context, model interface{}, userID int) error
	Update(ctx context.Context, model interface{}, userID int) error
	Delete(ctx context.Context, deviceID, userID int) (interface{}, error)
	GenerateDeviceCfg(ctx context.Context, deviceID, userID int) (string, error)
}

type deviceService struct {
	logger           *logrus.Entry
	deviceRepository repository.DeviceRepository
	sensorRepository repository.SensorRepository
	userRepository   repository.UserRepository
}

func NewDeviceService() DeviceService {
	return &deviceService{
		logger:           logger.NewLogrus("deviceService", os.Stdout),
		deviceRepository: repository.NewDeviceRepository(),
		sensorRepository: repository.NewSensorRepository(),
		userRepository:   repository.NewUserRepository(),
	}
}

func (d *deviceService) GenerateDeviceCfg(ctx context.Context, deviceID, userID int) (string, error) {
	d.logger.Debug("Generating device cfg...")
	cfgFileName := "device_cfg.yaml"
	dd, err := d.GetById(ctx, deviceID, userID)
	if err != nil {
		return "", err
	}

	f, err := os.Create(cfgFileName)
	if err != nil {
		return "", err
	}

	defer func() {
		f.Close()
	}()

	device := dto.Device{}
	err = mapstructure.Decode(dd, &device)
	if err != nil {
		return "", err
	}

	deviceBytes, err := yaml.Marshal(device)
	if err != nil {
		return "", err
	}

	_, err = f.Write(deviceBytes)
	if err != nil {
		return "", err
	}

	d.logger.Debug("Device cfg successfully generated")
	return cfgFileName, nil
}

func (d *deviceService) GetAll(ctx context.Context, userID int) (interface{}, error) {
	d.logger.Debug("Getting all devices")

	devices, err := d.deviceRepository.GetAll(ctx, userID)
	if err != nil {
		return nil, err
	}

	allDevices := []dto.Device{}
	err = mapstructure.Decode(devices, &allDevices)
	if err != nil {
		return nil, err
	}

	return allDevices, nil
}

func (d *deviceService) GetById(ctx context.Context, deviceID, userID int) (interface{}, error) {
	entityDevice, err := d.deviceRepository.GetByID(ctx, deviceID, userID)
	if err != nil {
		if errors.Is(err, global.ErrorObjectNotFound) {
			return nil, fmt.Errorf("device with ID: %d does not exist", deviceID)
		}
		return nil, err
	}

	device := dto.Device{}
	err = mapstructure.Decode(entityDevice, &device)
	if err != nil {
		return nil, err
	}

	return device, nil
}

func (d *deviceService) Add(ctx context.Context, model interface{}, userID int) error {
	device := entity.Device{}
	err := mapstructure.Decode(model, &device)
	if err != nil {
		return err
	}

	err = d.ifDeviceExist(ctx, device.Name)
	if errors.Is(err, global.ErrorDeviceWithNameAlreadyExist) {
		return fmt.Errorf("device with name %s already exists", device.Name)
	}

	if err != nil && !errors.Is(err, global.ErrorDeviceWithNameAlreadyExist) {
		return err
	}

	err = d.deviceRepository.Add(ctx, device, userID)
	if err != nil {
		return err
	}

	deviceID, err := d.deviceRepository.GetDeviceIDByName(ctx, device.Name)
	if err != nil {
		return err
	}

	device.ID = deviceID

	sensors, err := d.sensorRepository.GetAll(ctx)
	if err != nil {
		return err
	}

	device.Sensors = sensors.([]entity.Sensor)

	for _, sensor := range device.Sensors {
		err := d.deviceRepository.AddDeviceSensors(ctx, device.ID, sensor.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *deviceService) Update(ctx context.Context, model interface{}, userID int) error {
	device := entity.Device{}
	err := mapstructure.Decode(model, &device)
	if err != nil {
		return err
	}
	d.logger.Debugf("Updating device with ID: %d", device.ID)

	_, err = d.deviceRepository.GetByID(ctx, int(device.ID), userID)
	if err != nil {
		if errors.Is(err, global.ErrorObjectNotFound) {
			return fmt.Errorf("device with id: %d does not exist", device.ID)
		}
		return err
	}

	return d.deviceRepository.Update(ctx, device, userID)
}

func (d *deviceService) Delete(ctx context.Context, deviceID, userID int) (interface{}, error) {
	if deviceID == 1 && userID == 1 {
		return nil, fmt.Errorf("your are not allowed to delete your default device")
	}

	deviceForDelete, err := d.deviceRepository.GetByID(ctx, deviceID, userID)
	if err != nil {
		if errors.Is(err, global.ErrorObjectNotFound) {
			return nil, fmt.Errorf("device with id: %d does not exist", deviceID)
		}
		return nil, err
	}

	err = d.deviceRepository.Delete(ctx, deviceID, userID)
	if err != nil {
		return nil, err
	}

	device := dto.Device{}
	err = mapstructure.Decode(deviceForDelete, &device)
	if err != nil {
		return nil, err
	}

	return device, nil
}

func (d *deviceService) ifDeviceExist(ctx context.Context, deviceName string) error {
	checkForExistingDevice, err := d.deviceRepository.GetDeviceIDByName(ctx, deviceName)
	if err != nil && !errors.Is(err, global.ErrorObjectNotFound) {
		return err
	}

	if checkForExistingDevice != 0 {
		return global.ErrorDeviceWithNameAlreadyExist
	}

	return nil
}
