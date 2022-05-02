package sensor

import (
	"context"
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/Todorov99/sensorcli/pkg/util"
	"github.com/hashicorp/go-multierror"
	"github.com/shirou/gopsutil/v3/cpu"
)

const (
	usageSensor string = "CPU_USAGE"
)

type cpuUsageSensor struct {
	cpuUsage         string
	cpuCores         string
	cpuFrequency     string
	cpuUsageUnit     string
	cpuCoresUnit     string
	cpuFrequencyUnit string
	deviceID         int32
	group            string
	sensors          []Sensor
}

// CreateUsageSensor creates instance of usage sensor.
func CreateUsageSensor() ISensor {
	return &cpuUsageSensor{
		cpuUsageUnit:     "%",
		cpuCoresUnit:     "count",
		cpuFrequencyUnit: "GHz",
		group:            usageSensor,
	}
}

func (usageS *cpuUsageSensor) GetSensorData(ctx context.Context, format string) ([]Measurment, error) {
	return usageS.getUsageMeasurements(ctx, format)
}

func (usageS *cpuUsageSensor) ValidateFormat(format string) error {
	return util.ValidateFormat(format)
}

func (usageS *cpuUsageSensor) ValidateUnit() error {
	sensorLogger.Info("Validating usage sensor units...")
	var merr error

	currentDeviceSensors, err := device.GetDeviceSensorsByGroup(usageSensor)
	if err != nil {
		return fmt.Errorf("failed to get current device sensors: %w", err)
	}

	usageS.sensors = currentDeviceSensors

	for _, currentSensor := range usageS.sensors {
		if currentSensor.Unit != usageS.cpuCoresUnit &&
			currentSensor.Unit != usageS.cpuFrequencyUnit &&
			currentSensor.Unit != usageS.cpuUsageUnit {
			merr = multierror.Append(err, fmt.Errorf("invalid cpu usage unit %q", currentSensor.Unit))
		}

		if currentSensor.SensorGroups != usageS.group {
			merr = multierror.Append(err, fmt.Errorf("invalid usage sensor group %q", currentSensor.SensorGroups))
		}
	}

	return merr
}

func (usageS *cpuUsageSensor) SetSysInfoFile(filepath string) {
}

func (usageS cpuUsageSensor) getUsageMeasurements(ctx context.Context, format string) ([]Measurment, error) {
	sensorLogger.Info("Getting usage sensor measurements...")

	deviceID, err := device.GetDeviceID()
	if err != nil {
		return nil, fmt.Errorf("failed to get device id: %w", err)
	}

	sensors, err := device.GetDeviceSensorsByGroup(usageSensor)
	if err != nil {
		return nil, err
	}

	cpuInfo, err := usageS.newCPUUsageInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get cpuInfo: %w", err)
	}

	cpuInfo.deviceID = deviceID
	cpuInfo.sensors = sensors

	return newMeasurements(cpuInfo)
}

func (usageS cpuUsageSensor) newCPUUsageInfo(ctx context.Context) (cpuUsageSensor, error) {
	cores, frequency, err := getCPUCoresAndFrequency(ctx)
	if err != nil {
		return cpuUsageSensor{}, err
	}

	usage, err := getUsedPercent(ctx)
	if err != nil {
		return cpuUsageSensor{}, err
	}

	usageS.cpuCores = cores
	usageS.cpuFrequency = frequency
	usageS.cpuUsage = usage
	return usageS, nil
}

func getCPUCoresAndFrequency(ctx context.Context) (string, string, error) {
	var cpuInfo []cpu.InfoStat
	var err error

	if runtime.GOOS == "darwin" && runtime.GOARCH == "arm64" {
		cpuInfo, err = darwinArm64InfoWithContext(ctx)
		if err != nil {
			return "", "", fmt.Errorf("error getting cpu cores and frequency on darwin arm64: %w", err)
		}

		return strconv.FormatInt(int64(cpuInfo[0].Cores), 10), strconv.FormatFloat(cpuInfo[0].Mhz, 'f', 2, 64), nil
	}

	cpuInfo, err = cpu.InfoWithContext(ctx)
	if err != nil {
		return "", "", fmt.Errorf("error with getting cpu cores and frequency: %w", err)
	}

	return strconv.FormatInt(int64(cpuInfo[0].Cores), 10), strconv.FormatFloat(cpuInfo[0].Mhz, 'f', 2, 64), nil
}

func getUsedPercent(ctx context.Context) (string, error) {
	cpuUsedPercentage, err := cpu.PercentWithContext(ctx, time.Second, false)
	if err != nil {
		return "", fmt.Errorf("error in getting used cpu percent")
	}

	usedPercentage := cpuUsedPercentage[0] / 100

	return strconv.FormatFloat(usedPercentage, 'f', 2, 64), nil
}
