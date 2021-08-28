package sensor

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Todorov99/sensorcli/pkg/util"
	"github.com/hashicorp/go-multierror"
	"github.com/shirou/gopsutil/host"
)

const (
	tempSensor string = "CPU_TEMP"
)

type cpuTempSensor struct {
	cpuTemp  string
	deviceID string
	sensors  []sensor
}

// CreateTempSensor creates instance of temperature sensor.
func CreateTempSensor() ISensor {
	return &cpuTempSensor{}
}

func (tempS *cpuTempSensor) GetSensorData(ctx context.Context, format string) ([]Measurment, error) {
	sensorLogger.Info("Gerring sensor data...")
	cpuTemp, err := getTempMeasurments(ctx, format)
	if err != nil {
		msg := "failed to get temperature measurements: %w"
		sensorLogger.Errorf(msg, err)
		return nil, fmt.Errorf(msg, err)
	}

	return cpuTemp, nil
}

func (tempS *cpuTempSensor) ValidateFormat(format string) error {
	return util.ValidateFormat(format)
}

func (tempS *cpuTempSensor) ValidateUnit() error {
	sensorLogger.Info("Validating temperature sensor units...")
	var err error

	currentDeviceSensors, err := devices.getDeviceSensorsByGroup(tempSensor)
	if err != nil {
		return fmt.Errorf("failed to get current device sensors: %w", err)
	}

	tempS.sensors = currentDeviceSensors

	for _, currentSensor := range tempS.sensors {
		if currentSensor.Unit != "F" && currentSensor.Unit != "C" {
			err = multierror.Append(err, fmt.Errorf("invalid temperature unit %q", currentSensor.Unit))
		}
	}

	return err
}

func getTempMeasurments(ctx context.Context, format string) ([]Measurment, error) {
	sensorLogger.Info("Getting temperature sensor measurements...")

	cpuTempInfo, err := getTempFromSensor(ctx)
	if err != nil {
		return nil, err
	}

	deviceID, err := devices.getDeviceID()
	if err != nil {
		return nil, err
	}

	sensor, err := devices.getDeviceSensorsByGroup(tempSensor)
	if err != nil {
		return nil, err
	}

	cpuTempInfo.sensors = sensor
	cpuTempInfo.deviceID = deviceID

	return newMeasurements(cpuTempInfo), nil
}

func getTempFromSensor(ctx context.Context) (cpuTempSensor, error) {
	sensorLogger.Info("Getting temperature from sensor")
	sensorTeperatureInfo, err := host.SensorsTemperaturesWithContext(ctx)
	if err != nil {
		return cpuTempSensor{}, err
	}

	cpuTemp := sensorTeperatureInfo[0].Temperature
	sensorLogger.Info("Temperature from sensor is successfully got")

	return cpuTempSensor{
		cpuTemp: strconv.FormatFloat(cpuTemp, 'f', 1, 64),
	}, nil
}
