package proxy

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
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

// ServeHTTP implements the http.Handler interface for the proxy
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// For CONNECT method (HTTPS tunneling)
	if r.Method == http.MethodConnect {
		p.handleConnect(w, r)
		return
	}

	// For regular HTTP requests
	p.handleHTTP(w, r)
}

// handleHTTP handles regular HTTP requests
func (p *Proxy) handleHTTP(w http.ResponseWriter, r *http.Request) {
	// Create a new request to the target server
	targetURL, err := url.Parse(r.URL.String())
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	// Create the proxied request
	proxyReq, err := http.NewRequest(r.Method, targetURL.String(), r.Body)
	if err != nil {
		http.Error(w, "Failed to create proxy request", http.StatusInternalServerError)
		return
	}

	// Copy headers from original request
	for key, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(key, value)
		}
	}

	// Make the request to the target server
	resp, err := p.httpClient.Do(proxyReq)
	if err != nil {
		http.Error(w, "Failed to proxy request", http.StatusBadGateway)
		return
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Printf("Error closing response body: %v", closeErr)
		}
	}()

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Copy status code
	w.WriteHeader(resp.StatusCode)

	// Copy response body
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Printf("Error copying response body: %v", err)
	}
}

// handleConnect handles CONNECT method for HTTPS tunneling
func (p *Proxy) handleConnect(w http.ResponseWriter, r *http.Request) {
	// This is a simplified CONNECT handler
	// In a production proxy, you'd implement proper tunneling
	dest, err := net.Dial("tcp", r.Host)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer func() {
		if closeErr := dest.Close(); closeErr != nil {
			log.Printf("Error closing destination connection: %v", closeErr)
		}
	}()

	w.WriteHeader(http.StatusOK)

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}

	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer func() {
		if closeErr := clientConn.Close(); closeErr != nil {
			log.Printf("Error closing client connection: %v", closeErr)
		}
	}()

	// Start copying data between client and destination
	go func() {
		if _, err := io.Copy(dest, clientConn); err != nil {
			log.Printf("Error copying from client to destination: %v", err)
		}
	}()

	if _, err := io.Copy(clientConn, dest); err != nil {
		log.Printf("Error copying from destination to client: %v", err)
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

	return p.server.Shutdown(context.TODO())
}
