package config

import (
	"github.com/Todorov99/server/pkg/global"
)

var serverConfig *serverCfg

type serverCfg struct {
	databaseCfg   *databaseCfg
	jwtCfg        *jwtCfg
	mailSenderCfg *mailSender
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
		databaseCfg:   dbCfg,
		jwtCfg:        jwtConfig,
		mailSenderCfg: NewMailSender(applicationProperties),
	}

	serverConfig = serverCfg
}

func GetDatabaseCfg() *databaseCfg {
	return serverConfig.databaseCfg
}

func GetJWTCfg() *jwtCfg {
	return serverConfig.jwtCfg
}

func GetMailSenderCfg() *mailSender {
	return serverConfig.mailSenderCfg
}
