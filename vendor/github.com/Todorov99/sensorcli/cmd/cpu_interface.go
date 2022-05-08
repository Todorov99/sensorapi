package cmd

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
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

func loadDeviceConfig(cfgDirPath string) (*sensor.Device, error) {
	cmdLogger.Debugf("Loading device config...")
	cfgFileName := ""

	if cfgDirPath == "" && webHook == "" {
		fileName, err := filepath.Abs("./device_cfg.yaml")
		if err != nil {
			return nil, err
		}
		cfgFileName = fileName
	} else if webHook != "" && cfgDirPath != "" {
		cfgFileName = cfgDirPath + "/device_cfg.yaml"
		err := verifyChecksum(cfgFileName, cfgDirPath+"/.checksum")
		if err != nil {
			return nil, err
		}
	} else if cfgDirPath != "" {
		cfgFileName = cfgDirPath + "/device_cfg.yaml"
	}

	yamlFile, err := ioutil.ReadFile(cfgFileName)
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

func verifyChecksum(cfgFilePath, checkSumPath string) error {
	checkSumHashBytes, err := os.ReadFile(checkSumPath)
	if err != nil {
		return fmt.Errorf("config file has been changes. Get new config from the API and DO NOT change it")
	}

	f, err := os.Open(cfgFilePath)
	if err != nil {
		return err
	}

	defer func() {
		f.Close()
	}()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}

	cfgFileHash := fmt.Sprintf("%x", h.Sum(nil))
	if string(checkSumHashBytes) != cfgFileHash {
		return fmt.Errorf("device_cfg file has been changed")
	}

	return nil
}
