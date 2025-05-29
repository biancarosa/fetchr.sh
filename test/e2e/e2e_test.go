//go:build e2e

package e2e

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestE2E(t *testing.T) {
	t.Log("E2E test")
}

func TestProxyGoogleRequest(t *testing.T) {
	// Start the proxy server in the background
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Use a random available port
	port := getAvailablePort(t)

	// Get the absolute path to the fetchr binary in the project root
	binPath, err := filepath.Abs("../../fetchr")
	require.NoError(t, err, "Failed to get absolute path to fetchr binary")

	// Ensure the binary is executable
	err = os.Chmod(binPath, 0755)
	require.NoError(t, err, "Failed to make binary executable")

	// Start proxy server with the absolute path
	cmd := exec.CommandContext(ctx, binPath, "serve", "--port", port)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	require.NoError(t, err, "Failed to start proxy server")

	// Ensure the server has time to start
	time.Sleep(2 * time.Second)

	// Create an HTTP client that uses the proxy
	proxyURL, err := url.Parse("http://localhost:" + port)
	require.NoError(t, err)

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
		Timeout: 10 * time.Second,
	}

	// Make a request to Google through the proxy
	resp, err := client.Get("https://www.google.com")
	require.NoError(t, err, "Failed to make request through proxy")
	defer resp.Body.Close()

	// Verify the response status code
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected 200 OK response from Google")

	// Read and verify the response body contains expected Google content
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "Failed to read response body")

	assert.Contains(t, string(body), "<title", "Response should contain HTML title tag")
	assert.Contains(t, string(body), "Google", "Response should contain Google branding")
}

// getAvailablePort returns a random available port as a string
func getAvailablePort(t *testing.T) string {
	listener, err := net.Listen("tcp", ":0")
	require.NoError(t, err, "Failed to find available port")
	defer listener.Close()

	_, port, err := net.SplitHostPort(listener.Addr().String())
	require.NoError(t, err, "Failed to extract port number")

	return port
}
