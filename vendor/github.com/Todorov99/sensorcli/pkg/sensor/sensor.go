package sensor

type Sensor struct {
	ID           int32        `json:"id" yaml:"id" mapstructure:"id,omitempty"`
	Name         string       `json:"name" yaml:"name" mapstructure:"name,omitempty"`
	Description  string       `json:"description" yaml:"description" mapstructure:"description,omitempty"`
	Unit         string       `json:"unit" yaml:"unit" mapstructure:"unit,omitempty"`
	SensorGroups string       `json:"sensorGroups" yaml:"sensorGroups" mapstructure:"group_name,omitempty"`
	Measurments  []Measurment `json:"measurements" yaml:"measurements"`
}
