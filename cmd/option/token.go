package option

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// ParseTokenOption parses the token option string and returns the appropriate token.
// The token option can be one of the following:
//
//	none: request without token and follow oauth2
//	anonymous: get anonymous access token once and share between all instances
//	token=<token>: use provided token
func ParseTokenOption(tokenOption string, registry string) (string, error) {
	switch {
	case tokenOption == "none":
		return "", nil
	case tokenOption == "anonymous":
		return getAuthToken(registry, "")
	case strings.HasPrefix(tokenOption, "token="):
		return getAuthToken(registry, strings.TrimPrefix(tokenOption, "token="))
	default:
		return "", fmt.Errorf("invalid token option: %s", tokenOption)
	}
}

var getAuthToken = func(registry string, token string) (string, error) {
	// Get the authentication header
	client := http.DefaultClient
	req, err := http.NewRequest("HEAD", fmt.Sprintf("https://%s/v2/", registry), nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to perform request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	authHeader := resp.Header.Get("Www-Authenticate")
	if authHeader == "" {
		return "", fmt.Errorf("auth header not found")
	}

	// Parse the realm and service from the auth header
	realm := parseChallenge(authHeader, "realm")
	service := parseChallenge(authHeader, "service")
	if realm == "" || service == "" {
		return "", fmt.Errorf("failed to parse realm or service from auth header")
	}

	// Get the token using native Go HTTP client
	tokenURL := fmt.Sprintf("%s?service=%s&scope=repository:*:pull", realm, service)
	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}
	req, err = http.NewRequest("GET", tokenURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create token request: %v", err)
	}

	resp, err = client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to perform token request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code while fetching token: %d", resp.StatusCode)
	}

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return "", fmt.Errorf("failed to read token response body: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		return "", fmt.Errorf("failed to parse token response JSON: %v", err)
	}

	token, ok := result["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("access_token not found or invalid in response")
	}
	return token, nil
}

// parseChallenge extracts the value of a specific key from the WWW-Authenticate header
// It assumes the format is "key=value" and that the value is enclosed in double quotes.
// It returns an empty string if the key is not found or if the value is not properly formatted.
func parseChallenge(header, key string) string {
	prefix := fmt.Sprintf(`%s="`, key)
	start := strings.Index(header, prefix)
	if start == -1 {
		return ""
	}
	start += len(prefix)
	end := strings.Index(header[start:], `"`)
	if end == -1 {
		return ""
	}
	return header[start : start+end]
}
