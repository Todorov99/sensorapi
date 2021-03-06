package sensor

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Todorov99/sensorcli/pkg/util"
	"github.com/hashicorp/go-multierror"
	"github.com/shirou/gopsutil/v3/host"
)

const (
	tempSensor string = "CPU_TEMP"
)

type cpuTempSensor struct {
	cpuTempUnit     string
	cpuTemp         string
	deviceID        int32
	sensors         []Sensor
	group           string
	thermalFilePath string
}

// CreateTempSensor creates instance of temperature sensor.
func CreateTempSensor() ISensor {
	return &cpuTempSensor{
		cpuTempUnit: "C",
		group:       tempSensor,
	}
}

// GetSensorData gets the temperature sensor data
func (tempS *cpuTempSensor) GetSensorData(ctx context.Context, format string) ([]Measurment, error) {
	sensorLogger.Info("Gerring sensor data...")
	cpuTemp, err := tempS.getTempMeasurments(ctx, format)
	if err != nil {
		msg := "failed to get temperature measurements: %w"
		sensorLogger.Errorf(msg, err)
		return nil, fmt.Errorf(msg, err)
	}

	return cpuTemp, nil
}

// ValidateFormat validates the format to which the temperature should be parsed
func (tempS *cpuTempSensor) ValidateFormat(format string) error {
	return util.ValidateFormat(format)
}

// SetSysInfoFile sets the sys thermal info file from where the temperature
// could be measured in case any drivers are not installed on the sytem
func (tempS *cpuTempSensor) SetSysInfoFile(filepath string) {
	tempS.thermalFilePath = filepath
}

func (tempS *cpuTempSensor) ValidateUnit() error {
	sensorLogger.Info("Validating temperature sensor units...")
	var merr error

	currentDeviceSensors, err := device.GetDeviceSensorsByGroup(tempSensor)
	if err != nil {
		return fmt.Errorf("failed to get current device sensors: %w", err)
	}

	tempS.sensors = currentDeviceSensors

	for _, currentSensor := range tempS.sensors {
		if currentSensor.Unit != tempS.cpuTempUnit {
			merr = multierror.Append(err, fmt.Errorf("invalid temperature unit %q", currentSensor.Unit))
		}

		if currentSensor.SensorGroups != tempS.group {
			merr = multierror.Append(err, fmt.Errorf("invalid temperature sensor group %q", currentSensor.SensorGroups))
		}
	}

	return merr
}

func (tempS *cpuTempSensor) getTempMeasurments(ctx context.Context, format string, filePath ...string) ([]Measurment, error) {
	sensorLogger.Info("Getting temperature sensor measurements...")
	cpuTempInfo, err := tempS.getTempFromSensor(ctx)
	if err != nil {
		return nil, err
	}

	deviceID, err := device.GetDeviceID()
	if err != nil {
		return nil, err
	}

	sensor, err := device.GetDeviceSensorsByGroup(tempSensor)
	if err != nil {
		return nil, err
	}

	cpuTempInfo.sensors = sensor
	cpuTempInfo.deviceID = deviceID

	return newMeasurements(cpuTempInfo)
}

func (tempS cpuTempSensor) getTempFromSensor(ctx context.Context) (cpuTempSensor, error) {
	sensorLogger.Info("Getting temperature from sensor")

	if tempS.thermalFilePath != "" {
		temp, err := readSysFile(tempS.thermalFilePath)
		if err != nil {
			return tempS, err
		}
		tempS.cpuTemp = temp
		return tempS, nil
	}

	sensorTeperatureInfo, err := host.SensorsTemperaturesWithContext(ctx)
	if err != nil {
		return tempS, err
	}

	cpuTemp := sensorTeperatureInfo[0].Temperature
	sensorLogger.Info("Temperature from sensor is successfully got")

	tempS.cpuTemp = strconv.FormatFloat(cpuTemp, 'f', 1, 64)
	return tempS, nil
}

func readSysFile(filePath string) (string, error) {
	absolutePath, err := filepath.Abs(filePath)
	if err != nil {
		return "", err
	}

	b, err := os.ReadFile(absolutePath)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
