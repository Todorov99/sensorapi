package models

//Measurement model
type Measurement struct {
	MeasuredAt string `json:"measuredAt,omitempty"`
	Value      string `json:"value,omitempty"`
	SensorID   string `json:"sensorId,omitempty"`
	DeviceID   string `json:"deviceId,omitempty"`
}

type MeasurementBetweenTimestamp struct {
	StartTime string `json:"startTime,omitempty"`
	EndTime   string `json:"endTime,omitempty"`
	SensorID  string `json:"sensorId,omitempty"`
	DeviceID  string `json:"deviceId,omitempty"`
}
