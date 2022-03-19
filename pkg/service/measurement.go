package service

import (
	"context"
	"fmt"
	"time"

	sensorcmd "github.com/Todorov99/sensorcli/cmd"
	"github.com/Todorov99/sensorcli/pkg/sensor"
	"github.com/Todorov99/server/pkg/dto"
	"github.com/Todorov99/server/pkg/global"
	"github.com/Todorov99/server/pkg/repository"
	"github.com/mitchellh/mapstructure"
)

type MeasurementService interface {
	Monitor(ctx context.Context, duration string, sensorGroup map[string]string, valueCfg dto.ValueCfg, err chan error, response chan interface{}, done chan bool)
	GetSensorsCorrelationCoefficient(deviceID1 string, deviceID2 string, sensorID1 string, sensorID2 string, startTime string, endTime string) (float64, error)
	GetAverageValueOfMeasurements(deviceID string, sensorID string, startTime string, endTime string) (string, error)
	GetMeasurementsBetweenTimestamp(measurementsBetweeTimestamp dto.MeasurementBetweenTimestamp) ([]dto.Measurement, error)
	AddMeasurements(measurement dto.Measurement) error
}

type measurementService struct {
	repository repository.Repository
}

func NewMeasurementService() MeasurementService {
	return &measurementService{
		repository: repository.CreateMeasurementRepository(),
	}
}

func (m measurementService) Monitor(ctx context.Context, duration string, sensorGroup map[string]string, valueCfg dto.ValueCfg, errChan chan error, response chan interface{}, done chan bool) {
	defer func() {
		close(errChan)
		close(response)
		close(done)
	}()

	d, err := time.ParseDuration(duration)
	if err != nil {
		response <- nil
		errChan <- err
		return
	}

	monitorDuration := time.After(d)
	cpu := sensorcmd.NewCpu(sensorGroup)

	for {
		select {
		case <-monitorDuration:
			errChan <- nil
			done <- true
			return
		case <-ctx.Done():
			errChan <- err
			return
		default:
			measurements, err := cpu.GetMeasurements(ctx)
			if err != nil {
				errChan <- err
				return
			}

			metric, err := m.scanMetrics(measurements, valueCfg, true)
			if err != nil {
				errChan <- err
				response <- metric
				return
			}

			response <- measurements
		}
	}
}

func (m measurementService) GetSensorsCorrelationCoefficient(deviceID1, deviceID2, sensorID1, sensorID2, startTime, endTime string) (float64, error) {
	value, err := repository.GetSensorsCorrelationCoefficient(deviceID1, deviceID2, sensorID1, sensorID2, startTime, endTime)
	if err != nil {
		return 0.0, err
	}

	return value, nil
}

func (m measurementService) GetAverageValueOfMeasurements(deviceID string, sensorID string, startTime string, endTime string) (string, error) {
	averageValue, err := repository.GetAverageValueOfMeasurements(deviceID, sensorID, startTime, endTime)
	if err != nil {
		return "", err
	}

	return averageValue, nil
}

func (m measurementService) GetMeasurementsBetweenTimestamp(measurementsBetweeTimestamp dto.MeasurementBetweenTimestamp) ([]dto.Measurement, error) {
	timestampMeasurements, err := m.repository.GetByID(measurementsBetweeTimestamp.StartTime, measurementsBetweeTimestamp.EndTime, measurementsBetweeTimestamp.DeviceID, measurementsBetweeTimestamp.SensorID)
	if err != nil {
		return nil, err
	}
	measurements := []dto.Measurement{}

	err = mapstructure.Decode(timestampMeasurements, &measurements)
	if err != nil {
		return nil, err
	}

	if len(measurements) == 0 {
		return nil, fmt.Errorf("there are not any measurements in the %q - %q timestamp for sensor with ID: %q and device: %q", measurementsBetweeTimestamp.StartTime, measurementsBetweeTimestamp.EndTime, measurementsBetweeTimestamp.SensorID, measurementsBetweeTimestamp.DeviceID)
	}
	return measurements, nil
}

func (m measurementService) AddMeasurements(measurement dto.Measurement) error {
	return m.repository.Add(measurement.MeasuredAt, measurement.Value,
		measurement.SensorID, measurement.DeviceID)
}

func (m measurementService) scanMetrics(metrics []sensor.Measurment, valueCfg dto.ValueCfg, addToDb bool) (interface{}, error) {
	for _, metr := range metrics {
		if addToDb {
			err := m.repository.Add(metr.MeasuredAt.Format(time.RFC3339), metr.Value, metr.SensorID, metr.DeviceID)
			if err != nil {
				return nil, err
			}
		}

		switch metr.SensorID {
		case global.TempSensor:
			if metr.Value > valueCfg.TempMax {
				return metr, fmt.Errorf("cpu temperature: %q is over than the expected maximum value of %q", metr.Value, valueCfg.TempMax)
			}
			continue
		case global.FrequencySensor:
			if metr.Value > valueCfg.CPUFrequencyMax {
				return metr, fmt.Errorf("cpu frequency: %q is over than the expected maximum value of %q", metr.Value, valueCfg.CPUFrequencyMax)
			}
			continue
		case global.UsageSensor:
			if metr.Value > valueCfg.CPUUsageMax {
				return metr, fmt.Errorf("cpu usage: %q is over than the expected maximum value of %q", metr.Value, valueCfg.CPUUsageMax)
			}
			continue
		case global.MemoryAvailable:
			if metr.Value > valueCfg.MemAvailableMax {
				return metr, fmt.Errorf("the available memory: %q is over than the expected maximum value of %q", metr.Value, valueCfg.MemAvailableMax)
			}
			continue
		case global.MemoryUsed:
			if metr.Value > valueCfg.MemUsedMax {
				return metr, fmt.Errorf("the used memory: %q is over than the expected maximum value of %q", metr.Value, valueCfg.MemUsedMax)
			}
			continue
		case global.MemoryUsedParcent:
			continue
		case global.CoresSensor:
			continue
		case global.TotalMemory:
			continue
		}
	}
	return nil, nil
}
