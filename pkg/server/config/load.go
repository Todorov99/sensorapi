package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Todorov99/sensorcli/pkg/logger"
	_ "github.com/lib/pq"
	"gopkg.in/yaml.v2"
)

var configLogger = logger.NewLogrus("config", os.Stdout)

func LoadApplicationProperties(propsFile string) (*ApplicationProperties, error) {
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
