package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Todorov99/sensorapi/pkg/dto"
	"github.com/Todorov99/sensorapi/pkg/global"
	"github.com/Todorov99/sensorapi/pkg/repository"
	"github.com/mitchellh/mapstructure"
)

type sensorService struct {
	sensorRepository repository.SensorRepository
}

func NewSensorService() IService {
	return &sensorService{
		sensorRepository: repository.NewSensorRepository(),
	}
}

func (s *sensorService) GetAll(ctx context.Context) (interface{}, error) {
	sensors, err := s.sensorRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return sensors, nil
}

func (s *sensorService) GetById(ctx context.Context, ID int) (interface{}, error) {
	sensor, err := s.sensorRepository.GetByID(ctx, ID)
	if err != nil {
		if errors.Is(err, global.ErrorObjectNotFound) {
			return nil, fmt.Errorf("sensor with ID: %d does not exist", ID)
		}
		return nil, err
	}

	return sensor, nil
}

func (s *sensorService) Add(ctx context.Context, model interface{}) error {
	sensor := dto.Sensor{}
	err := mapstructure.Decode(model, &sensor)
	if err != nil {
		return err
	}

	sensorID, err := s.sensorRepository.GetSensorIDByName(ctx, sensor.Name)
	if err != nil {
		return err
	}

	if sensorID != "" {
		return fmt.Errorf("sensor with name: %s already exists", sensor.Name)
	}

	return s.sensorRepository.Add(ctx, sensor)
}

func (s *sensorService) Update(ctx context.Context, model interface{}) error {
	sensor := dto.Sensor{}
	err := mapstructure.Decode(model, &sensor)
	if err != nil {
		return err
	}
	_, err = s.sensorRepository.GetByID(ctx, int(sensor.ID))
	if err != nil {
		if errors.Is(err, global.ErrorObjectNotFound) {
			return fmt.Errorf("sensor with ID: %d does not exist", sensor.ID)
		}
		return err
	}
	return s.sensorRepository.Update(ctx, sensor)
}

func (s *sensorService) Delete(ctx context.Context, ID int) (interface{}, error) {
	sensorForDelete, err := s.sensorRepository.GetByID(ctx, ID)
	if err != nil {
		if errors.Is(err, global.ErrorObjectNotFound) {
			return nil, fmt.Errorf("sensor with ID: %d does not exist", ID)
		}
		return nil, err
	}

	return sensorForDelete, s.sensorRepository.Delete(ctx, ID)
}
