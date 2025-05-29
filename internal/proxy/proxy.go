package proxy

import (
	"fmt"
	"net/http"
	"time"
)

// Config holds the proxy configuration
type Config struct {
	Port     int
	LogLevel string
	Metrics  bool
	Health   bool
}

// Proxy represents the HTTP proxy server
type Proxy struct {
	config     *Config
	server     *http.Server
	httpClient *http.Client
}

// New creates a new Proxy instance
func New(config *Config) *Proxy {
	return &Proxy{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Start starts the proxy server
func (p *Proxy) Start() error {
	if p.server == nil {
		return fmt.Errorf("server not initialized")
	}

	return p.server.ListenAndServe()
}

// Stop stops the proxy server
func (p *Proxy) Stop() error {
	if p.server == nil {
		return fmt.Errorf("server not initialized")
	}

	return nil
}
