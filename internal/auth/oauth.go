package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// GetAuthHeader tries to authenticate with the registry and get the authentication header.
// If the authentication is successful, it returns the an empty challenge.
func GetAuthHeader(registry string) (string, error) {
	req, err := http.NewRequest("HEAD", fmt.Sprintf("https://%s/v2/", registry), nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to perform request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return "", nil
	}

	if resp.StatusCode != http.StatusUnauthorized {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	authHeader := resp.Header.Get("Www-Authenticate")
	if authHeader == "" {
		return "", fmt.Errorf("auth header not found")
	}
	return authHeader, nil
}

// ExchangeToken exchanges the token using native Go HTTP client.
func ExchangeToken(realm string, service string, token string) (string, error) {
	// Get the token using native Go HTTP client
	tokenURL := fmt.Sprintf("%s?service=%s&scope=repository:*:pull", realm, service)
	req, err := http.NewRequest("GET", tokenURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create exchange token request: %v", err)
	}
	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	resp, err := http.DefaultClient.Do(req)
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

// ParseRealmAndService gets the realm and service from the auth header.
func ParseRealmAndService(authHeader string) (string, string, error) {
	// Parse the realm and service from the auth header
	realm := parseChallenge(authHeader, "realm")
	service := parseChallenge(authHeader, "service")
	// Parse the realm and service from the auth header
	if realm == "" || service == "" {
		return "", "", fmt.Errorf("failed to parse realm or service from auth header")
	}
	return realm, service, nil
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
