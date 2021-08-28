package util

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Todorov99/sensorcli/pkg/logger"
	yaml "gopkg.in/yaml.v2"
)

var utilLogger logger.Logger = logger.NewLogger("./util")

const (
	unitTypeFahrenheit string = "F"
	unitTypeCelsius    string = "C"
	celsiusName        string = "Celsius"
	fahrenheitName     string = "Fahrenheit"

	gigaBytes string = "GigaBytes"
	bytes     string = "Bytes"
	kiloBytes string = "KiloBytes"
	megaBytes string = "MegaBytes"

	jsonType string = "JSON"
	yamlType string = "YAML"
)

// ParseTempAccordingToUnit returns correct parsed temperature.
func ParseTempAccordingToUnit(unit string, temperature string) (string, error) {
	val, err := strconv.ParseFloat(temperature, 64)
	if err != nil {
		utilLogger.Error(err)
		return "", err
	}

	switch unit {
	case unitTypeFahrenheit:
		utilLogger.Info("Parsing temperatur to fahrenheit unit")
		return strconv.FormatFloat(fahrenheit(val), 'f', 1, 64), nil
	case unitTypeCelsius:
		utilLogger.Info("Parsing temperatur to celsius unit")
		return strconv.FormatFloat(val, 'f', 1, 64), nil
	}

	return "", fmt.Errorf("invalid unit")
}

// ParseTempAccordingToUnit returns correct parsed temperature.
func ParseMemoryUsageAccordingToUnit(unit string, value string) (string, error) {
	val, err := strconv.ParseFloat(value, 64)
	if err != nil {
		utilLogger.Error(err)
		return "", err
	}

	switch unit {
	case gigaBytes:
		return strconv.FormatFloat(gigaByte(val), 'f', 1, 64), nil
	case megaBytes:
		return strconv.FormatFloat(megabytes(val), 'f', 1, 64), nil
	case kiloBytes:
		return strconv.FormatFloat(kilobytes(val), 'f', 1, 64), nil
	case bytes:
		return strconv.FormatFloat(val, 'f', 1, 64), nil
	case "%":
		return strconv.FormatFloat(val, 'f', 1, 64), nil

	}

	return "", fmt.Errorf("invalid unit")
}

// ParseDataAccordingToFormat returns corect parsed data.
func ParseDataAccordingToFormat(format string, data interface{}) string {
	if format == yamlType {
		utilLogger.Info("Parsing data to yaml")
		return yamlParser(data)
	}

	utilLogger.Info("Parsing data to json")
	return jsonParser(data)
}

func kilobytes(bytes float64) float64 {
	return bytes / 1024
}

func megabytes(bytes float64) float64 {
	return bytes / 1024 / 1024
}

func gigaByte(bytes float64) float64 {
	return bytes / 1024 / 1024 / 1024
}
func fahrenheit(temperature float64) float64 {
	return temperature*(9/5) + 32
}

func jsonParser(data interface{}) string {
	jsonData, _ := json.Marshal(data)
	return string(jsonData)
}

func yamlParser(data interface{}) string {
	yamlData, _ := yaml.Marshal(data)
	return string(yamlData)
}
