package entity

type Device struct {
	ID          int32  `mapstructure:"id,omitempty"`
	Name        string `mapstructure:"name,omitempty"`
	Description string `mapstructure:"description,omitempty"`
	Sensors     []Sensor
}
