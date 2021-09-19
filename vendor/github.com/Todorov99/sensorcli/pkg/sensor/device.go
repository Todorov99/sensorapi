package sensor

// Device model
type Device struct {
	ID          string   `json:"id" yaml:"id"`
	Name        string   `json:"name" yaml:"name"`
	Description string   `json:"description" yaml:"description"`
	Sensors     []sensor `json:"sensors" yaml:"sensors"`
}
