package config

import (
	"github.com/Todorov99/server/pkg/global"
)

var serverConfig *serverCfg

type serverCfg struct {
	databaseCfg *DatabaseCfg
	jwtCfg      *JWTCfg
}

func init() {
	applicationProperties, err := LoadApplicationProperties(global.ApplicationPropertyFile)
	if err != nil {
		panic(err)
	}

	dbCfg, err := NewDatabaseCfg(applicationProperties)
	if err != nil {
		panic(err)
	}

	jwtConfig, err := NewJWTCfg(applicationProperties)
	if err != nil {
		panic(err)
	}

	serverCfg := &serverCfg{
		databaseCfg: dbCfg,
		jwtCfg:      jwtConfig,
	}

	serverConfig = serverCfg
}

func GetDatabaseCfg() *DatabaseCfg {
	return serverConfig.databaseCfg
}

func GetJWTCfg() *JWTCfg {
	return serverConfig.jwtCfg
}
