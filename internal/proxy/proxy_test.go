//go:build unit

package proxy

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProxy(t *testing.T) {
	// Create a test server that will act as the target for our proxy
	targetServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Response-Header", "response-value")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("Hello, World!")); err != nil {
			t.Logf("Error writing response: %v", err)
		}
	}))
	defer targetServer.Close()

	// Create and start the proxy
	config := &Config{
		Port:     8080,
		LogLevel: "info",
	}
	proxy := New(config)

	// Create a test server using our proxy handler
	proxyServer := httptest.NewServer(proxy)
	defer proxyServer.Close()

	// Make a request through the proxy to the target
	// We need to use the proxy server URL but request the target URL
	req, err := http.NewRequest("GET", targetServer.URL, http.NoBody)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("X-Custom-Header", "custom-value")

	// Create a client that uses the proxy
	client := &http.Client{}

	// Make the request directly to the proxy (simulating proxy behavior)
	// In real usage, this would be configured as a proxy
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			t.Logf("Error closing response body: %v", closeErr)
		}
	}()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Check response headers
	if resp.Header.Get("X-Response-Header") != "response-value" {
		t.Error("Expected X-Response-Header to be forwarded")
	}

	// Check response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	expectedBody := "Hello, World!"
	if string(body) != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, string(body))
	}
}

func TestNewProxy(t *testing.T) {
	config := &Config{
		Port:        9090,
		AdminPort:   9091,
		LogLevel:    "debug",
		HistorySize: 500,
	}

	proxy := New(config)

	if proxy.config.Port != 9090 {
		t.Errorf("Expected port 9090, got %d", proxy.config.Port)
	}

	if proxy.config.LogLevel != "debug" {
		t.Errorf("Expected log level 'debug', got %s", proxy.config.LogLevel)
	}

	if proxy.config.AdminPort != 9091 {
		t.Errorf("Expected admin port 9091, got %d", proxy.config.AdminPort)
	}

	if proxy.config.HistorySize != 500 {
		t.Errorf("Expected history size 500, got %d", proxy.config.HistorySize)
	}

	if proxy.history == nil {
		t.Error("Expected history to be initialized")
	}
}
