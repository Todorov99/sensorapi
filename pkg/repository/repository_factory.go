package repository

import (
	"os"

	"github.com/Todorov99/sensorcli/pkg/logger"
	"github.com/Todorov99/server/pkg/database"
)

var repositoryLogger = logger.NewLogrus("repositroy", os.Stdout)

// CreateMeasurementRepository creates measurement reposiroty.
func CreateMeasurementRepository() Repository {
	return &measurementRepository{
		postgreClient: database.GetDatabaseCfg().GetPostgreClient(),
		influxClient:  database.GetDatabaseCfg().GetInfluxClient(),
		org:           database.GetDatabaseCfg().GetInfluxOrg(),
		bucket:        database.GetDatabaseCfg().GetInfluxBucket(),
	}
}

// CreateSensorRepository creates sensor repository.
func CreateSensorRepository() Repository {
	return &sensorRepository{
		postgreClient: database.GetDatabaseCfg().GetPostgreClient(),
	}
}

// CreateDeviceRepository creates device repository.
func CreateDeviceRepository() Repository {
	return &deviceRepository{
		postgreClient: database.GetDatabaseCfg().GetPostgreClient(),
	}
}

// CreateDeviceRepository creates user repository.
func CreateUserRepository() Repository {
	return &userRepository{
		postgreClient: database.GetDatabaseCfg().GetPostgreClient(),
	}
}
