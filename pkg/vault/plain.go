package vault

import (
	"fmt"
	"os"
	"sync"

	"github.com/Todorov99/serverapi/pkg/global"
	"gopkg.in/yaml.v2"
)

type plain struct {
	vaultPath string
	mx        sync.RWMutex
}

func newPlainVault() Vault {
	return &plain{
		vaultPath: global.VaultPath,
	}
}

// Get gets a secret from the vault. If the provided secret ID
// does not exist in the vault an error is returned
func (p *plain) Get(secretID string) (Secret, error) {
	vaultLogger.Debugf("Retrieving secret from the vault with ID: %q...", secretID)
	secrets, err := p.read()
	if err != nil {
		return Secret{}, err
	}

	p.mx.RLock()
	defer p.mx.RUnlock()
	secret, ok := secrets[secretID]
	if !ok {
		return Secret{}, fmt.Errorf("secret with ID: %q does not exist in the vault", secretID)
	}

	vaultLogger.Debugf("Secret: %q successfully retrieved", secretID)
	return secret, nil
}

func (p *plain) Store(secret Secret) error {
	vaultLogger.Debug("Saving secret in the vault")
	p.mx.Lock()
	defer p.mx.Unlock()
	f, err := os.Open(p.vaultPath)
	if err != nil {
		return nil
	}
	defer func() {
		f.Close()
	}()

	secretBytes, err := yaml.Marshal(secret)
	if err != nil {
		return err
	}

	_, err = f.Write(secretBytes)
	if err != nil {
		return err
	}
	return nil
}

func (p *plain) read() (map[string]Secret, error) {
	vaultLogger.Debugf("Reading from: %q...", p.vaultPath)

	p.mx.RLock()
	defer p.mx.RUnlock()
	secrets := map[string]Secret{}

	b, err := os.ReadFile(p.vaultPath)
	if err != nil {
		return nil, fmt.Errorf("failed reading from %q: %w", p.vaultPath, err)
	}

	vaultSecrets := []Secret{}
	err = yaml.Unmarshal(b, &vaultSecrets)
	if err != nil {
		return nil, err
	}

	for _, s := range vaultSecrets {
		secrets[s.ID] = s
	}
	vaultLogger.Debug("Vault file successfully retrieved...")
	return secrets, nil
}
