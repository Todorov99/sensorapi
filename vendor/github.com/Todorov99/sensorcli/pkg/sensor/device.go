package sensor

import "fmt"

// Device model
type Device struct {
	ID          int32    `json:"id" yaml:"id" mapstructure:"id,omitempty"`
	Name        string   `json:"name" yaml:"name" mapstructure:"name,omitempty"`
	Description string   `json:"description" yaml:"description" mapstructure:"description,omitempty"`
	Sensors     []Sensor `json:"sensors" yaml:"sensors" mapstructure:"sensors,omitempty"`
}

// GetDeviceID gets device id.
func (d *Device) GetDeviceID() (int32, error) {
	sensorLogger.Info("Getting device id...")
	if device.ID == 0 {
		return 0, fmt.Errorf("there is not available device")
	}

	return device.ID, nil
}

func (d *Device) GetDeviceSensorsByGroup(sensorGroup string) ([]Sensor, error) {
	sensorLogger.Info("Getting device %q sensors by %q group...", d.Name, sensorGroup)
	var sensors []Sensor

	for _, currentSensor := range d.Sensors {
		if currentSensor.SensorGroups == sensorGroup {
			sensors = append(sensors, currentSensor)
		}
	}

	return sensors, nil
}
