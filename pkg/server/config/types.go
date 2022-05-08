package config

type ApplicationProperties struct {
	Services  Services `yaml:"services,omitempty"`
	Security  Security `yaml:"security,omitempty"`
	VaultType string   `yaml:"vaultType,omitempty"`
	User      User     `yaml:"user,omitempty"`
}

type Security struct {
	TLS            TLS            `yaml:"tls,omitempty"`
	Authentication Authentication `yaml:"authentication,omitempty"`
}

type TLS struct {
	CertFile   string `yaml:"certFile,omitempty"`
	PrivateKey string `yaml:"privateKey,omitempty"`
	RootCACert string `yaml:"rootCACert,omitempty"`
	RootCAKey  string `yaml:"rootCAKey,omitempty"`
}

type Services struct {
	InfluxProps  InfluxProperties     `yaml:"influxdb,omitempty"`
	PostgreProps PostgreProperties    `yaml:"postgresdb,omitempty"`
	MailSender   MailSenderProperties `yaml:"mailsender,omitempty"`
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

type Authentication struct {
	JWT JWTAuth `yaml:"JWT,omitempty"`
}

type JWTAuth struct {
	JWTAudienceSecret string `yaml:"jwtAudienceSecret,omitempty"`
	JWTIssuerSecret   string `yaml:"jwtIssuerSecret,omitempty"`
	JWTSigningKey     string `yaml:"jwtSigningKey,omitempty"`
	ExpirationTime    string `yaml:"expirationTime,omitempty"`
}

type User struct {
	UserSecret string `yaml:"userSecret"`
	FirstName  string `yaml:"firstname"`
	LastName   string `yaml:"lastname"`
	Email      string `yaml:"email"`
}
