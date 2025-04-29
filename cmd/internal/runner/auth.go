package runner

import "github.com/billy-playground/registry-load-tester/internal/auth"

// AuthRunner can be used to start a new test instance to download blobs and manifests.
type AuthRunner struct {
	realm        string
	service      string
	registry     string
	refreshToken string
}

func NewAuthRunner(realm string, service string, registry string, refreshToken string) *AuthRunner {
	return &AuthRunner{
		realm:        realm,
		service:      service,
		registry:     registry,
		refreshToken: refreshToken,
	}
}

// StartNew starts a new test instance to download blobs and manifests.
func (r *AuthRunner) StartNew() error {
	_, err := auth.ExchangeToken(r.realm, r.service, r.refreshToken)
	return err
}
