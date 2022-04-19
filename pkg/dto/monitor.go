package dto

type MonitorDto struct {
	Duration       string            `json:"duration,omitempty"`
	DeltaDuration  string            `json:"deltaDuration,omitempty"`
	SensorGroups   map[string]string `json:"sensorGroups,omitempty"`
	DeviceID       int               `json:"deviceID,omitempty"`
	GenerateReport bool              `json:"generateReport,omitempty"`
	SendReport     bool              `json:"sendReport,omitempty"`
	MetricValueCfg ValueCfg          `json:"metricValueCfg,omitempty"`
}

type ValueCfg struct {
	TempMax         string `json:"tempMaxValue,omitempty"`
	CPUUsageMax     string `json:"usageMaxValue,omitempty"`
	MemAvailableMax string `json:"memAvailableMaxValue,omitempty"`
	CPUFrequencyMax string `json:"cpuFrequencyMaxValue,omitempty"`
	MemUsedMax      string `json:"memUsedMaxValue,omitempty"`
	MemUsedPercent  string `json:"memUsedPercent,omitempty"`
}
