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
	cpuUsage     string
	cpuCores     string
	cpuFrequency string
	deviceID     string
	sensors      []sensor
}

// CreateUsageSensor creates instance of usage sensor.
func CreateUsageSensor() ISensor {
	return &cpuUsageSensor{}
}

func (usageS *cpuUsageSensor) GetSensorData(ctx context.Context, format string) ([]Measurment, error) {
	return getUsageMeasurements(ctx, format)
}

func (usageS *cpuUsageSensor) ValidateFormat(format string) error {
	return util.ValidateFormat(format)
}

func (usageS *cpuUsageSensor) ValidateUnit() error {
	sensorLogger.Info("Validating usage sensor units...")
	var err error

	currentDeviceSensors, err := devices.getDeviceSensorsByGroup(usageSensor)
	if err != nil {
		return fmt.Errorf("failed to get current device sensors: %w", err)
	}

	usageS.sensors = currentDeviceSensors

	for _, currentSensor := range usageS.sensors {
		if currentSensor.Unit != "GHz" &&
			currentSensor.Unit != "%" &&
			currentSensor.Unit != "count" &&
			currentSensor.Unit != "Hz" {
			err = multierror.Append(err, fmt.Errorf("invalid temperature unit %q", currentSensor.Unit))
		}
	}

	return err
}

func (usageS *cpuUsageSensor) SetSysInfoFile(filepath string) {
}

func getUsageMeasurements(ctx context.Context, format string) ([]Measurment, error) {
	sensorLogger.Info("Getting usage sensor measurements...")

	deviceID, err := devices.getDeviceID()
	if err != nil {
		return nil, fmt.Errorf("failed to get device id: %w", err)
	}

	sensors, err := devices.getDeviceSensorsByGroup(usageSensor)
	if err != nil {
		return nil, err
	}

	cpuInfo, err := newCPUUsageInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get cpuInfo: %w", err)
	}

	cpuInfo.deviceID = deviceID
	cpuInfo.sensors = sensors

	return newMeasurements(cpuInfo), nil
}

func newCPUUsageInfo(ctx context.Context) (cpuUsageSensor, error) {
	cores, frequency, err := getCPUCoresAndFrequency(ctx)
	if err != nil {
		return cpuUsageSensor{}, err
	}

	usage, err := getUsedPercent(ctx)
	if err != nil {
		return cpuUsageSensor{}, err
	}

	return cpuUsageSensor{
		cpuCores:     cores,
		cpuFrequency: frequency,
		cpuUsage:     usage,
	}, nil
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
