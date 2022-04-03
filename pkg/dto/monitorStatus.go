package dto

import (
	"github.com/Todorov99/sensorcli/pkg/sensor"
)

type MonitorStatus struct {
	StartTime           string              `json:"startedAt,omitempty"`
	FinishedAt          string              `json:"finishedAt,omitempty"`
	Status              string              `json:"status,omitempty"`
	ReportFile          string              `json:"reportFile,omitempty"`
	Error               string              `json:"error,omitempty"`
	CriticalMeasurement []sensor.Measurment `json:"criticalMeasurement,omitempty"`
	Measurements        []sensor.Measurment `json:"measurements,omitempty"`
}
