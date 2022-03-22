package dto

// Sensor model
type Sensor struct {
	ID           int32  `json:"id,omitempty" mapstructure:"id,omitempty"`
	Name         string `json:"name,omitempty" mapstructure:"name,omitempty"`
	DeviceId     string `json:"deviceId,omitempty" mapstructure:"deviceId,omitempty"`
	Description  string `json:"description,omitempty" mapstructure:"description,omitempty"`
	Unit         string `json:"unit,omitempty" mapstructure:"unit,omitempty"`
	SensorGroups string `json:"sensorGroups,omitempty" mapstructure:"group_name,omitempty"`
}
