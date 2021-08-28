package sensor

import (
	"time"

	"github.com/Todorov99/sensorcli/pkg/util"
)

// Measurment is struct type that holds inforamation about different sensor metrics
type Measurment struct {
	MeasuredAt time.Time `json:"measuredAt" yaml:"measuredAt"`
	Value      string    `json:"value" yaml:"value"`
	SensorID   string    `json:"sensorId" yaml:"sensorId"`
	DeviceID   string    `json:"deviceId" yaml:"deviceId"`
}

func newMeasurement(value string, sensorID string, deviceID string) Measurment {
	return Measurment{time.Now(), value, sensorID, deviceID}
}

func newMeasurements(info interface{}) []Measurment {
	var m = []Measurment{}
	switch v := info.(type) {
	case cpuUsageSensor:
		for _, s := range v.sensors {
			switch s.Name {
			case cpuCores:
				m = append(m, newMeasurement(v.cpuCores, s.ID, v.deviceID))
			case cpuFrequency:
				m = append(m, newMeasurement(v.cpuFrequency, s.ID, v.deviceID))
			case cpuUsage:
				m = append(m, newMeasurement(v.cpuUsage, s.ID, v.deviceID))
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
				m = append(m, newMeasurement(val, s.ID, v.deviceID))
			case memoryUsed:
				val, err := util.ParseMemoryUsageAccordingToUnit(s.Unit, v.usedMemory)
				if err != nil {
					sensorLogger.Error(err)
				}
				m = append(m, newMeasurement(val, s.ID, v.deviceID))
			case memoryTotal:
				val, err := util.ParseMemoryUsageAccordingToUnit(s.Unit, v.totalMemory)
				if err != nil {
					sensorLogger.Error(err)
				}
				m = append(m, newMeasurement(val, s.ID, v.deviceID))
			case memoryUsedPercent:
				val, err := util.ParseMemoryUsageAccordingToUnit(s.Unit, v.usedPercentMemory)
				if err != nil {
					sensorLogger.Error(err)
				}
				m = append(m, newMeasurement(val, s.ID, v.deviceID))
			}

		}
	case cpuTempSensor:
		for _, s := range v.sensors {
			switch s.Name {
			case cpuTemp:
				val, err := util.ParseTempAccordingToUnit(s.Unit, v.cpuTemp)
				if err != nil {
					sensorLogger.Error(err)
				}
				m = append(m, newMeasurement(val, s.ID, v.deviceID))
			}

		}
	}

	return m
}
