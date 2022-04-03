package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	sensorcmd "github.com/Todorov99/sensorcli/cmd"
	"github.com/Todorov99/sensorcli/pkg/sensor"
	"github.com/Todorov99/sensorcli/pkg/writer"
	"github.com/Todorov99/server/pkg/dto"
	"github.com/Todorov99/server/pkg/entity"
	"github.com/Todorov99/server/pkg/global"
	"github.com/Todorov99/server/pkg/repository"
	"github.com/hashicorp/go-multierror"
	"github.com/mitchellh/mapstructure"
)

type MeasurementService interface {
	Monitor(ctx context.Context, deviceID int, duration, deltaDuration string, sensorGroup map[string]string, valueCfg dto.ValueCfg, generateReport bool) (<-chan bool, error)
	GetMonitorStatus() dto.MonitorStatus
	GetSensorsCorrelationCoefficient(ctx context.Context, deviceID1 string, deviceID2 string, sensorID1 string, sensorID2 string, startTime string, endTime string) (float64, error)
	GetAverageValueOfMeasurements(ctx context.Context, deviceID string, sensorID string, startTime string, endTime string) (string, error)
	GetMeasurementsBetweenTimestamp(ctx context.Context, measurementsBetweeTimestamp dto.MeasurementBetweenTimestamp) ([]dto.Measurement, error)
	AddMeasurements(ctx context.Context, measurement dto.Measurement) error
}

type measurementService struct {
	measurementRepository repository.MeasurementRepository
	deviceRepository      repository.DeviceRepository
	sensorRepository      repository.SensorRepository
}

type monitorState struct {
	startTime             string
	finishedAt            string
	done                  bool
	alreadyStartedProcess bool
	monitorError          error
	reportFile            string
	measurements          []sensor.Measurment
	criticalMeasurements  []sensor.Measurment
}

func NewMeasurementService() MeasurementService {
	return &measurementService{
		measurementRepository: repository.NewMeasurementRepository(),
		sensorRepository:      repository.NewSensorRepository(),
		deviceRepository:      repository.NewDeviceRepository(),
	}
}

var monState monitorState

func (m measurementService) Monitor(ctx context.Context, deviceID int, duration, deltaDuration string, sensorGroup map[string]string, valueCfg dto.ValueCfg, generateReport bool) (<-chan bool, error) {
	done := make(chan bool)
	monState = monitorState{}
	monState.alreadyStartedProcess = true
	monState.startTime = time.Now().Format(time.RFC3339)

	device, err := m.deviceRepository.GetByID(ctx, deviceID)
	if err != nil {
		if errors.Is(err, global.ErrorObjectNotFound) {
			return nil, fmt.Errorf("device with ID: %d does not exist", deviceID)
		}
		return nil, err
	}

	d, err := time.ParseDuration(duration)
	if err != nil {
		return nil, err
	}

	monitorDuration := time.After(d)
	cpu := sensorcmd.NewCpu(sensorGroup)
	reportFilename := "measurement" + time.Now().Format("2006-01-02-15:04:05") + ".xlsx"
	reportWriter := writer.New(reportFilename)
	dDuration, err := time.ParseDuration(deltaDuration)
	if err != nil {
		return nil, err
	}

	t := time.NewTicker(dDuration)

	go func() {
		defer func() {
			close(done)
		}()

		for {
			select {
			case <-monitorDuration:
				if generateReport {
					monState.reportFile = reportFilename
				}

				done <- true
				monState.done = true
				monState.finishedAt = time.Now().Format(time.RFC3339)
				return
			case <-ctx.Done():
				monState.monitorError = ctx.Err()
				monState.done = true
				monState.finishedAt = time.Now().Format(time.RFC3339)
				return
			case <-t.C:
				serviceLogger.Debug("Getting sensor measurements...")
				measurements, err := cpu.GetMeasurements(ctx, device)
				if err != nil {
					monState.monitorError = err
					monState.done = true
					monState.finishedAt = time.Now().Format(time.RFC3339)
					done <- true
					return
				}

				metric, err := m.scanMetrics(ctx, measurements, valueCfg, true)
				if err != nil {
					monState.monitorError = err
					monState.done = true
					monState.finishedAt = time.Now().Format(time.RFC3339)
					monState.criticalMeasurements = metric

					done <- true
					return
				}

				if generateReport {
					go func() {
						err := reportWriter.WritoToXslx(measurements)
						if err != nil {
							monState.monitorError = err
							monState.done = true

							done <- true
							return
						}
					}()
				}

				monState.measurements = append(monState.measurements, measurements...)
			}
		}
	}()

	return done, nil
}

