package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Todorov99/sensorapi/pkg/entity"
	"github.com/Todorov99/sensorapi/pkg/global"
	"github.com/Todorov99/sensorapi/pkg/repository/query"
	"github.com/Todorov99/sensorapi/pkg/server/config"
	"github.com/Todorov99/sensorcli/pkg/logger"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/sirupsen/logrus"
)

type MeasurementRepository interface {
	Add(ctx context.Context, measuerement entity.Measurement) error
	GetMeasurementsFromStartingTime(ctx context.Context, startTime string, userID int) ([]interface{}, error)
	GetMeasurementsBetweenTimestampByDeviceIDBySensorID(ctx context.Context, startTime, endTime, deviceID, sensorID string, userID int) ([]interface{}, error)
	GetMeasurementsValuesBetweenTimestampByDeviceIDAndSensorID(ctx context.Context, startTime, endTime, deviceID, sensorID string, userID int) ([]interface{}, error)
	GetMeasurementsAverageValueBetweenTimestampByDeviceIDAndSensorID(ctx context.Context, startTime, endTime, deviceID, sensorID string, userID int) (string, error)
	CountMeasurementsBetweenTimestampByDeviceIDBySensorID(ctx context.Context, startTime, endTime, deviceID, sensorID string, userID int) (float64, error)
}

type measurementRepository struct {
	logger        *logrus.Entry
	postgreClient *sql.DB
	influxClient  influxdb2.Client
	org           string
	bucket        string
}

func NewMeasurementRepository() MeasurementRepository {
	return &measurementRepository{
		logger:        logger.NewLogrus("deviceRepository", os.Stdout),
		postgreClient: config.GetDatabaseCfg().GetPostgreClient(),
		influxClient:  config.GetDatabaseCfg().GetInfluxClient(),
		org:           config.GetDatabaseCfg().GetInfluxOrg(),
		bucket:        config.GetDatabaseCfg().GetInfluxBucket(),
	}
}

func (m *measurementRepository) GetMeasurementsFromStartingTime(ctx context.Context, startTime string, userID int) ([]interface{}, error) {
	m.logger.Infof("Getting metrics starting from %s for user with ID: %d", startTime, userID)
	influxQuery := fmt.Sprintf(query.GetAllMeasurementsFromStartTime, m.bucket, startTime, userID)
	return executeSelectQueryInflux(ctx, influxQuery, true, m.influxClient, m.org, m.bucket)
}

func (m *measurementRepository) GetMeasurementsBetweenTimestampByDeviceIDBySensorID(ctx context.Context, startTime, endTime, deviceID, sensorID string, userID int) ([]interface{}, error) {
	start, err := parseToRFC3339(startTime)
	if err != nil {
		return nil, err
	}

	end, err := parseToRFC3339(endTime)
	if err != nil {
		return nil, err
	}

	influxQuery := fmt.Sprintf(query.GetMeasurementsBeetweenTimestampByDeviceIdAndSensorId, m.bucket, start, end, deviceID, sensorID, userID)

	measurements, err := executeSelectQueryInflux(ctx, influxQuery, true, m.influxClient, m.org, m.bucket)
	if err != nil {
		return nil, err
	}

	return measurements, nil
}

func (m *measurementRepository) Add(ctx context.Context, measurement entity.Measurement) error {
	rfc3339Timestamp, err := parseToRFC3339(measurement.MeasuredAt)
	if err != nil {
		return fmt.Errorf("invalid measurement timestamp: %w", err)
	}

	measurement.MeasuredAt = rfc3339Timestamp
	writePointToBatch(measurement, m.influxClient, m.org, m.bucket)
	return nil
}

func (m *measurementRepository) GetMeasurementsAverageValueBetweenTimestampByDeviceIDAndSensorID(ctx context.Context, startTime, endTime, deviceID, sensorID string, userID int) (string, error) {
	m.logger.Infof("Getting average value of measurements between %s - %s", startTime, endTime)
	start, err := parseToRFC3339(startTime)
	if err != nil {
		return "", err
	}

	end, err := parseToRFC3339(endTime)
	if err != nil {
		return "", err
	}

	influxQuery := fmt.Sprintf(query.GetAverageValueOfMeasurementsBetweenTimeStampByDeviceIdAndSensorId, m.bucket, start, end, deviceID, sensorID, userID)
	average, err := executeSelectQueryInflux(ctx, influxQuery, true, m.influxClient, m.org, m.bucket)
	if err != nil {
		if strings.Contains(err.Error(), "cannot query an empty range") {
			return "", fmt.Errorf("not available maeasurements in the concrete timestamp %s - %s for user with ID: %d", startTime, endTime, userID)
		}
		return "", err
	}

	if len(average) == 0 {
		return "", fmt.Errorf("not available maeasurements in the concrete timestamp %s - %s for user with ID: %d", startTime, endTime, userID)
	}
	return average[0].(entity.Measurement).Value, nil
}

func (m *measurementRepository) CountMeasurementsBetweenTimestampByDeviceIDBySensorID(ctx context.Context, startTime, endTime, deviceID, sensorID string, userID int) (float64, error) {
	m.logger.Debugf("Getting the count of measurement values between %s - %s for user with ID: %d values...", startTime, endTime, userID)
	start, err := parseToRFC3339(startTime)
	if err != nil {
		return 0, err
	}

	end, err := parseToRFC3339(endTime)
	if err != nil {
		return 0, err
	}

	valueCount, err := executeSelectQueryInflux(ctx, fmt.Sprintf(query.CountMeasurementValues, m.bucket, start, end, deviceID, sensorID, userID), false, m.influxClient, m.org, m.bucket)
	if err != nil {
		return 0, err
	}

	if len(valueCount) == 0 {
		return 0, fmt.Errorf("not available maeasurements in the concrete timestamp %s - %s for user with ID: %d", startTime, endTime, userID)
	}

	return parseFloat(valueCount[0]), nil
}

func (m *measurementRepository) GetMeasurementsValuesBetweenTimestampByDeviceIDAndSensorID(ctx context.Context, startTime, endTime, deviceID, sensorID string, userID int) ([]interface{}, error) {
	start, err := parseToRFC3339(startTime)
	if err != nil {
		return nil, err
	}

	end, err := parseToRFC3339(endTime)
	if err != nil {
		return nil, err
	}

	values, err := executeSelectQueryInflux(ctx, fmt.Sprintf(query.GetMeasurementValuesByDeviceAndSensorIdBeetweenTimestamp, m.bucket, start, end, deviceID, sensorID, userID), false, m.influxClient, m.org, m.bucket)
	if err != nil {
		return nil, err
	}

	return values, nil
}

func parseFloat(v interface{}) float64 {
	switch v.(type) {
	case int64:
		return float64(v.(int64))
	case int32:
		return float64(v.(int32))
	case float64:
		return float64(v.(float64))
	case float32:
		return float64(v.(float64))
	}

	return 0
}

func parseToRFC3339(t string) (string, error) {
	tTime, err := time.Parse(global.TimeFormat, t)
	if err != nil {
		return "", fmt.Errorf("invalid timeformat. Expected %s, received: %s", global.TimeFormat, t)
	}
	tTime = tTime.Add(-time.Hour * 3)
	return tTime.Format(time.RFC3339), nil
}
