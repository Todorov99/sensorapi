package config

// import (
// 	_ "github.com/lib/pq"
// )

// type DatabaseClients struct {
// 	InfluxDbClient influxdb2.Client
// 	PostreDbClient *sql.DB
// }

type ApplicationProperties struct {
	InfluxProps  InfluxProperties  `yaml:"influxdb,omitempty"`
	PostgreProps PostgreProperties `yaml:"postgresdb,omitempty"`
	VaultType    string            `yaml:"vaultType,omitempty"`
}

type PostgreProperties struct {
	DatabaseName   string `yaml:"databaseName,omitempty"`
	PasswordSecret string `yaml:"passwordSecret,omitempty"`
	ServiceName    string `yaml:"serviceName,omitempty"`
	SSLMode        string `yaml:"sslmode,omitempty"`
	Port           string `yaml:"port,omitempty"`
}

type InfluxProperties struct {
	DatabaseName string `yaml:"databaseName,omitempty"`
	TokenSecret  string `yaml:"tokenSecret,omitempty"`
	ServiceName  string `yaml:"serviceName,omitempty"`
	Org          string `yaml:"org,omitempty"`
	Bucket       string `yaml:"bucket,omitempty"`
	Port         string `yaml:"port,omitempty"`
}
