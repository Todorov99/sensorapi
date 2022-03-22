package entity

type Sensor struct {
	ID          int32  `mapstructure:"id,omitempty"`
	Name        string `mapstructure:"name,omitempty"`
	Description string `mapstructure:"description,omitempty"`
	Unit        string `mapstructure:"unit,omitempty"`
	GroupName   string `mapstructure:"group_name,omitempty"`
}
