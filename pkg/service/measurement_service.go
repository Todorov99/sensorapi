package service

import (
	"context"
	"fmt"
	"time"

	sensorcmd "github.com/Todorov99/sensorcli/cmd"
	"github.com/Todorov99/sensorcli/pkg/sensor"
	"github.com/Todorov99/server/pkg/global"
	"github.com/Todorov99/server/pkg/models"
	"github.com/Todorov99/server/pkg/repository"
)

type MeasurementService interface {
	Monitor(ctx context.Context, duration string, sensorGroup []string, err chan error, response chan interface{}, done chan bool)
}

type measurementService struct {
	valueCfg   models.ValueCfg
	repository repository.Repository
}

func NewMeasurementService(valueCfg models.ValueCfg) MeasurementService {
	return &measurementService{
		valueCfg:   valueCfg,
		repository: repository.CreateMeasurementRepository(),
	}
}

func (m measurementService) Monitor(ctx context.Context, duration string, sensorGroup []string, errChan chan error, response chan interface{}, done chan bool) {
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

			metric, err := m.scanMetrics(measurements, true)
			if err != nil {
				errChan <- err
				response <- metric
				return
			}

			response <- measurements
		}
	}
}

func (m measurementService) scanMetrics(metrics []sensor.Measurment, addToDb bool) (interface{}, error) {

	for _, metr := range metrics {
		if addToDb {
			err := m.repository.Add(metr.MeasuredAt.Format(time.RFC3339), metr.Value, metr.SensorID, metr.DeviceID)
			if err != nil {
				return nil, err
			}
		}

		switch metr.SensorID {
		case global.TempSensor:
			if metr.Value > m.valueCfg.TempMax {
				return metr, fmt.Errorf("overheating warning")
			}
			continue
		case global.FrequencySensor:
			continue
		case global.UsageSensor:
			if metr.Value > m.valueCfg.UsageMax {
				return metr, fmt.Errorf("usage problem")
			}
			continue
		case global.MemoryAvailable:
			if metr.Value > m.valueCfg.MemAvailableMax {
				return metr, fmt.Errorf("mem available problem")
			}
			continue
		case global.MemoryUsedParcent:
			continue
		case global.MemoryUsed:
			continue
		case global.CoresSensor:
			continue
		case global.TotalMemory:
			continue
		}
	}
	return nil, nil
}