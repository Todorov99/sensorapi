package repository

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/Todorov99/server/pkg/repository/query"
	"github.com/mitchellh/mapstructure"
)

func executeSelectQuery(query string, postgreClient *sql.DB, entity interface{}, args ...interface{}) error {
	var objects []map[string]interface{}

	rowsRs, err := postgreClient.Query(query, args...)
	//TODO make proper check
	if rowsRs == nil {
		return nil
	}

	columns, cerr := rowsRs.ColumnTypes()
	if cerr != nil {
		return cerr
	}

	for rowsRs.Next() {

		if len(columns) != 1 {
			// Scan needs an array of pointers to the values it is setting
			// This creates the object and sets the values correctly
			vals := make([]interface{}, len(columns))
			object := map[string]interface{}{}
			for i, column := range columns {
				object[column.Name()] = reflect.New(column.ScanType()).Interface()
				vals[i] = object[column.Name()]
			}

			err = rowsRs.Scan(vals...)
			if err != nil {
				return err
			}
			objects = append(objects, object)
		} else {
			err = rowsRs.Scan(entity)
			if err != nil {
				return err
			}
		}

	}

	if err != nil {
		return err
	}

	if len(columns) != 1 {

		switch reflect.Indirect(reflect.ValueOf(entity)).Kind() {
		case reflect.Slice:
			err = mapstructure.Decode(objects, entity)
			if err != nil {
				return err
			}
		case reflect.Struct:
			err = mapstructure.Decode(objects[0], entity)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported type")

		}

	}

	return nil
}

func executeModifyingQuery(query string, postgreClient *sql.DB, args ...interface{}) error {
	_, err := postgreClient.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed executing query %q with arguments %q: %w", query, args, err)
	}

	return nil
}

func checkForExistingSensorByID(id string, postgreClient *sql.DB) bool {
	var existingSensor int64
	_ = executeSelectQuery(query.GetSensorNameByID, postgreClient, &existingSensor, id)

	return existingSensor == 0
}

func checkForExistingDeviceByID(id string, postgreClient *sql.DB) bool {
	var existingDevice int64
	_ = executeSelectQuery(query.GetDeviceNameByID, postgreClient, &existingDevice, id)
	return existingDevice == 0
}

func checkForExistingDevicesAndSensors(deviceID string, sensorID string, postgreClient *sql.DB) error {
	repositoryLogger.Info("Checking for existing device and sensor...")
	if !checkForExistingDeviceByID(deviceID, postgreClient) {
		return fmt.Errorf("failed getting device with %s ID", deviceID)
	}

	if !checkForExistingSensorByID(sensorID, postgreClient) {
		return fmt.Errorf("failed getting sensor with %s ID", sensorID)
	}

	return nil
}
