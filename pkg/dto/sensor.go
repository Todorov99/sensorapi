package dto

// Sensor model
type Sensor struct {
	ID           string `json:"id,omitempty" yaml:"id,omitempty"`
	Name         string `json:"name,omitempty" yaml:"name,omitempty"`
	DeviceId     string `json:"deviceId,omitempty" yaml:"deviceId,omitempty"`
	Description  string `json:"description,omitempty" yaml:"description,omitempty"`
	Unit         string `json:"unit,omitempty" yaml:"unit,omitempty"`
	SensorGroups string `json:"sensorGroups,omitempty" yaml:"sensorGroups,omitempty"`
}
