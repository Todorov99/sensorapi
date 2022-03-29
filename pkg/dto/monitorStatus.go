package dto

import (
	"time"

	"github.com/Todorov99/sensorcli/pkg/sensor"
)

type MonitorStatus struct {
	StartTime    time.Time           `json:"startTime,omitempty"`
	Status       string              `json:"status,omitempty"`
	ReportFile   string              `json:"reportFile,omitempty"`
	Error        string              `json:"error,omitempty"`
	Measurements []sensor.Measurment `json:"measurements,omitempty"`
}
