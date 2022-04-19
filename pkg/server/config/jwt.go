package config

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Todorov99/sensorapi/pkg/entity"
	"github.com/Todorov99/sensorapi/pkg/vault"
	"github.com/dgrijalva/jwt-go"
)

type jwtCfg struct {
	jwtAudSecret        string
	jwtIssSecret        string
	jwtSigningKeySecret string
	secretVault         vault.Vault
	jwtExpTime          time.Duration
}

func NewJWTCfg(applicationProperties *ApplicationProperties) (*jwtCfg, error) {
	configLogger.Debug("Initializing JWT config...")
	vault, err := vault.New(applicationProperties.VaultType)
	if err != nil {
		return nil, err
	}
	jwtAuthProps := applicationProperties.Authorization.JWT

	expTime, err := time.ParseDuration(jwtAuthProps.ExpirationTime)
	if err != nil {
		return nil, err
	}

	configLogger.Debug("Initializing JWT config finished successfully")
	return &jwtCfg{
		jwtAudSecret:        jwtAuthProps.JWTAudienceSecret,
		jwtIssSecret:        jwtAuthProps.JWTIssuerSecret,
		jwtSigningKeySecret: jwtAuthProps.JWTSigningKey,
		secretVault:         vault,
		jwtExpTime:          expTime,
	}, nil
}

func (j *jwtCfg) GenerateJWT(userEntity entity.User) (string, error) {
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

func (j *jwtCfg) GetJWTAudience() string {
	audSecret, err := j.secretVault.Get(j.jwtAudSecret)
	if err != nil {
		return ""
	}
	return audSecret.Value
}

func (j *jwtCfg) GetJWTSigningKey() []byte {
	signingKeySecret, err := j.secretVault.Get(j.jwtSigningKeySecret)
	if err != nil {
		return nil
	}
	return []byte(signingKeySecret.Value)
}

func (j *jwtCfg) GetJWTIssuer() string {
	issSecret, err := j.secretVault.Get(j.jwtIssSecret)
	if err != nil {
		return ""
	}
	return issSecret.Value
}

func (j *jwtCfg) GetJWTExpTimeDuration() time.Duration {
	return j.jwtExpTime
}

func (j *jwtCfg) RenewSigningKey(signKeySecret vault.Secret) error {
	configLogger.Debug("Renewing the JWT signing key...")
	err := j.secretVault.Store(signKeySecret)
	if err != nil {
		return err
	}

	configLogger.Debug("JWT signing key successfully updated")
	return nil
}

// HTTPError represents a HTTP response error
type HTTPError struct {
	err        error
	statusCode int
}

func (h HTTPError) Error() error {
	return h.err
}

func (h HTTPError) StatusCode() int {
	return h.statusCode
}

// ParseToken parses the token provided in the request header
func ParseToken(w http.ResponseWriter, r *http.Request) (*jwt.Token, HTTPError) {
	bareerToken := r.Header.Get("Authorization")
	if bareerToken == "" {
		return nil, HTTPError{
			err:        fmt.Errorf("no Authorization toke provided"),
			statusCode: http.StatusForbidden,
		}
	}

	t := strings.Split(bareerToken, " ")[1]

	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		jwtCfg := GetJWTCfg()
		tokenClaims := token.Claims.(jwt.MapClaims)
		if !tokenClaims.VerifyAudience(jwtCfg.GetJWTAudience(), false) {
			return nil, errors.New("invalid aud of the JWT")
		}

		if !tokenClaims.VerifyIssuer(jwtCfg.GetJWTIssuer(), false) {
			return nil, errors.New("invalid iss of the JWT")
		}

		return jwtCfg.GetJWTSigningKey(), nil
	})

	if err != nil {
		return nil, HTTPError{
			err:        err,
			statusCode: http.StatusForbidden,
		}
	}

	return token, HTTPError{
		err:        nil,
		statusCode: http.StatusOK,
	}
}

func GetJWTEmailClaim(token *jwt.Token) string {
	tokenClaims := token.Claims.(jwt.MapClaims)
	emailClaim := tokenClaims["email"].(string)
	return emailClaim
}
