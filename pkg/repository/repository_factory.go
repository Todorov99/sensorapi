package repository

import (
	"os"

	"github.com/Todorov99/sensorcli/pkg/logger"
	"github.com/Todorov99/server/pkg/server/config"
)

var repositoryLogger = logger.NewLogrus("repositroy", os.Stdout)

// CreateMeasurementRepository creates measurement reposiroty.
func CreateMeasurementRepository() Repository {
	return &measurementRepository{
		postgreClient: config.GetDatabaseCfg().GetPostgreClient(),
		influxClient:  config.GetDatabaseCfg().GetInfluxClient(),
		org:           config.GetDatabaseCfg().GetInfluxOrg(),
		bucket:        config.GetDatabaseCfg().GetInfluxBucket(),
	}
}

// CreateSensorRepository creates sensor repository.
func CreateSensorRepository() Repository {
	return &sensorRepository{
		postgreClient: config.GetDatabaseCfg().GetPostgreClient(),
	}
}

// CreateDeviceRepository creates device repository.
func CreateDeviceRepository() Repository {
	return &deviceRepository{
		postgreClient: config.GetDatabaseCfg().GetPostgreClient(),
	}
}

// CreateDeviceRepository creates user repository.
func CreateUserRepository() Repository {
	return &userRepository{
		postgreClient: config.GetDatabaseCfg().GetPostgreClient(),
	}
}
