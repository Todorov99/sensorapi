package models

type ValueCfg struct {
	TempMax         string `json:"tempMaxValue,omitempty"`
	UsageMax        string `json:"usageMaxValue,omitempty"`
	MemAvailableMax string `json:"memAvailableMaxValue,omitempty"`
}

func NewValueCfg(tempMax, usageMax, memAvailableMax string) *ValueCfg {
	return &ValueCfg{
		TempMax:         tempMax,
		UsageMax:        usageMax,
		MemAvailableMax: memAvailableMax,
	}
}
