package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/Todorov99/sensorcli/pkg/sensor"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
)

type Cpu interface {
	//GetMeasurements is func that gets concrete cpu sensor measurements
	GetMeasurements(ctx context.Context, device interface{}) ([]sensor.Measurment, error)
}

type cpuSensor struct {
	groups map[string]string
}

func NewCpu(sensorGroup map[string]string) Cpu {
	return &cpuSensor{
		groups: sensorGroup,
	}
}

func (c *cpuSensor) GetMeasurements(ctx context.Context, device interface{}) ([]sensor.Measurment, error) {
	d := &sensor.Device{}
	err := mapstructure.Decode(device, d)
	if err != nil {
		return nil, err
	}

	sensor.SetDevice(d)

	return getMultipleSensorsMeasurements(ctx, c.groups)
}

func getSensorMeasurements(ctx context.Context, sensorGroup, sensorSysFile string) ([]sensor.Measurment, error) {
	if sensorGroup == "" {
		return nil, fmt.Errorf("invalid sensor group")
	}

	sensorType, err := sensor.NewSensor(sensorGroup)
	if err != nil {
		return nil, err
	}

	if sensorSysFile != "" {
		sensorType.SetSysInfoFile(sensorSysFile)
	}

	err = sensorType.ValidateUnit()
	if err != nil {
		cmdLogger.Error(err)
		return nil, err
	}

	err = sensorType.ValidateFormat(format)
	if err != nil {
		cmdLogger.Error(err)
		return nil, err
	}

	sensorMeasurements, err := sensorType.GetSensorData(ctx, format)
	if err != nil {
		cmdLogger.Error(err)
		return nil, err
	}

	return sensorMeasurements, nil
}

func loadDeviceConfig() (*sensor.Device, error) {
	cmdLogger.Debugf("Loading device config...")
	fileName, err := filepath.Abs("./device.yaml")
	if err != nil {
		return nil, err
	}

	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	device := &sensor.Device{}
	fileErr := yaml.Unmarshal(yamlFile, device)
	if fileErr != nil {
		return nil, err
	}

	return device, nil
}
