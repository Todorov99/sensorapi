package entity

type SensorGroup struct {
	ID   int32  `mapstructure:"id,omitempty"`
	Name string `mapstructure:"name,omitempty"`
}
