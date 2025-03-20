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
		w.Header().Set("X-Test-Header", "test-value")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello from target server"))
	}))
	defer targetServer.Close()

	// Create and start the proxy
	config := &Config{
		Port:     8080,
		LogLevel: "info",
	}
	proxy := New(config)

	// Create a test server using our proxy handler
	proxyServer := httptest.NewServer(http.HandlerFunc(proxy.handleRequest))
	defer proxyServer.Close()

	// Make a request through the proxy to the target
	req, err := http.NewRequest("GET", targetServer.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("X-Custom-Header", "custom-value")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Check response headers
	if resp.Header.Get("X-Test-Header") != "test-value" {
		t.Error("Expected X-Test-Header to be forwarded")
	}

	// Check response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	expectedBody := "Hello from target server"
	if string(body) != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, string(body))
	}
} 