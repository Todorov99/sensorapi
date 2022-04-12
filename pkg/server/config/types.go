package config

type ApplicationProperties struct {
	InfluxProps   InfluxProperties     `yaml:"influxdb,omitempty"`
	PostgreProps  PostgreProperties    `yaml:"postgresdb,omitempty"`
	Authorization Authorization        `yaml:"authorization,omitempty"`
	MailSender    MailSenderProperties `yaml:"mailsender,omitempty"`
	VaultType     string               `yaml:"vaultType,omitempty"`
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

type MailSenderProperties struct {
	ServiceName string `yaml:"serviceName,omitempty"`
	Port        string `yaml:"port,omitempty"`
}

type Authorization struct {
	JWT JWTAuthorization `yaml:"JWT,omitempty"`
}

type JWTAuthorization struct {
	JWTAudienceSecret string `yaml:"jwtAudienceSecret,omitempty"`
	JWTIssuerSecret   string `yaml:"jwtIssuerSecret,omitempty"`
	JWTSigningKey     string `yaml:"jwtSigningKey,omitempty"`
	ExpirationTime    string `yaml:"expirationTime,omitempty"`
}
