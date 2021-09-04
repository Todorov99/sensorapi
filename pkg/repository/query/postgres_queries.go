package query

const (
	GetSensorNameByID = "Select name from sensor where id=$1"
	GetDeviceNameByID = "Select name from device where id=$1"
)
