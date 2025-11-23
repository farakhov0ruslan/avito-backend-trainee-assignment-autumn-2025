package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

const (
	maxRetries     = 30
	retryDelay     = 2 * time.Second
	requestTimeout = 5 * time.Second
)

var (
	baseURL    string
	httpClient *http.Client
)

func TestMain(m *testing.M) {
	// Load test environment configuration
	// Try different paths since working directory may vary
	envPaths := []string{".env.e2e", "../../.env.e2e"}
	loaded := false
	for _, path := range envPaths {
		if err := godotenv.Load(path); err == nil {
			loaded = true
			break
		}
	}
	if !loaded {
		fmt.Fprintf(os.Stderr, "Warning: Could not load .env.e2e from any path\n")
	}

	// Get server port from env or use default
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8082"
	}
	baseURL = fmt.Sprintf("http://localhost:%s", port)

	fmt.Printf("Using base URL: %s\n", baseURL)

	// Create HTTP client with timeout
	httpClient = &http.Client{
		Timeout: requestTimeout,
	}

	// Wait for API to be ready
	fmt.Println("Waiting for API server to be ready...")
	if err := waitForAPI(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to API: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("API server is ready!")

	// Run tests
	code := m.Run()

	os.Exit(code)
}

// waitForAPI waits for the API server to become healthy
func waitForAPI() error {
	healthURL := fmt.Sprintf("%s/health", baseURL)

	for i := 0; i < maxRetries; i++ {
		resp, err := httpClient.Get(healthURL)
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return nil
		}
		if resp != nil {
			resp.Body.Close()
		}

		time.Sleep(retryDelay)
	}

	return fmt.Errorf("API did not become healthy after %d attempts", maxRetries)
}

// Helper functions for making HTTP requests

// doRequest performs an HTTP request and returns the response
func doRequest(method, path string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, baseURL+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return httpClient.Do(req)
}

// doGet performs a GET request with query parameters
func doGet(path string, queryParams map[string]string) (*http.Response, error) {
	u, err := url.Parse(baseURL + path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	if len(queryParams) > 0 {
		q := u.Query()
		for key, value := range queryParams {
			q.Set(key, value)
		}
		u.RawQuery = q.Encode()
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	return httpClient.Do(req)
}

// parseResponse parses JSON response into the given structure
func parseResponse(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(bodyBytes, v); err != nil {
		return fmt.Errorf("failed to unmarshal response (body: %s): %w", string(bodyBytes), err)
	}

	return nil
}

// assertStatusCode checks if the response has the expected status code
func assertStatusCode(t *testing.T, resp *http.Response, expected int) {
	t.Helper()
	if resp.StatusCode != expected {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Errorf("Expected status code %d, got %d. Response body: %s",
			expected, resp.StatusCode, string(bodyBytes))
	}
}

// assertErrorCode checks if the error response has the expected error code
func assertErrorCode(t *testing.T, resp *http.Response, expectedCode string) {
	t.Helper()
	defer resp.Body.Close()

	var errorResp struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := parseResponse(resp, &errorResp); err != nil {
		t.Fatalf("Failed to parse error response: %v", err)
	}

	if errorResp.Error.Code != expectedCode {
		t.Errorf("Expected error code %s, got %s (message: %s)",
			expectedCode, errorResp.Error.Code, errorResp.Error.Message)
	}
}
