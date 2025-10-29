package secrets

import (
	"fmt"
	keyring "github.com/zalando/go-keyring"
)

// Service name used in OS keychain entries for codzure.
const serviceName = "codzure"

// SaveKey stores the API key for a given profile in the OS keychain.
func SaveKey(profile string, apiKey string) error {
	if profile == "" {
		return fmt.Errorf("profile name required for keyring entry")
	}
	if apiKey == "" {
		return fmt.Errorf("API key cannot be empty")
	}
	return keyring.Set(serviceName, profile, apiKey)
}

// GetKey retrieves the API key for a given profile from the OS keychain.
func GetKey(profile string) (string, error) {
	if profile == "" {
		return "", fmt.Errorf("profile name required for keyring lookup")
	}
	return keyring.Get(serviceName, profile)
}

// DeleteKey removes the API key entry for the given profile from the OS keychain.
func DeleteKey(profile string) error {
	if profile == "" {
		return fmt.Errorf("profile name required for keyring delete")
	}
	return keyring.Delete(serviceName, profile)
}
