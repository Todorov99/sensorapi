package models

// Device model
type Device struct {
	ID          string   `json:"id,omitempty" yaml:"id,omitempty"`
	Name        string   `json:"name,omitempty" yaml:"name,omitempty"`
	Description string   `json:"description,omitempty" yaml:"description,omitempty"`
	Sensors     []Sensor `json:"sensors,omitempty" yaml:"sensors,omitempty"`
}
