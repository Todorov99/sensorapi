package query

const (
	GetSensorAndDeviceBeetweenTimestampQuery = `from(bucket: "my-bucket")
	|> range(start: %s, stop: %s)
	|> filter(fn: (r) => r["_measurement"] == "sensor")
	|> filter(fn: (r) => r["deviceID"] == "%s" and r["sensorID"] == "%s")
	`
)
