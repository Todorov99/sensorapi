package dto

type ValueCfg struct {
	TempMax         string `json:"tempMaxValue,omitempty"`
	CPUUsageMax     string `json:"usageMaxValue,omitempty"`
	MemAvailableMax string `json:"memAvailableMaxValue,omitempty"`
	CPUFrequencyMax string `json:"cpuFrequencyMaxValue,omitempty"`
	MemUsedMax      string `json:"memUsedMaxValue,omitempty"`
	MemUsedPercent  string `json:"memUsedPercent,omitempty"`
}
