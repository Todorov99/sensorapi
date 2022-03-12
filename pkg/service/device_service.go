package service

import (
	"github.com/Todorov99/server/pkg/models"
	"github.com/Todorov99/server/pkg/repository"
)

type DeviceService interface {
	GetAllDevices() (interface{}, error)
	GetDeviceById(deviceID string) (interface{}, error)
	AddDevice(models.Device) error
	UpdateDevice(device models.Device) error
	DeleteDevice(deviceID string) (interface{}, error)
}

type deviceService struct {
	deviceRepository repository.Repository
}

func NewDeviceService() DeviceService {
	return &deviceService{
		deviceRepository: repository.CreateDeviceRepository(),
	}
}

func (d *deviceService) GetAllDevices() (interface{}, error) {
	return d.deviceRepository.GetAll()
}

func (d *deviceService) GetDeviceById(deviceID string) (interface{}, error) {
	return d.deviceRepository.GetByID(deviceID)
}

func (d *deviceService) AddDevice(device models.Device) error {
	return d.deviceRepository.Add(device.Name, device.Description)
}

func (d *deviceService) UpdateDevice(device models.Device) error {
	return d.deviceRepository.Update(device.Name, device.Description, device.ID)
}

func (d *deviceService) DeleteDevice(deviceID string) (interface{}, error) {
	return d.deviceRepository.Delete(deviceID)
}
