package vault

import (
	"fmt"

	"github.com/Todorov99/server/pkg/global"
)

type Vault interface {
	// Get gets a secret from the vault by provided secret ID
	Get(secretID string) (Secret, error)
}

type Secret struct {
	ID    string `yaml:"id"`
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

// New returns an instance of the requeried vault type
func New(vaultType string) (Vault, error) {
	switch vaultType {
	case global.PlainVaultType:
		return newPlainVault(), nil
	}

	return nil, fmt.Errorf("there is not %q existing vault", vaultType)
}
