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
	Port      int
	AdminPort int
	LogLevel  string
}

// Proxy represents the HTTP proxy server
type Proxy struct {
	config      *Config
	server      *http.Server
	adminServer *http.Server
	httpClient  *http.Client
}

// New creates a new Proxy instance
func New(config *Config) *Proxy {
	proxy := &Proxy{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	// Initialize the main HTTP proxy server
	proxy.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: proxy,
	}

	// Initialize the admin server if admin port is specified
	if config.AdminPort > 0 {
		adminMux := http.NewServeMux()

		// Always enable both health and metrics when admin port is specified
		adminMux.HandleFunc("/healthz", proxy.handleHealth)
		adminMux.HandleFunc("/metrics", proxy.handleMetrics)

		proxy.adminServer = &http.Server{
			Addr:    fmt.Sprintf(":%d", config.AdminPort),
			Handler: adminMux,
		}
	}

	return proxy
}

// ServeHTTP implements the http.Handler interface for the proxy
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Debug logging for received requests
	if p.config.LogLevel == "debug" {
		log.Printf("Received request: %s %s", r.Method, r.URL.String())
	}

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

	// Debug logging for completed requests
	if p.config.LogLevel == "debug" {
		log.Printf("HTTP request completed: %s %s -> %d", r.Method, r.URL.String(), resp.StatusCode)
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

// handleHealth handles health check requests
func (p *Proxy) handleHealth(w http.ResponseWriter, r *http.Request) {
	// Add CORS headers to allow requests from the dashboard
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight requests
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(`{"status":"healthy","proxy":"fetchr.sh"}`)); err != nil {
		log.Printf("Error writing health response: %v", err)
	}
}

// handleMetrics handles metrics requests
func (p *Proxy) handleMetrics(w http.ResponseWriter, r *http.Request) {
	// Add CORS headers to allow requests from the dashboard
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight requests
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	// Simple metrics for now - can be expanded later
	metrics := `# HELP fetchr_requests_total Total number of requests handled
# TYPE fetchr_requests_total counter
fetchr_requests_total 0

# HELP fetchr_proxy_status Status of the proxy server
# TYPE fetchr_proxy_status gauge
fetchr_proxy_status 1
`
	if _, err := w.Write([]byte(metrics)); err != nil {
		log.Printf("Error writing metrics response: %v", err)
	}
}

// Start starts the proxy server and admin server (if configured)
func (p *Proxy) Start() error {
	if p.server == nil {
		return fmt.Errorf("server not initialized")
	}

	// Start admin server in background if configured
	if p.adminServer != nil {
		go func() {
			log.Printf("Starting admin server on port %d", p.config.AdminPort)
			if err := p.adminServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Printf("Admin server error: %v", err)
			}
		}()
	}

	log.Printf("Starting proxy server on port %d", p.config.Port)
	return p.server.ListenAndServe()
}

// Stop stops both the proxy server and admin server
func (p *Proxy) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var proxyErr, adminErr error

	if p.server != nil {
		proxyErr = p.server.Shutdown(ctx)
	}

	if p.adminServer != nil {
		adminErr = p.adminServer.Shutdown(ctx)
	}

	// Return the first error encountered
	if proxyErr != nil {
		return fmt.Errorf("proxy server shutdown error: %v", proxyErr)
	}
	if adminErr != nil {
		return fmt.Errorf("admin server shutdown error: %v", adminErr)
	}

	return nil
}
