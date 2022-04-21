package query

//PostreSQL queries
const (
	//Sensor queries
	GetSensorNameByID       = "Select name from sensor where id=$1"
	GetAllSensors           = "SELECT distinct s.id, s.name, s.description, s.unit, ss.group_name FROM device_sensor as ds join sensor as s on ds.sensor_id = s.id join sensor_groups as ss on s.sensor_groups_id = ss.id ORDER BY s.id"
	GetAllSensorsBySensorID = "SELECT distinct s.id, s.name, s.description, s.unit, ss.group_name FROM device_sensor as ds join sensor as s on ds.sensor_id = s.id join sensor_groups as ss on s.sensor_groups_id = ss.id where s.id = $1"
	GetAllSensorsByDeviceID = `SELECT distinct s.id, s.name, s.description, s.unit, ss.group_name FROM device_sensor as ds join sensor as s on ds.sensor_id = s.id join sensor_groups as ss on s.sensor_groups_id = ss.id where ds.device_id = $1 ORDER BY s.id`

	AddSensor              = `INSERT INTO sensor(name,description,device_id,sensor_groups_id,unit) VALUES ($1,$2,$3,$4,$5)`
	UpdateSensor           = `UPDATE sensor set name=$1,description=$2,sensor_groups_id=$3,unit=$4 where id=$5`
	DeleteSensor           = `DELETE from sensor where id=$1`
	GetSensorIDByGroupName = "SELECT id from sensor_groups where group_name=$1"
	GetSensorIDByDeviceID  = "SELECT s.id from sensor as s where s.device_id=$1"
	GetSensorByID          = "SELECT s.id, s.name, s.description, s.unit, ss.group_name FROM sensor as s join sensor_groups as ss on s.sensor_groups_id = ss.id where s.id=$1"
	GetSensorByName        = "SELECT id from sensor where name=$1"

	//Device queries
	GetDeviceNameByID = "Select name from device where id=$1"
	GetDeviceIDByName = "SELECT id FROM device where name=$1"
	InsertDevice      = `INSERT INTO device(name,description,user_id) VALUES ($1,$2,$3)`

	UpdateUserDevice = `UPDATE device set user_id=$1`

	UpdateDevice = `UPDATE device set name=$1,description=$2 where id=$3 and user_id=$4`
	DeleteDevice = `DELETE from device where id=$1 and user_id=$2`

	GetDeviceByID      = `SELECT distinct d.id, d.name, d.description from device as d where d.id=$1 and d.user_id=$2`
	GetAllDevices      = `SELECT distinct d.id, d.name, d.description from device_sensor as dv join device as d on dv.device_id = d.id where d.user_id =$1 ORDER BY d.id`
	GetHighestDeviceID = "SELECT max(id) + 1 from device"

	//User queries
	GetHighestUserID = "SELECT max(id) + 1 from users"
	GetUserIDByName  = "SELECT id FROM users where user_name=$1"
	GetUserIDByEmail = "SELECT id FROM users where email=$1"
	GetUserByID      = "SELECT s.user_name s.pass FROM user as u where ID=$1"
	GetUserByName    = "SELECT * FROM users where user_name=$1"
	InsertUser       = `INSERT INTO users(user_name,pass,first_name,last_name,email) VALUES ($1,$2,$3,$4,$5)`

	//Device_sensor queries
	InsertDeviceSensors = `INSERT INTO device_sensor(device_id,sensor_id) VALUES ($1,$2)`
)

//InfluxDB 2.0 queries
const (
	GetMeasurementsBeetweenTimestampByDeviceIdAndSensorId = `from(bucket: "%s")
	|> range(start: 1999-09-04T23:30:00Z)
	|> filter(fn: (r) => r["_time"] >= %s and r["_time"] <= %s)
	|> filter(fn: (r) => r["_measurement"] == "sensor")
	|> filter(fn: (r) => r["deviceID"] == "%s" and r["sensorID"] == "%s")
	`
	GetAverageValueOfMeasurementsBetweenTimeStampByDeviceIdAndSensorId = `from(bucket: "%s")
	|> range(start: 1999-09-04T23:30:00Z)
	|> filter(fn: (r) => r["_time"] >= %s and r["_time"] <= %s)
	|> filter(fn: (r) => r["_measurement"] == "sensor")
	|> filter(fn: (r) => r["deviceID"] == "%s" and r["sensorID"] == "%s")
	|> mean()
	`

	GetMeasurementValuesByDeviceAndSensorIdBeetweenTimestamp = `from(bucket: "%s")
	|> range(start: 1999-09-04T23:30:00Z)
	|> filter(fn: (r) => r["_time"] >= %s and r["_time"] <= %s)
	|> filter(fn: (r) => r["_measurement"] == "sensor")
	|> filter(fn: (r) => r["deviceID"] == "%s" and r["sensorID"] == "%s")
	|> keep(columns: ["_value"])
	`

	CountMeasurementValues = `from(bucket: "%s")
	|> range(start: 1999-09-04T23:30:00Z)
	|> filter(fn: (r) => r["_time"] >= %s and r["_time"] <= %s)
	|> filter(fn: (r) => r["_measurement"] == "sensor")
	|> filter(fn: (r) => r["deviceID"] == "%s" and r["sensorID"] == "%s")
	|> count()
	`

	GetAllMeasurementsFromStartTime = `from(bucket: "%s")
	|> range(start: 1999-09-04T23:30:00Z)
	|> filter(fn: (r) => r["_time"] >= %s)
	|> filter(fn: (r) => r["_measurement"] == "sensor")
	`
)
