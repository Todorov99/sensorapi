package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Todorov99/server/pkg/dto"
	"github.com/Todorov99/server/pkg/entity"
	"github.com/Todorov99/server/pkg/global"
	"github.com/Todorov99/server/pkg/repository"
	"github.com/mitchellh/mapstructure"
)

type deviceService struct {
	deviceRepository repository.DeviceRepository
}

func NewDeviceService() IService1 {
	return &deviceService{
		deviceRepository: repository.NewDeviceRepository(),
	}
}

func (d *deviceService) GetAll(ctx context.Context) (interface{}, error) {
	serviceLogger.Debug("Getting all devices")
	devices, err := d.deviceRepository.GetAll(ctx)
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

func (d *deviceService) GetById(ctx context.Context, deviceID int) (interface{}, error) {
	entityDevice, err := d.deviceRepository.GetByID(ctx, deviceID)
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

func (d *deviceService) Add(ctx context.Context, model interface{}) error {
	device := entity.Device{}
	err := mapstructure.Decode(model, &device)
	if err != nil {
		return err
	}

	checkForExistingDevice, err := d.deviceRepository.GetDeviceIDByName(ctx, device.Name)
	if err != nil {
		return err
	}

	if checkForExistingDevice != "" {
		return fmt.Errorf("device with name: %q already exists", device.Name)
	}

	return d.deviceRepository.Add(ctx, device)
}

func (d *deviceService) Update(ctx context.Context, model interface{}) error {
	device := entity.Device{}
	err := mapstructure.Decode(model, &device)
	if err != nil {
		return err
	}
	serviceLogger.Debugf("Updating device with ID: %d", device.ID)

	_, err = d.deviceRepository.GetByID(ctx, int(device.ID))
	if err != nil {
		if errors.Is(err, global.ErrorObjectNotFound) {
			return fmt.Errorf("device with id: %d does not exist", device.ID)
		}
		return err
	}

	return d.deviceRepository.Update(ctx, device)
}

func (d *deviceService) Delete(ctx context.Context, deviceID int) (interface{}, error) {
	deviceForDelete, err := d.deviceRepository.GetByID(ctx, deviceID)
	if err != nil {
		if errors.Is(err, global.ErrorObjectNotFound) {
			return nil, fmt.Errorf("device with id: %d does not exist", deviceID)
		}
		return nil, err
	}

	//TODO check whether the device has connected sensors
	err = d.deviceRepository.Delete(ctx, deviceID)
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
