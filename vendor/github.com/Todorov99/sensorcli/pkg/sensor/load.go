package sensor

import (
	"github.com/Todorov99/sensorcli/pkg/logger"
)

const (
	cpuTemperature    string = "cpuTemperature"
	cpuCores          string = "cpuCores"
	cpuUsagePecentage string = "cpuUsagePercentage"
	cpuFrequency      string = "cpuFrequency"
	memoryTotal       string = "memoryTotal"
	memoryAvailable   string = "memoryAvailable"
	memoryUsed        string = "memoryUsed"
	memoryUsedPercent string = "memoryUsedPercentage"
)

var device *Device

var sensorLogger = logger.NewLogrus("./sensor")

func SetDevice(d *Device) {
	device = d
}
