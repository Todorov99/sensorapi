package models

//Measurement model
type Measurement struct {
	MeasuredAt string `json:"measuredAt,omitempty" yaml:"measuredAt,omitempty"`
	Value      string `json:"value,omitempty" yaml:"value,omitempty"`
	SensorID   string `json:"sensorId,omitempty" yaml:"sensorId,omitempty"`
	DeviceID   string `json:"deviceId,omitempty" yaml:"deviceId,omitempty"`
}
