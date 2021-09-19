package sensor

import "fmt"

// Diveces models
type Diveces struct {
	Devices []Device `json:"devices" yaml:"devices"`
}

// GetDeviceID gets device id.
func (d *Diveces) getDeviceID() (string, error) {
	sensorLogger.Info("Getting device id...")
	if devices.Devices[0].ID == "" {
		return "", fmt.Errorf("there is not available device")
	}

	return devices.Devices[0].ID, nil
}

func (d *Diveces) getDeviceSensorsByGroup(sensorGroup string) ([]sensor, error) {
	sensorLogger.Info("Getting device %q sensors by %q group...", d.Devices[0].Name, sensorGroup)
	var sensorsIds []sensor

	for _, currentSensor := range d.Devices[0].Sensors {

		if currentSensor.SensorGroups[0] == sensorGroup {

			_, err := currentSensor.getSensorIDAccordingToSensorName(currentSensor.Name, currentSensor.ID)
			if err != nil {
				sensorLogger.Error(err)
				return nil, fmt.Errorf("failed to get device sensors by %q group: %w", sensorGroup, err)
			}

			sensorsIds = append(sensorsIds, currentSensor)
		}
	}

	return sensorsIds, nil
}
