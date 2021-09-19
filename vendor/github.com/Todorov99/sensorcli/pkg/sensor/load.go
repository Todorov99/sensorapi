package sensor

import (
	"io/ioutil"
	"path/filepath"

	"github.com/Todorov99/sensorcli/pkg/logger"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
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

var modelFilePath string = "./model.yaml"

var devices *Diveces
var sensorLogger *logrus.Entry

func init() {
	sensorLogger = logger.NewLogrus("./sensor")

	fileName, err := filepath.Abs(modelFilePath)
	if err != nil {
		sensorLogger.Panic(err)
	}

	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		sensorLogger.Panic(err)
	}

	fileErr := yaml.Unmarshal(yamlFile, &devices)
	if fileErr != nil {
		sensorLogger.Panic(err)
	}

}

// //ReadFileSystemFile reads server temperature from filesystem file.
func ReadFileSystemFile(fileSystemPath string) (float64, error) {
	fileName, err := filepath.Abs(fileSystemPath)
	if err != nil {
		return 0, err
	}

	fileContent, err := ioutil.ReadFile(fileName)
	if err != nil {
		return 0, err
	}

	var temp float64

	err = yaml.Unmarshal(fileContent, &temp)
	if err != nil {
		return 0, err
	}

	return temp, nil
}
