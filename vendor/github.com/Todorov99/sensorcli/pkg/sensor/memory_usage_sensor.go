package sensor

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Todorov99/sensorcli/pkg/util"
	"github.com/hashicorp/go-multierror"
	"github.com/shirou/gopsutil/v3/mem"
)

const (
	memorySensor string = "MEMORY_USAGE"
)

type cpuMemorySensor struct {
	totalMemory       string
	availableMemory   string
	usedMemory        string
	usedPercentMemory string
	deviceID          int32
	sensors           []Sensor
}

// CreateMemorySensor creates instance of memory sensor.
func CreateMemorySensor() ISensor {
	return &cpuMemorySensor{}
}

func (memoryS *cpuMemorySensor) GetSensorData(ctx context.Context, format string) ([]Measurment, error) {
	sensorLogger.Info("Gerring sensor data...")
	memoryUsageData, err := getMemoryUsageData(ctx, format)
	if err != nil {
		msg := "failed to get memory usage data: %w"
		sensorLogger.Errorf(msg, err)
		return nil, fmt.Errorf(msg, err)
	}

	return memoryUsageData, nil
}

func (memoryS *cpuMemorySensor) ValidateFormat(format string) error {
	return util.ValidateFormat(format)
}

func (memoryS *cpuMemorySensor) ValidateUnit() error {
	sensorLogger.Info("Validating memory sensor units...")
	var err error

	currentDeviceSensors, err := device.GetDeviceSensorsByGroup(memorySensor)
	if err != nil {
		return fmt.Errorf("failed to get current device sensors: %w", err)
	}

	memoryS.sensors = currentDeviceSensors

	for _, currentSensor := range memoryS.sensors {
		if currentSensor.Unit != "MegaBytes" &&
			currentSensor.Unit != "%" && currentSensor.Unit != "Bytes" &&
			currentSensor.Unit != "Kilobytes" && currentSensor.Unit != "GigaBytes" {
			err = multierror.Append(err, fmt.Errorf("invalid memory unit: %q", currentSensor.Unit))
		}
	}

	return err
}

func (memoryS *cpuMemorySensor) SetSysInfoFile(filepath string) {
}

func getMemoryUsageData(ctx context.Context, format string) ([]Measurment, error) {
	sensorLogger.Info("Getting memory usage data...")

	deviceID, err := device.GetDeviceID()
	if err != nil {
		msg := "failed to get deviceID: %w"
		sensorLogger.Errorf(msg, err)
		return nil, fmt.Errorf(msg, err)
	}

	sensors, err := device.GetDeviceSensorsByGroup(memorySensor)
	if err != nil {
		msg := "failed to get sensorID: %w"
		sensorLogger.Errorf(msg, err)
		return nil, fmt.Errorf(msg, err)
	}

	memoryUsageValues, err := getMemoryUsageValues(ctx)
	if err != nil {
		msg := "failed to get memory usage data: %w"
		sensorLogger.Errorf(msg, err)
		return nil, fmt.Errorf(msg, err)
	}

	memoryUsageValues.sensors = sensors
	memoryUsageValues.deviceID = deviceID

	return newMeasurements(memoryUsageValues), nil
}

func getMemoryUsageValues(ctx context.Context) (cpuMemorySensor, error) {
	sensorLogger.Info("Getting memory usage data...")

	memory, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		msg := "failed to get virtual memory: %w"
		sensorLogger.Errorf(msg, err)
		return cpuMemorySensor{}, fmt.Errorf(msg, err)
	}

	return cpuMemorySensor{
		totalMemory:       strconv.FormatUint(memory.Total, 10),
		availableMemory:   strconv.FormatUint(memory.Available, 10),
		usedMemory:        strconv.FormatUint(memory.Used, 10),
		usedPercentMemory: strconv.FormatFloat(memory.UsedPercent, 'f', 2, 64),
	}, nil
}
