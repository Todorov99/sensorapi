package dto

// Sensor model
type Sensor struct {
	ID   int32  `json:"id,omitempty" yaml:"id,omitempty"  mapstructure:"id,omitempty"`
	Name string `json:"name,omitempty" yaml:"name,omitempty"  mapstructure:"name,omitempty"`
	//	DeviceId     string `json:"deviceId,omitempty" yaml:"deviceID,omitempty"  mapstructure:"deviceId,omitempty"`
	Description  string `json:"description,omitempty" yaml:"description,omitempty" mapstructure:"description,omitempty"`
	Unit         string `json:"unit,omitempty" yaml:"unit,omitempty"  mapstructure:"unit,omitempty"`
	SensorGroups string `json:"sensorGroups,omitempty" yaml:"sensorGroups,omitempty"  mapstructure:"group_name,omitempty"`
}
