package service

import (
	"github.com/Todorov99/server/pkg/dto"
	"github.com/Todorov99/server/pkg/repository"
	"github.com/mitchellh/mapstructure"
)

type sensorService struct {
	sensorRepository repository.Repository
}

func NewSensorService() IService {
	return &sensorService{
		sensorRepository: repository.CreateSensorRepository(),
	}
}

func (s *sensorService) GetAll() (interface{}, error) {
	sensors, err := s.sensorRepository.GetAll()
	if err != nil {
		return nil, err
	}

	return sensors, nil
}

func (s *sensorService) GetById(senosorID string) (interface{}, error) {
	sensor, err := s.sensorRepository.GetByID(senosorID)
	if err != nil {
		return nil, err
	}

	return sensor, nil
}

func (s *sensorService) Add(model interface{}) error {
	sensor := dto.Sensor{}
	err := mapstructure.Decode(model, &sensor)
	if err != nil {
		return err
	}

	return s.sensorRepository.Add(sensor.Name, sensor.Description, sensor.DeviceId, sensor.SensorGroups, sensor.Unit)
}

func (s *sensorService) Update(model interface{}) error {
	sensor := dto.Sensor{}
	err := mapstructure.Decode(model, &sensor)
	if err != nil {
		return err
	}

	return s.sensorRepository.Update(sensor.Name, sensor.Description, sensor.Unit, sensor.SensorGroups, sensor.ID)
}

func (s *sensorService) Delete(ID string) (interface{}, error) {
	return s.sensorRepository.Delete(ID)
}
