package config

import (
	"database/sql"
	"fmt"

	"github.com/Todorov99/server/pkg/vault"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type DatabaseCfg struct {
	influxDbClient influxdb2.Client
	influxOrg      string
	influxBucket   string
	postreDbClient *sql.DB
}

func NewDatabaseCfg(applicationProperties *ApplicationProperties) (*DatabaseCfg, error) {
	configLogger.Debug("Getting databaseClients...")

	vault, err := vault.New(applicationProperties.VaultType)
	if err != nil {
		return nil, err
	}

	configLogger.Debug("Initializing influx db 2.0 client")
	influxAddress := fmt.Sprintf("http://%s:%s/", applicationProperties.InfluxProps.ServiceName, applicationProperties.InfluxProps.Port)
	tokenSecret, err := vault.Get(applicationProperties.InfluxProps.TokenSecret)
	if err != nil {
		return nil, err
	}

	influxdbClient := influxdb2.NewClient(influxAddress, tokenSecret.Value)

	configLogger.Debug("Initializing postgres DB client")
	postgreSecret, err := vault.Get(applicationProperties.PostgreProps.PasswordSecret)
	if err != nil {
		return nil, err
	}

	postgreConnectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		postgreSecret.Name, postgreSecret.Value, applicationProperties.PostgreProps.ServiceName, applicationProperties.PostgreProps.Port, applicationProperties.PostgreProps.DatabaseName, applicationProperties.PostgreProps.SSLMode)
	postgreClient, err := sql.Open("postgres", postgreConnectionString)
	if err != nil {
		return nil, err
	}

	err = postgreClient.Ping()
	if err != nil {
		return nil, err
	}

	configLogger.Debug("Database Config successfully initialized")
	return &DatabaseCfg{
		influxDbClient: influxdbClient,
		influxOrg:      applicationProperties.InfluxProps.Org,
		influxBucket:   applicationProperties.InfluxProps.Bucket,
		postreDbClient: postgreClient,
	}, nil
}

func (d *DatabaseCfg) GetInfluxClient() influxdb2.Client {
	return d.influxDbClient
}

func (d *DatabaseCfg) GetPostgreClient() *sql.DB {
	return d.postreDbClient
}

func (d *DatabaseCfg) GetInfluxOrg() string {
	return d.influxOrg
}

func (d *DatabaseCfg) GetInfluxBucket() string {
	return d.influxBucket
}
