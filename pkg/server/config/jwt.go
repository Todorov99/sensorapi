package config

import (
	"fmt"
	"time"

	"github.com/Todorov99/server/pkg/entity"
	"github.com/Todorov99/server/pkg/vault"
	"github.com/dgrijalva/jwt-go"
)

type JWTCfg struct {
	jwtAudSecret        string
	jwtIssSecret        string
	jwtSigningKeySecret string
	secretVault         vault.Vault
	jwtExpTime          time.Duration
}

func NewJWTCfg(applicationProperties *ApplicationProperties) (*JWTCfg, error) {
	configLogger.Debug("Initializing JWT config...")
	vault, err := vault.New(applicationProperties.VaultType)
	if err != nil {
		return nil, err
	}
	jwtAuthProps := applicationProperties.Authorization.JWT

	fmt.Println(jwtAuthProps)
	expTime, err := time.ParseDuration(jwtAuthProps.ExpirationTime)
	if err != nil {
		return nil, err
	}

	configLogger.Debug("Initializing JWT config finished successfully")
	return &JWTCfg{
		jwtAudSecret:        jwtAuthProps.JWTAudienceSecret,
		jwtIssSecret:        jwtAuthProps.JWTIssuerSecret,
		jwtSigningKeySecret: jwtAuthProps.JWTSigningKey,
		secretVault:         vault,
		jwtExpTime:          expTime,
	}, nil
}

func (j *JWTCfg) GenerateJWT(userEntity entity.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	audSecret, err := j.secretVault.Get(j.jwtAudSecret)
	if err != nil {
		return "", err
	}

	issSecret, err := j.secretVault.Get(j.jwtIssSecret)
	if err != nil {
		return "", err
	}

	signingKeySecret, err := j.secretVault.Get(j.jwtSigningKeySecret)
	if err != nil {
		return "", err
	}

	claims["email"] = userEntity.Email
	claims["name"] = userEntity.FirstName
	// Aud: the services from which the token could be used
	claims["aud"] = audSecret.Value
	// OpenID  "https://ttodorov.com"
	claims["issuer"] = issSecret.Value
	claims["exp"] = time.Now().Add(j.jwtExpTime).Unix()

	tokenString, err := token.SignedString([]byte(signingKeySecret.Value))
	if err != nil {
		return "", fmt.Errorf("wrong signing key: %w", err)
	}

	return tokenString, nil
}

func (j *JWTCfg) GetJWTAudience() string {
	audSecret, err := j.secretVault.Get(j.jwtAudSecret)
	if err != nil {
		return ""
	}
	return audSecret.Value
}

func (j *JWTCfg) GetJWTSigningKey() []byte {
	signingKeySecret, err := j.secretVault.Get(j.jwtSigningKeySecret)
	if err != nil {
		return nil
	}
	return []byte(signingKeySecret.Value)
}

func (j *JWTCfg) GetJWTIssuer() string {
	issSecret, err := j.secretVault.Get(j.jwtIssSecret)
	if err != nil {
		return ""
	}
	return issSecret.Value
}

func (j *JWTCfg) GetJWTExpTimeDuration() time.Duration {
	return j.jwtExpTime
}

func (j *JWTCfg) RenewSigningKey(signKeySecret vault.Secret) error {
	configLogger.Debug("Renewing the JWT signing key...")
	err := j.secretVault.Store(signKeySecret)
	if err != nil {
		return err
	}

	configLogger.Debug("JWT signing key successfully updated")
	return nil
}
