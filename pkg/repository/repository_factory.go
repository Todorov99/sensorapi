package repository

import "github.com/Todorov99/sensorcli/pkg/logger"

var repositoryLogger = logger.NewLogger("./repositroy")

// CreateMeasurementRepository creates measurement reposiroty.
func CreateMeasurementRepository() Repository {
	return &measurementRepository{}
}

// CreateSensorRepository creates sensor repository.
func CreateSensorRepository() Repository {
	return &sensorRepository{}
}

// CreateDeviceRepository creates device repository.
func CreateDeviceRepository() Repository {
	return &deviceRepository{}
}
