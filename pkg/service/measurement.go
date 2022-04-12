package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/Todorov99/sensorapi/pkg/dto"
	"github.com/Todorov99/sensorapi/pkg/entity"
	"github.com/Todorov99/sensorapi/pkg/global"
	"github.com/Todorov99/sensorapi/pkg/mailsender"
	"github.com/Todorov99/sensorapi/pkg/repository"
	sensorcmd "github.com/Todorov99/sensorcli/cmd"
	"github.com/Todorov99/sensorcli/pkg/logger"
	"github.com/Todorov99/sensorcli/pkg/sensor"
	"github.com/Todorov99/sensorcli/pkg/writer"
	"github.com/hashicorp/go-multierror"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

const (
	StateInProgress = "In progress"
	StateFinished   = "Finished"
	StateError      = "Error"
)

var monState monitorState

type MeasurementService interface {
	Monitor(ctx context.Context, cfg MonitorCfg) (<-chan bool, error)
	GetMonitorStatus() dto.MonitorStatus
	GetSensorsCorrelationCoefficient(ctx context.Context, deviceID1 string, deviceID2 string, sensorID1 string, sensorID2 string, startTime string, endTime string) (float64, error)
	GetAverageValueOfMeasurements(ctx context.Context, deviceID string, sensorID string, startTime string, endTime string) (string, error)
	GetMeasurementsBetweenTimestamp(ctx context.Context, measurementsBetweeTimestamp dto.MeasurementBetweenTimestamp) ([]dto.Measurement, error)
	AddMeasurements(ctx context.Context, measurement dto.Measurement) error
}

type measurementService struct {
	logger                *logrus.Entry
	measurementRepository repository.MeasurementRepository
	deviceRepository      repository.DeviceRepository
	sensorRepository      repository.SensorRepository
	mailsenderClt         *mailsender.Client
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

type MonitorCfg struct {
	DeviceID           int
	Duration           string
	DeltaDuration      string
	SnsorGroups        map[string]string
	CriticalMetricsCfg dto.ValueCfg
	GenerateReport     bool
	SendReport         bool
}

func NewMeasurementService() MeasurementService {
	return &measurementService{
		logger:                logger.NewLogrus("measurementService", os.Stdout),
		measurementRepository: repository.NewMeasurementRepository(),
		sensorRepository:      repository.NewSensorRepository(),
		deviceRepository:      repository.NewDeviceRepository(),
		mailsenderClt:         mailsender.New(),
	}
}

func (m measurementService) Monitor(ctx context.Context, cfg MonitorCfg) (<-chan bool, error) {
	m.logger.Debug("Starting monitoring...")
	done := make(chan bool)
	startTime := time.Now().Format(global.TimeFormat)
	reportFilename := "measurement_" + startTime + ".xlsx"

	monState = monitorState{}
	monState.alreadyStartedProcess = true
	monState.startTime = startTime

	device, err := m.deviceRepository.GetByID(ctx, cfg.DeviceID)
	if err != nil {
		if errors.Is(err, global.ErrorObjectNotFound) {
			return nil, fmt.Errorf("device with ID: %d does not exist", cfg.DeviceID)
		}
		return nil, err
	}

	d, err := time.ParseDuration(cfg.Duration)
	if err != nil {
		return nil, err
	}

	monitorDuration := time.After(d)
	cpu := sensorcmd.NewCpu(cfg.SnsorGroups)

	reportWriter := writer.New(reportFilename)
	dDuration, err := time.ParseDuration(cfg.DeltaDuration)
	if err != nil {
		return nil, err
	}

	t := time.NewTicker(dDuration)

	reportSender := dto.MailSenderDto{
		Subject: "Measurement report",
		To: []string{
			"todor.mtodorov01@gmail.com",
		},
	}

	go func() {
		defer func() {
			close(done)
		}()

		for {
			select {
			case <-monitorDuration:
				if cfg.GenerateReport {
					monState.reportFile = reportFilename
				}

				if cfg.SendReport {
					reportSender.Body = "Your measurements finished successfully without any critical measurements"
					m.mailsenderClt.SendWithAttachments(ctx, reportSender, []string{reportFilename})
				}

				done <- true
				monState.done = true
				monState.finishedAt = time.Now().Format(global.TimeFormat)
				return
			case <-ctx.Done():
				monState.monitorError = ctx.Err()
				monState.done = true
				monState.finishedAt = time.Now().Format(global.TimeFormat)
				return
			case <-t.C:
				m.logger.Debug("Getting sensor measurements...")
				measurements, err := cpu.GetMeasurements(ctx, device)
				if err != nil {
					monState.monitorError = err
					monState.done = true
					monState.finishedAt = time.Now().Format(global.TimeFormat)
					done <- true
					return
				}

				metric, err := m.scanMetrics(ctx, measurements, cfg.CriticalMetricsCfg, true)
				if err != nil {
					var merr error
					merr = multierror.Append(merr, err)

					monState.monitorError = merr
					monState.done = true
					monState.finishedAt = time.Now().Format(global.TimeFormat)
					monState.criticalMeasurements = metric

					critMetricReportFilename := "critical_metrics_from_" + startTime + ".xlsx"

					criticalMetricReportWriter := writer.New(critMetricReportFilename)
					err := criticalMetricReportWriter.WritoToXslx(metric)
					if err != nil {
						merr = multierror.Append(merr, err)
						monState.monitorError = merr

						done <- true
						return
					}

					reportAttachments := []string{
						critMetricReportFilename,
					}

					if cfg.GenerateReport {
						if _, err := os.Stat(reportFilename); err == nil {
							reportAttachments = append(reportAttachments, reportFilename)
						}
					}

					reportSender.Body = "Critical measurements occurred during your monitor timeframe"

					err = m.mailsenderClt.SendWithAttachments(ctx, reportSender, reportAttachments)
					if err != nil {
						m.logger.Warnf("Failed sending email with critical measurement report: %w", err)
					}

					done <- true
					return
				}

				if cfg.GenerateReport {
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
			Status:       StateInProgress,
			Measurements: monState.measurements,
		}
	}

	if monState.done && monState.monitorError != nil {
		return dto.MonitorStatus{
			StartTime:           monState.startTime,
			Status:              StateError,
			FinishedAt:          monState.finishedAt,
			CriticalMeasurement: monState.criticalMeasurements,
			Measurements:        monState.measurements,
			Error:               monState.monitorError.Error(),
		}
	}

	return dto.MonitorStatus{
		StartTime:    monState.startTime,
		FinishedAt:   monState.finishedAt,
		Status:       StateFinished,
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
				MeasuredAt: metr.MeasuredAt,
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
	m.logger.Info("Getting correlation coficient...")
	err := m.ifDeviceAndSensorExists(context.Background(), deviceID1, sensorID1)
	if err != nil {
		return 0, err
	}
	m.logger.Infof("Getting values for deviceID: %s and sensorID %s...", deviceID1, sensorID1)
	firstSensorValues, err := m.measurementRepository.
		GetMeasurementsValuesBetweenTimestampByDeviceIDAndSensorID(
			ctx, startTime, endTime, deviceID1, sensorID1)
	if err != nil {
		return 0, err
	}

	m.logger.Infof("Getting values for deviceID: %s and sensorID %s...", deviceID2, sensorID2)
	secondSensorValues, err := m.measurementRepository.
		GetMeasurementsValuesBetweenTimestampByDeviceIDAndSensorID(
			ctx, startTime, endTime, deviceID2, sensorID2)
	if err != nil {
		return 0, err
	}

	m.logger.Info("Getting the count of values...")
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
