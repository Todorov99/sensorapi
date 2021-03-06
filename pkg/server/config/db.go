package config

import (
	"database/sql"
	"fmt"

	"github.com/Todorov99/sensorapi/pkg/vault"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type databaseCfg struct {
	influxDbClient influxdb2.Client
	influxOrg      string
	influxBucket   string
	postreDbClient *sql.DB
}

func NewDatabaseCfg(applicationProperties *ApplicationProperties) (*databaseCfg, error) {
	configLogger.Debug("Getting databaseClients...")

	vault, err := vault.New(applicationProperties.VaultType)
	if err != nil {
		return nil, err
	}

	configLogger.Debug("Initializing influx db 2.0 client")
	influxAddress := fmt.Sprintf("http://%s:%s/", applicationProperties.Services.InfluxProps.ServiceName, getEnv(applicationProperties.Services.InfluxProps.Port))
	tokenSecret, err := vault.Get(applicationProperties.Services.InfluxProps.TokenSecret)
	if err != nil {
		return nil, err
	}

	influxdbClient := influxdb2.NewClient(influxAddress, tokenSecret.Value)

	configLogger.Debug("Initializing postgres DB client")
	postgreSecret, err := vault.Get(applicationProperties.Services.PostgreProps.PasswordSecret)
	if err != nil {
		return nil, err
	}

	postgreConnectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		postgreSecret.Name, postgreSecret.Value, applicationProperties.Services.PostgreProps.ServiceName, getEnv(applicationProperties.Services.PostgreProps.Port), applicationProperties.Services.PostgreProps.DatabaseName, applicationProperties.Services.PostgreProps.SSLMode)
	postgreClient, err := sql.Open("postgres", postgreConnectionString)
	if err != nil {
		return nil, err
	}

	err = postgreClient.Ping()
	if err != nil {
		return nil, err
	}

	configLogger.Debug("Database Config successfully initialized")
	return &databaseCfg{
		influxDbClient: influxdbClient,
		influxOrg:      applicationProperties.Services.InfluxProps.Org,
		influxBucket:   applicationProperties.Services.InfluxProps.Bucket,
		postreDbClient: postgreClient,
	}, nil
}

func (d *databaseCfg) GetInfluxClient() influxdb2.Client {
	return d.influxDbClient
}

func (d *databaseCfg) GetPostgreClient() *sql.DB {
	return d.postreDbClient
}

func (d *databaseCfg) GetInfluxOrg() string {
	return d.influxOrg
}

func (d *databaseCfg) GetInfluxBucket() string {
	return d.influxBucket
}
