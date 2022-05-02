package config

import (
	"github.com/Todorov99/sensorapi/pkg/global"
)

var serverConfig *serverCfg

type serverCfg struct {
	databaseCfg   *databaseCfg
	jwtCfg        *jwtCfg
	mailSenderCfg *mailSender
	vault         string
	userCfg       User
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
		vault:         applicationProperties.VaultType,
		userCfg:       applicationProperties.User,
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

func GetUserCfg() User {
	return serverConfig.userCfg
}

func GetVault() string {
	return serverConfig.vault
}