func (m measurementService) GetMonitorStatus() dto.MonitorStatus {
	if !monState.alreadyStartedProcess {
		return dto.MonitorStatus{
			Status: "Monitor process haven't been started yet",
		}
	}

	if !monState.done {
		return dto.MonitorStatus{
			StartTime:    monState.startTime,
			FinishedAt:   monState.finishedAt,
			Status:       "In progress",
			Measurements: monState.measurements,
		}
	}

	if monState.done && monState.monitorError != nil {
		return dto.MonitorStatus{
			StartTime:           monState.startTime,
			Status:              "Finished with error",
			FinishedAt:          monState.finishedAt,
			CriticalMeasurement: monState.criticalMeasurements,
			Measurements:        monState.measurements,
			Error:               monState.monitorError.Error(),
		}
	}

	return dto.MonitorStatus{
		StartTime:    monState.startTime,
		FinishedAt:   monState.finishedAt,
		Status:       "Monitoring finished successfully",
		Measurements: monState.measurements,
		ReportFile:   monState.reportFile,
	}
}

func (m measurementService) GetAverageValueOfMeasurements(ctx context.Context, deviceID string, sensorID string, startTime string, endTime string) (string, error) {
	err := m.ifDeviceAndSensorExists(ctx, deviceID, sensorID)
	if err != nil {
		return "", err
	}

	averageValue, err := m.measurementRepository.GetMeasurementsAverageValueBetweenTimestampByDeviceIDAndSensorID(ctx, startTime, endTime, deviceID, sensorID)
	if err != nil {
		return "", err
	}

	return averageValue, nil
}

