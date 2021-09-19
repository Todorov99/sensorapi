package query

//PostresSQL queries
const (
	//Sensor queries
	GetSensorNameByID       = "Select name from sensor where id=$1"
	GetAllSensors           = "SELECT s.id, s.name, s.description, s.unit, ss.group_name FROM sensor as s join sensor_groups as ss on s.sensor_groups_id = ss.id"
	GetAllSensorsBySensorID = "SELECT s.id, s.name, s.description, s.unit, ss.group_name FROM sensor as s join sensor_groups as ss on s.sensor_groups_id = ss.id where s.id = $1"
	GetAllSensorsByDeviceID = "SELECT s.id, s.name, s.description, s.unit, ss.group_name FROM sensor as s join sensor_groups as ss on s.sensor_groups_id = ss.id where s.device_id = $1"
	AddSensor               = `INSERT INTO sensor(id,name,description,device_id,sensor_groups_id,unit) VALUES ($1,$2,$3,$4,$5,$6)`
	UpdateSensor            = `UPDATE sensor set name=$1,description=$2,sensor_groups_id=$3,unit=$4 where id=$5`
	DeleteSensor            = `DELETE from sensor where id=$1`
	GetSensorIDByGroupName  = "SELECT id from sensor_groups where group_name=$1"
	GetSensorIDByDeviceID   = "SELECT s.id from sensor as s where s.device_id=$1"
	GetSensorByID           = "SELECT s.id, s.name, s.description, s.unit, ss.group_name FROM sensor as s join sensor_groups as ss on s.sensor_groups_id = ss.id where s.id=$1"
	GetSensorByName         = "SELECT id from sensor where name=$1"

	//Device queries
	GetDeviceNameByID  = "Select name from device where id=$1"
	GetDeviceIDByName  = "SELECT id FROM device where name=$1"
	InsertDevice       = `INSERT INTO device(id,name,description) VALUES ($1,$2,$3)`
	UpdateDevice       = `UPDATE device set name=$1,description=$2 where id=$3`
	DeleteDevice       = `DELETE from device where id=$1`
	GetDeviceByID      = `SELECT d.id, d.name, d.description from device as d where d.id=$1`
	GetAllDevices      = `SELECT d.id, d.name, d.description from device as d`
	GetHighestDeviceID = "SELECT max(id) + 1 from device"
)

//InfluxDB 2.0 queries
const (
	GetMeasurementsBeetweenTimestampByDeviceIdAndSensorId = `from(bucket: "my-bucket")
	|> range(start: %s, stop: %s)
	|> filter(fn: (r) => r["_measurement"] == "sensor")
	|> filter(fn: (r) => r["deviceID"] == "%s" and r["sensorID"] == "%s")
	`
	GetAverageValueOfMeasurementsBetweenTimeStampByDeviceIdAndSensorId = `from(bucket: "my-bucket")
	|> range(start: %s, stop: %s)
	|> filter(fn: (r) => r["_measurement"] == "sensor")
	|> filter(fn: (r) => r["deviceID"] == "%s" and r["sensorID"] == "%s")
	|> mean()
	`

	GetMeasurementValuesByDeviceAndSensorIdBeetweenTimestamp = `from(bucket: "my-bucket")
	|> range(start: %s, stop: %s)
	|> filter(fn: (r) => r["_measurement"] == "sensor")
	|> filter(fn: (r) => r["deviceID"] == "%s" and r["sensorID"] == "%s")
	|> keep(columns: ["_value"])
	`

	CountMeasurementValues = `from(bucket: "my-bucket")
	|> range(start: %s, stop: %s)
	|> filter(fn: (r) => r["_measurement"] == "sensor")
	|> filter(fn: (r) => r["deviceID"] == "%s" and r["sensorID"] == "%s")
	|> count()
	`

	GetAllMeasurementsFromStartTime = `from(bucket: "my-bucket")
	|> range(start: 2021-09-04T23:30:00Z)
	|> filter(fn: (r) => r["_measurement"] == "sensor")
	`
)
