package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/Todorov99/sensorcli/pkg/sensor"
	"github.com/Todorov99/sensorcli/pkg/util"
	"github.com/Todorov99/sensorcli/pkg/writer"
)

type Cpu interface {
	//GetMeasurements is func that gets concrete cpu sensor measurements
	GetMeasurements(ctx context.Context) ([]sensor.Measurment, error)
}

type cpuSensor struct {
	groups []string
}

func NewCpu(sensorGroup []string) Cpu {
	return &cpuSensor{
		groups: sensorGroup,
	}
}

func (c *cpuSensor) GetMeasurements(ctx context.Context) ([]sensor.Measurment, error) {
	return getMultipleSensorsMeasurements(ctx, c.groups)
}

func getMultipleSensorsMeasurements(ctx context.Context, groups []string) ([]sensor.Measurment, error) {
	var multipleSensorsData []sensor.Measurment

	for _, group := range groups {

		var currentSensorGroupData []sensor.Measurment

		currentSensorGroupData, err := getSensorMeasurements(ctx, group)
		if err != nil {
			return nil, err
		}

		multipleSensorsData = append(multipleSensorsData, currentSensorGroupData...)
	}

	return multipleSensorsData, nil
}

func getSensorMeasurements(ctx context.Context, sensorGroup string) ([]sensor.Measurment, error) {
	if sensorGroup == "" {
		cmdLogger.Errorf("invalid sensor group")
		return nil, fmt.Errorf("invalid sensor group")
	}

	sensorType, err := sensor.NewSensor(sensorGroup)
	if err != nil {
		cmdLogger.Error(err)
		return nil, err
	}

	err = sensorType.ValidateUnit()
	if err != nil {
		cmdLogger.Error(err)
		return nil, err
	}

	err = sensorType.ValidateFormat(format)
	if err != nil {
		cmdLogger.Error(err)
		return nil, err
	}

	sensorMeasurements, err := sensorType.GetSensorData(ctx, format)
	if err != nil {
		cmdLogger.Error(err)
		return nil, err
	}

	return sensorMeasurements, nil
}

func getMeasurementsInDeltaDuration(ctx context.Context, sensorData []sensor.Measurment, deltaDuration time.Duration) error {
	cmdLogger.Info("Getting measurements in delta duration...")

	measurementDuration := time.After(deltaDuration)
	done := make(chan bool)
	sensorsData := sendSensorData(sensorData, done)

	defer func() {
		close(done)
	}()

	for {
		select {
		case data := <-sensorsData:

			if webHook != "" {
				webHookURL(webHook, util.ParseDataAccordingToFormat("JSON", data))
			}

			if file != "" {

				go func() {
					var sensorsData []string
					sensorsData = append(sensorsData, util.ParseDataAccordingToFormat(format, sensorData))
					err := writer.WriteOutputToCSV(sensorsData, file)
					if err != nil {
						cmdLogger.Errorf("error during writing in CSV file: %w", err)
					}
				}()
			}

			fmt.Println(util.ParseDataAccordingToFormat(format, data))
		case <-measurementDuration:
			done <- true
			return nil
		case <-ctx.Done():
			done <- true
			cmdLogger.Error(ctx.Err())
			return ctx.Err()
		}
	}

}

func sendSensorData(sensorsInfo []sensor.Measurment, done chan bool) <-chan sensor.Measurment {
	cmdLogger.Info("Sending sensor data...")

	out := make(chan sensor.Measurment)

	go func() {
		for _, currentSensorInfo := range sensorsInfo {
			out <- currentSensorInfo
		}

		if <-done {
			close(out)
			return
		}
	}()

	return out
}