func (m measurementService) GetMeasurementsBetweenTimestamp(ctx context.Context, measurementsBetweeTimestamp dto.MeasurementBetweenTimestamp) ([]dto.Measurement, error) {
	err := m.ifDeviceAndSensorExists(ctx, measurementsBetweeTimestamp.DeviceID, measurementsBetweeTimestamp.SensorID)
	if err != nil {
		return nil, err
	}

	timestampMeasurements, err := m.measurementRepository.
		GetMeasurementsBetweenTimestampByDeviceIDBySensorID(
			ctx,
			measurementsBetweeTimestamp.StartTime,
			measurementsBetweeTimestamp.EndTime,
			measurementsBetweeTimestamp.DeviceID,
			measurementsBetweeTimestamp.SensorID,
		)
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

func (m measurementService) AddMeasurements(ctx context.Context, measurement dto.Measurement) error {
	measurementEntity := entity.Measurement{}
	err := mapstructure.Decode(measurement, &measurementEntity)
	if err != nil {
		return err
	}
	return m.measurementRepository.Add(ctx, measurementEntity)
}

func (m measurementService) scanMetrics(ctx context.Context, metrics []sensor.Measurment, valueCfg dto.ValueCfg, addToDb bool) ([]sensor.Measurment, error) {
	criticalMeasurements := []sensor.Measurment{}
	var merr error
	for _, metr := range metrics {
		if addToDb {
			measurementEntity := entity.Measurement{
				MeasuredAt: metr.MeasuredAt.Format(time.RFC3339),
				Value:      metr.Value,
				SensorID:   metr.SensorID,
				DeviceID:   metr.DeviceID,
			}

			err := m.measurementRepository.Add(ctx, measurementEntity)
			if err != nil {
				return nil, err
			}
		}

		switch metr.SensorID {
		case global.TempSensor:
			if metr.Value > valueCfg.TempMax {
				criticalMeasurements = append(criticalMeasurements, metr)
				merr = multierror.Append(merr, fmt.Errorf("cpu temperature: %q is over than the expected maximum value of %q", metr.Value, valueCfg.TempMax))
			}
			continue
		case global.FrequencySensor:
			if metr.Value > valueCfg.CPUFrequencyMax {
				criticalMeasurements = append(criticalMeasurements, metr)
				merr = multierror.Append(merr, fmt.Errorf("cpu frequency: %q is over than the expected maximum value of %q", metr.Value, valueCfg.CPUFrequencyMax))
			}
			continue
		case global.UsageSensor:
			if metr.Value > valueCfg.CPUUsageMax {
				criticalMeasurements = append(criticalMeasurements, metr)
				merr = multierror.Append(merr, fmt.Errorf("cpu usage: %q is over than the expected maximum value of %q", metr.Value, valueCfg.CPUUsageMax))
			}
			continue
		case global.MemoryAvailable:
			if metr.Value > valueCfg.MemAvailableMax {
				criticalMeasurements = append(criticalMeasurements, metr)
				merr = multierror.Append(merr, fmt.Errorf("the available memory: %q is over than the expected maximum value of %q", metr.Value, valueCfg.MemAvailableMax))
			}
			continue
		case global.MemoryUsed:
			if metr.Value > valueCfg.MemUsedMax {
				criticalMeasurements = append(criticalMeasurements, metr)
				merr = multierror.Append(merr, fmt.Errorf("the used memory: %q is over than the expected maximum value of %q", metr.Value, valueCfg.MemUsedMax))
			}
			continue
		case global.MemoryUsedParcent:
			if metr.Value > valueCfg.MemUsedPercent {
				criticalMeasurements = append(criticalMeasurements, metr)
				merr = multierror.Append(merr, fmt.Errorf("the used memory percentage: %q is over than the expected maximum value of %q", metr.Value, valueCfg.MemUsedPercent))
			}
			continue
		case global.CoresSensor:
			continue
		case global.TotalMemory:
			continue
		}
	}
	return criticalMeasurements, merr
}

// GetSensorsCorrelationCoefficient gets Pearson's correlation coefficient between two sensors.
func (m measurementService) GetSensorsCorrelationCoefficient(ctx context.Context, deviceID1, deviceID2, sensorID1, sensorID2, startTime, endTime string) (float64, error) {
	serviceLogger.Info("Getting correlation coficient...")
	err := m.ifDeviceAndSensorExists(context.Background(), deviceID1, sensorID1)
	if err != nil {
		return 0, err
	}
	serviceLogger.Infof("Getting values for deviceID: %s and sensorID %s...", deviceID1, sensorID1)
	firstSensorValues, err := m.measurementRepository.
		GetMeasurementsValuesBetweenTimestampByDeviceIDAndSensorID(
			ctx, startTime, endTime, deviceID1, sensorID1)
	if err != nil {
		return 0, err
	}

	serviceLogger.Infof("Getting values for deviceID: %s and sensorID %s...", deviceID2, sensorID2)
	secondSensorValues, err := m.measurementRepository.
		GetMeasurementsValuesBetweenTimestampByDeviceIDAndSensorID(
			ctx, startTime, endTime, deviceID2, sensorID2)
	if err != nil {
		return 0, err
	}

	serviceLogger.Info("Getting the count of values...")
	countOfMeasurements, err := m.measurementRepository.
		CountMeasurementsBetweenTimestampByDeviceIDBySensorID(
			ctx, startTime, endTime, deviceID1, sensorID1,
		)
	if err != nil {
		return 0, err
	}

	return correlationCoefficient(firstSensorValues, secondSensorValues, countOfMeasurements), nil
}

func correlationCoefficient(firstSensorValues []interface{}, secondSensorValues []interface{}, valueCount float64) float64 {

	sumFirstSensor := 0.0
	sumSecondSensor := 0.0
	sumBothSensorValues := 0.0
	squareSumFirstSensor := 0.0
	squareSumSecondSensor := 0.0

	for i := 0; i < int(valueCount)-1; i++ {

		if i == len(firstSensorValues) || i == len(secondSensorValues) {
			break
		}

		sumFirstSensor = sumFirstSensor + firstSensorValues[i].(float64)

		sumSecondSensor = sumSecondSensor + secondSensorValues[i].(float64)

		sumBothSensorValues = sumBothSensorValues + firstSensorValues[i].(float64)*secondSensorValues[i].(float64)

		squareSumFirstSensor = squareSumFirstSensor + firstSensorValues[i].(float64)*firstSensorValues[i].(float64)
		squareSumSecondSensor = squareSumSecondSensor + secondSensorValues[i].(float64)*secondSensorValues[i].(float64)
	}

	return float64((valueCount*sumBothSensorValues - sumFirstSensor*sumSecondSensor)) /
		(math.Sqrt(float64((valueCount*squareSumFirstSensor - sumFirstSensor*sumFirstSensor) *
			(valueCount*squareSumSecondSensor - sumSecondSensor*sumSecondSensor))))
}

func (m *measurementService) ifDeviceAndSensorExists(ctx context.Context, deviceID, sensorID string) error {
	dID, err := strconv.Atoi(deviceID)
	if err != nil {
		return err
	}
	_, err = m.deviceRepository.GetDeviceNameByID(ctx, dID)
	if err != nil {
		if errors.Is(err, global.ErrorObjectNotFound) {
			return fmt.Errorf("device with id: %d does not exist", dID)
		}

		return err
	}

	sID, err := strconv.Atoi(sensorID)
	if err != nil {
		return err
	}

	_, err = m.sensorRepository.GetByID(context.Background(), sID)
	if err != nil {
		if errors.Is(err, global.ErrorObjectNotFound) {
			return fmt.Errorf("sensor with id: %d does not exist", sID)
		}

		return err
	}

	return nil
}
