package sensor

import (
	"github.com/Todorov99/sensorcli/pkg/logger"
)

const (
	cpuTemp           string = "cpuTemp"
	cpuCores          string = "cpuCores"
	cpuUsage          string = "cpuUsage"
	cpuFrequency      string = "cpuFrequency"
	memoryTotal       string = "memoryTotal"
	memoryAvailable   string = "memoryAvailable"
	memoryUsed        string = "memoryUsed"
	memoryUsedPercent string = "memoryUsedPercent"
)

//var modelFilePath string = "./device.yaml"

var device *Device

//var sensorLogger *logrus.Entry
var sensorLogger = logger.NewLogrus("./sensor")

func SetDevice(d *Device) {
	device = d
}
