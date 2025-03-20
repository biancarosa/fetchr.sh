package api

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// RequestConfig holds the configuration for an API request
type RequestConfig struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    io.Reader
	Timeout time.Duration
}

// Response represents the API response
type Response struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
}

// MakeRequest makes a single API request through the proxy
func MakeRequest(proxyURL string, config RequestConfig) (*Response, error) {
	// Create proxy URL
	proxy, err := url.Parse(proxyURL)
	if err != nil {
		return nil, fmt.Errorf("invalid proxy URL: %v", err)
	}

	// Create HTTP client with proxy
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxy),
		},
		Timeout: config.Timeout,
	}

	// Create request
	req, err := http.NewRequest(config.Method, config.URL, config.Body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Add headers
	for key, value := range config.Headers {
		req.Header.Add(key, value)
	}

	// Make request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       body,
	}, nil
}
