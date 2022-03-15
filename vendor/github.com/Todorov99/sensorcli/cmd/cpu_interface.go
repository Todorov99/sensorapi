package cmd

import (
	"context"
	"fmt"

	"github.com/Todorov99/sensorcli/pkg/sensor"
)

type Cpu interface {
	//GetMeasurements is func that gets concrete cpu sensor measurements
	GetMeasurements(ctx context.Context) ([]sensor.Measurment, error)
}

type cpuSensor struct {
	groups map[string]string
}

func NewCpu(sensorGroup map[string]string) Cpu {
	return &cpuSensor{
		groups: sensorGroup,
	}
}

func (c *cpuSensor) GetMeasurements(ctx context.Context) ([]sensor.Measurment, error) {
	return getMultipleSensorsMeasurements(ctx, c.groups)
}

func SetSensorGroupSysFile(filePath, sensorGroup string) {

}

func getSensorMeasurements(ctx context.Context, sensorGroup, sensorSysFile string) ([]sensor.Measurment, error) {
	if sensorGroup == "" {
		cmdLogger.Errorf("invalid sensor group")
		return nil, fmt.Errorf("invalid sensor group")
	}

	sensorType, err := sensor.NewSensor(sensorGroup)
	if err != nil {
		cmdLogger.Error(err)
		return nil, err
	}

	if sensorSysFile != "" {
		sensorType.SetSysInfoFile(sensorSysFile)
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
