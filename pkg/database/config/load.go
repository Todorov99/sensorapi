package config

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Todorov99/sensorcli/pkg/logger"
	"github.com/Todorov99/server/pkg/vault"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	_ "github.com/lib/pq"
	"gopkg.in/yaml.v2"
)

var configLogger = logger.NewLogrus("config", os.Stdout)

type DatabaseCfg struct {
	influxDbClient influxdb2.Client
	influxOrg      string
	influxBucket   string
	postreDbClient *sql.DB
}

func NewDatabaseClients(propsFile string) (*DatabaseCfg, error) {
	configLogger.Debug("Getting databaseClients...")

	applicationProperties, err := loadApplicationProperties(propsFile)
	if err != nil {
		return nil, err
	}

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

	configLogger.Debug("Database clients successfully initialized")
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

func loadApplicationProperties(propsFile string) (*ApplicationProperties, error) {
	appPropersties := &ApplicationProperties{}
	absoluteFilePath, err := filepath.Abs(propsFile)
	if err != nil {
		return nil, fmt.Errorf("failed getting absolute path form: %q", propsFile)
	}

	configLogger.Debugf("Loading property file: %q...", absoluteFilePath)
	b, err := os.ReadFile(absoluteFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed reading config file from: %q", absoluteFilePath)
	}

	err = yaml.Unmarshal(b, appPropersties)
	if err != nil {
		return nil, err
	}
	configLogger.Debug("Property file successfully loaded")
	return appPropersties, nil
}
