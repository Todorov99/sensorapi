package database

import (
	"github.com/Todorov99/server/pkg/database/config"
	"github.com/Todorov99/server/pkg/global"
)

var databaseCfg *config.DatabaseCfg

func init() {
	dbClients, err := config.NewDatabaseClients(global.ApplicationPropertyFile)
	if err != nil {
		panic(err)
	}

	databaseCfg = dbClients
}

func GetDatabaseCfg() *config.DatabaseCfg {
	return databaseCfg
}
