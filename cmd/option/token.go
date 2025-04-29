package option

import (
	"fmt"
	"strings"

	"github.com/billy-playground/registry-load-tester/internal/auth"
)

// Token represents the token option for the registry load tester.
type Token struct {
	tokenModeInput string
	AccessToken    string
}

// SetFlag sets the token mode for the token option.
func (t *Token) SetFlag(tokenMode string) {
	t.tokenModeInput = tokenMode
}

// Parse retrieves the appropriate token based on the token mode.
// The token option can be one of the following:
//
//	none: request without token and follow oauth2
//	anonymous: get anonymous access token once and share between all instances
//	token=<token>: use provided token
func (t *Token) Parse(registry string) (err error) {
	switch {
	case t.tokenModeInput == "none":
		return nil
	case t.tokenModeInput == "anonymous":
		t.AccessToken, err = getAuthToken(registry, "")
		return err
	case strings.HasPrefix(t.tokenModeInput, "token="):
		t.AccessToken, err = getAuthToken(registry, strings.TrimPrefix(t.tokenModeInput, "token="))
		return err
	default:
		return fmt.Errorf("invalid token option: %s", t.tokenModeInput)
	}
}

var getAuthToken = func(registry string, refreshToken string) (string, error) {
	authHeader, err := auth.GetAuthHeader(registry)
	if err != nil {
		return "", err
	}

	if authHeader == "" {
		// no access token obtained since the registry requires no authN at all
		return "", nil
	}
	realm, service, err := auth.ParseRealmAndService(authHeader)
	if err != nil {
		return "", err
	}

	return auth.ExchangeToken(realm, service, refreshToken)
}
