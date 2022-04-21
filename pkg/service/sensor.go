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

type SensorService interface {
	GetAll(ctx context.Context) ([]dto.Sensor, error)
	GetById(ctx context.Context, ID int) (dto.Sensor, error)
	Update(ctx context.Context, model interface{}) error
}

type sensorService struct {
	sensorRepository repository.SensorRepository
	userRepository   repository.UserRepository
}

func NewSensorService() SensorService {
	return &sensorService{
		sensorRepository: repository.NewSensorRepository(),
		userRepository:   repository.NewUserRepository(),
	}
}

func (s *sensorService) GetAll(ctx context.Context) ([]dto.Sensor, error) {
	sensors, err := s.sensorRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	sensorsDto := []dto.Sensor{}
	err = mapstructure.Decode(sensors, &sensorsDto)
	if err != nil {
		return nil, err
	}

	return sensorsDto, nil
}

func (s *sensorService) GetById(ctx context.Context, ID int) (dto.Sensor, error) {
	sensor, err := s.sensorRepository.GetByID(ctx, ID)
	if err != nil {
		if errors.Is(err, global.ErrorObjectNotFound) {
			return dto.Sensor{}, fmt.Errorf("sensor with ID: %d does not exist", ID)
		}
		return dto.Sensor{}, err
	}

	sensorsDto := dto.Sensor{}
	err = mapstructure.Decode(sensor, &sensorsDto)
	if err != nil {
		return dto.Sensor{}, err
	}

	return sensorsDto, nil
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
