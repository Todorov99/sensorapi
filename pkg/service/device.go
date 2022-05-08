package service

import (
	"archive/zip"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Todorov99/sensorapi/pkg/dto"
	"github.com/Todorov99/sensorapi/pkg/entity"
	"github.com/Todorov99/sensorapi/pkg/global"
	"github.com/Todorov99/sensorapi/pkg/repository"
	"github.com/Todorov99/sensorapi/pkg/server/config"
	"github.com/Todorov99/sensorcli/pkg/logger"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type DeviceService interface {
	GetAll(ctx context.Context, userID int) (interface{}, error)
	GetById(ctx context.Context, ID, userID int) (interface{}, error)
	Add(ctx context.Context, model interface{}, userID int) error
	Update(ctx context.Context, model interface{}, userID int) error
	Delete(ctx context.Context, deviceID, userID int) (interface{}, error)
	GenerateDeviceCfg(ctx context.Context, deviceID, userID int, binaryOS string) (string, error)
}

type deviceService struct {
	logger           *logrus.Entry
	deviceRepository repository.DeviceRepository
	sensorRepository repository.SensorRepository
	userRepository   repository.UserRepository
}

func NewDeviceService() DeviceService {
	return &deviceService{
		logger:           logger.NewLogrus("deviceService", os.Stdout),
		deviceRepository: repository.NewDeviceRepository(),
		sensorRepository: repository.NewSensorRepository(),
		userRepository:   repository.NewUserRepository(),
	}
}

func (d *deviceService) GenerateDeviceCfg(ctx context.Context, deviceID, userID int, binaryOS string) (string, error) {
	d.logger.Debug("Generating device cfg...")
	err := os.MkdirAll(global.CliResourceDir, 0777)
	if err != nil {
		return "", err
	}

	err = copyRootCA(global.CliResourceDir)
	if err != nil {
		return "", err
	}

	err = copyBinary(global.CliBinariesDir, binaryOS, global.CliResourceDir)
	if err != nil {
		return "", err
	}

	dd, err := d.GetById(ctx, deviceID, userID)
	if err != nil {
		return "", err
	}

	f, err := os.Create(global.CfgFileName)
	if err != nil {
		return "", err
	}

	defer func() {
		f.Close()
	}()

	device := dto.Device{}
	err = mapstructure.Decode(dd, &device)
	if err != nil {
		return "", err
	}

	deviceBytes, err := yaml.Marshal(device)
	if err != nil {
		return "", err
	}

	_, err = f.Write(deviceBytes)
	if err != nil {
		return "", err
	}

	err = createFileChecksum(global.CfgFileName)
	if err != nil {
		return "", err
	}

	err = zipSource(global.CliResourceDir, global.CliZipCfg)
	if err != nil {
		return "", err
	}

	d.logger.Debug("Device cfg successfully generated")
	return global.CliZipCfg, nil
}

func (d *deviceService) GetAll(ctx context.Context, userID int) (interface{}, error) {
	d.logger.Debug("Getting all devices")

	devices, err := d.deviceRepository.GetAll(ctx, userID)
	if err != nil {
		return nil, err
	}

	allDevices := []dto.Device{}
	err = mapstructure.Decode(devices, &allDevices)
	if err != nil {
		return nil, err
	}

	return allDevices, nil
}

func (d *deviceService) GetById(ctx context.Context, deviceID, userID int) (interface{}, error) {
	entityDevice, err := d.deviceRepository.GetByID(ctx, deviceID, userID)
	if err != nil {
		if errors.Is(err, global.ErrorObjectNotFound) {
			return nil, fmt.Errorf("device with ID: %d does not exist", deviceID)
		}
		return nil, err
	}

	device := dto.Device{}
	err = mapstructure.Decode(entityDevice, &device)
	if err != nil {
		return nil, err
	}

	return device, nil
}

func (d *deviceService) Add(ctx context.Context, model interface{}, userID int) error {
	device := entity.Device{}
	err := mapstructure.Decode(model, &device)
	if err != nil {
		return err
	}

	err = d.ifDeviceExist(ctx, device.Name, userID)
	if errors.Is(err, global.ErrorDeviceWithNameAlreadyExist) {
		return fmt.Errorf("device with name %s already exists", device.Name)
	}

	if err != nil && !errors.Is(err, global.ErrorDeviceWithNameAlreadyExist) {
		return err
	}

	err = d.deviceRepository.Add(ctx, device, userID)
	if err != nil {
		return err
	}

	deviceID, err := d.deviceRepository.GetDeviceIDByName(ctx, device.Name, userID)
	if err != nil {
		return err
	}

	device.ID = deviceID

	sensors, err := d.sensorRepository.GetAll(ctx)
	if err != nil {
		return err
	}

	device.Sensors = sensors.([]entity.Sensor)

	for _, sensor := range device.Sensors {
		err := d.deviceRepository.AddDeviceSensors(ctx, device.ID, sensor.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *deviceService) Update(ctx context.Context, model interface{}, userID int) error {
	device := entity.Device{}
	err := mapstructure.Decode(model, &device)
	if err != nil {
		return err
	}
	d.logger.Debugf("Updating device with ID: %d", device.ID)

	_, err = d.deviceRepository.GetByID(ctx, int(device.ID), userID)
	if err != nil {
		if errors.Is(err, global.ErrorObjectNotFound) {
			return fmt.Errorf("device with id: %d does not exist", device.ID)
		}
		return err
	}

	return d.deviceRepository.Update(ctx, device, userID)
}

func (d *deviceService) Delete(ctx context.Context, deviceID, userID int) (interface{}, error) {
	if deviceID == 1 && userID == 1 {
		return nil, fmt.Errorf("your are not allowed to delete your default device")
	}

	deviceForDelete, err := d.deviceRepository.GetByID(ctx, deviceID, userID)
	if err != nil {
		if errors.Is(err, global.ErrorObjectNotFound) {
			return nil, fmt.Errorf("device with id: %d does not exist", deviceID)
		}
		return nil, err
	}

	err = d.deviceRepository.Delete(ctx, deviceID, userID)
	if err != nil {
		return nil, err
	}

	device := dto.Device{}
	err = mapstructure.Decode(deviceForDelete, &device)
	if err != nil {
		return nil, err
	}

	return device, nil
}

func (d *deviceService) ifDeviceExist(ctx context.Context, deviceName string, userId int) error {
	checkForExistingDevice, err := d.deviceRepository.GetDeviceIDByName(ctx, deviceName, userId)
	if err != nil && !errors.Is(err, global.ErrorObjectNotFound) {
		return err
	}

	if checkForExistingDevice != 0 {
		return global.ErrorDeviceWithNameAlreadyExist
	}

	return nil
}

func zipSource(source, target string) error {
	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := zip.NewWriter(f)
	defer writer.Close()

	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Method = zip.Deflate

		header.Name, err = filepath.Rel(filepath.Dir(source), path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			header.Name += "/"
		}

		headerWriter, err := writer.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(headerWriter, f)
		return err
	})
}

func copyRootCA(dest string) error {
	destCertDir := dest + "/.sec"
	err := os.MkdirAll(destCertDir, 0700)
	if err != nil {
		return err
	}

	return copyFile(global.CertificatesPath+"/"+config.GetTLSCfg().RootCACert, destCertDir+"/"+"rootCACert.pem")
}

func copyBinary(sourceDir, OS, dest string) error {
	dirEntries, err := os.ReadDir(sourceDir)
	if err != nil {
		return err
	}

	var osBasedBinaryFilename string
	for _, f := range dirEntries {
		if strings.HasSuffix(f.Name(), OS) {
			osBasedBinaryFilename = f.Name()
			break
		}
	}

	if osBasedBinaryFilename == "" {
		return fmt.Errorf("invalid OS: %s", OS)
	}

	return copyFile(sourceDir+"/"+osBasedBinaryFilename, dest+"/"+osBasedBinaryFilename)
}

func createFileChecksum(filepath string) error {
	h := sha256.New()
	f, err := os.Open(filepath)
	if err != nil {
		return err
	}

	defer func() {
		f.Close()
	}()

	if _, err := io.Copy(h, f); err != nil {
		return err
	}

	checkSumFile, err := os.Create(global.DeviceCfgChecksum)
	if err != nil {
		return err
	}
	_, err = checkSumFile.WriteString(fmt.Sprintf("%x", h.Sum(nil)))
	if err != nil {
		return err
	}

	return nil
}

func copyFile(source, dest string) error {
	original, err := os.Open(source)
	if err != nil {
		return err
	}
	defer original.Close()

	new, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer new.Close()

	_, err = io.Copy(new, original)
	if err != nil {
		return err
	}

	return nil
}
