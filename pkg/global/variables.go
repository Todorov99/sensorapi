package global

const (
	TempSensor        = "1"
	UsageSensor       = "2"
	CoresSensor       = "3"
	FrequencySensor   = "4"
	TotalMemory       = "5"
	MemoryAvailable   = "6"
	MemoryUsed        = "7"
	MemoryUsedParcent = "8"
)

const (
	CpuTempGroup  = "CPU_TEMP"
	CpuUsageGroup = "CPU_USAGE"
	MemoryGroup   = "MEMORY_USAGE"
)

const (
	VaultPath               = "./cfg/vault.yaml"
	PlainVaultType          = "plain"
	ApplicationPropertyFile = "./cfg/application_properties.yaml"
	CfgFileName             = "./cliresources/device_cfg.yaml"
	CliResourceDir          = "./cliresources"
	CliBinariesDir          = "./binaries"
	CliZipCfg               = "cli_cfg.zip"
	DeviceCfgChecksum       = "./cliresources/.checksum"
)

const (
	TimeFormat = "2006-01-02-15:04:05"
)
