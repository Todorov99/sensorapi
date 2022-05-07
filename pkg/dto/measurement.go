package dto

//Measurement model
type Measurement struct {
	MeasuredAt string `json:"measuredAt,omitempty"`
	Value      string `json:"value,omitempty"`
	SensorID   string `json:"sensorId,omitempty"`
	DeviceID   string `json:"deviceId,omitempty"`
	UserID     int    `json:"userId,omitempty"`
}
