package dto

// Device model
type Device struct {
	ID          int32    `json:"id,omitempty" mapstructure:"id,omitempty"`
	Name        string   `json:"name,omitempty" mapstructure:"name,omitempty"`
	Description string   `json:"description,omitempty" mapstructure:"description,omitempty"`
	Sensors     []Sensor `json:"sensors,omitempty" mapstructure:"sensors,omitempty"`
}

type AddUpdateDeviceDto struct {
	ID          int    `json:"id,omitempty" mapstructure:"id,omitempty"`
	Name        string `json:"name" mapstructure:"name"`
	Description string `json:"description" mapstructure:"description"`
}
