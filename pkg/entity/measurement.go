package entity

type Measurement struct {
	MeasuredAt string `mapstructure:"measuredAt,omitempty"`
	Value      string `mapstructure:"value,omitempty"`
	SensorID   string `mapstructure:"sensorId,omitempty"`
	DeviceID   string `mapstructure:"deviceId,omitempty"`
	UserID     int    `mapstructure:"userId,omitempty"`
}
