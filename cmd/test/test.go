package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/billy-playground/registry-load-tester/cmd/internal/runner"
)

func main() {
	// Default number of instances is 50
	numInstances := 50
	if len(os.Args) > 1 {
		fmt.Sscanf(os.Args[1], "%d", &numInstances)
	}

	// Registry
	registry := "mcr.azure.cn"

	// Get anonymous token
	token, err := getAuthToken(registry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting auth token: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("json_file,total_size,download_milliseconds")

	// Generate a random array of picked JSON files with length numInstances
	files := make([]string, numInstances)

	allFiles, err := filepath.Glob(filepath.Join("assets/images", "*.json"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading image JSON files: %v\n", err)
		os.Exit(1)
	}
	if len(allFiles) == 0 {
		fmt.Fprintf(os.Stderr, "No JSON files found in assets/images\n")
		os.Exit(1)
	}
	for i := range numInstances {
		files[i] = allFiles[rand.Intn(len(allFiles))]
	}

	// Run instances in parallel
	var wg sync.WaitGroup
	testRunner := runner.NewRunner(token)
	for i := 0; i < numInstances; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = testRunner.StartNew(files[i])
		}()
	}

	// Wait for all instances to complete
	wg.Wait()
}

func getAuthToken(registry string) (string, error) {
	// Get the authentication header
	client := &http.Client{}
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
