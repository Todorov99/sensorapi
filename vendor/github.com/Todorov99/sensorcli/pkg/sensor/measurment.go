package sensor

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Todorov99/sensorcli/pkg/util"
)

const (
	TimeFormat = "2006-01-02-15:04:05"
)

// Measurment is struct type that holds inforamation about different sensor metrics
type Measurment struct {
	MeasuredAt string `json:"measuredAt" yaml:"measuredAt"`
	Value      string `json:"value" yaml:"value"`
	SensorID   string `json:"sensorId" yaml:"sensorId"`
	DeviceID   string `json:"deviceId" yaml:"deviceId"`
	Unit       string `json:"unit" yaml:"unit"`
}

func newMeasurement(value, unit string, sensorID int32, deviceID int32) Measurment {
	return Measurment{
		MeasuredAt: time.Now().Format(TimeFormat),
		Value:      value,
		SensorID:   strconv.FormatInt(int64(sensorID), 10),
		DeviceID:   strconv.FormatInt(int64(deviceID), 10),
		Unit:       unit,
	}
}

func newMeasurements(info interface{}) ([]Measurment, error) {
	var m = []Measurment{}
	switch v := info.(type) {
	case cpuUsageSensor:
		for _, s := range v.sensors {
			switch s.Name {
			case cpuCores:
				m = append(m, newMeasurement(v.cpuCores, v.cpuCoresUnit, s.ID, v.deviceID))
			case cpuFrequency:
				m = append(m, newMeasurement(v.cpuFrequency, v.cpuFrequencyUnit, s.ID, v.deviceID))
			case cpuUsagePecentage:
				m = append(m, newMeasurement(v.cpuUsage, v.cpuUsageUnit, s.ID, v.deviceID))
			default:
				return nil, fmt.Errorf("invalid sensor name: %q", s.Name)
			}
		}
	case cpuMemorySensor:
		for _, s := range v.sensors {
			switch s.Name {
			case memoryAvailable:
				val, err := util.ParseMemoryUsageAccordingToUnit(s.Unit, v.availableMemory)
				if err != nil {
					sensorLogger.Error(err)
				}
				m = append(m, newMeasurement(val, v.availableMemoryUnit, s.ID, v.deviceID))
			case memoryUsed:
				val, err := util.ParseMemoryUsageAccordingToUnit(s.Unit, v.usedMemory)
				if err != nil {
					sensorLogger.Error(err)
				}
				m = append(m, newMeasurement(val, v.usedMemoryUnit, s.ID, v.deviceID))
			case memoryTotal:
				val, err := util.ParseMemoryUsageAccordingToUnit(s.Unit, v.totalMemory)
				if err != nil {
					sensorLogger.Error(err)
				}
				m = append(m, newMeasurement(val, v.totalMemoryUnit, s.ID, v.deviceID))
			case memoryUsedPercent:
				val, err := util.ParseMemoryUsageAccordingToUnit(s.Unit, v.usedPercentMemory)
				if err != nil {
					sensorLogger.Error(err)
				}
				m = append(m, newMeasurement(val, v.usedPercentMemoryUnit, s.ID, v.deviceID))
			default:
				return nil, fmt.Errorf("invalid sensor name: %q", s.Name)
			}

		}
	case cpuTempSensor:
		for _, s := range v.sensors {
			switch s.Name {
			case cpuTemperature:
				val, err := util.ParseTempAccordingToUnit(s.Unit, v.cpuTemp)
				if err != nil {
					sensorLogger.Error(err)
				}
				m = append(m, newMeasurement(val, v.cpuTempUnit, s.ID, v.deviceID))
			default:
				return nil, fmt.Errorf("invalid sensor name: %q", s.Name)
			}
		}
	}

	return m, nil
}
