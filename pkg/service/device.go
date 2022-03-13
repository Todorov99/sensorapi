package service

import (
	"github.com/Todorov99/server/pkg/models"
	"github.com/Todorov99/server/pkg/repository"
	"github.com/mitchellh/mapstructure"
)

type deviceService struct {
	deviceRepository repository.Repository
}

func NewDeviceService() IService {
	return &deviceService{
		deviceRepository: repository.CreateDeviceRepository(),
	}
}

func (d *deviceService) GetAll() (interface{}, error) {
	return d.deviceRepository.GetAll()
}

func (d *deviceService) GetById(deviceID string) (interface{}, error) {
	return d.deviceRepository.GetByID(deviceID)
}

func (d *deviceService) Add(model interface{}) error {
	device := models.Device{}
	err := mapstructure.Decode(model, &device)
	if err != nil {
		return err
	}
	return d.deviceRepository.Add(device.Name, device.Description)
}

func (d *deviceService) Update(model interface{}) error {
	device := models.Device{}
	err := mapstructure.Decode(model, &device)
	if err != nil {
		return err
	}
	return d.deviceRepository.Update(device.Name, device.Description, device.ID)
}

func (d *deviceService) Delete(deviceID string) (interface{}, error) {
	return d.deviceRepository.Delete(deviceID)
}
