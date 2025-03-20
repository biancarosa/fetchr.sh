package proxy

import (
	"fmt"
	"io"
	"log"
	"net"
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
	p.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", p.config.Port),
		Handler: http.HandlerFunc(p.handleRequest),
	}

	log.Printf("Starting proxy server on port %d", p.config.Port)
	return p.server.ListenAndServe()
}

// Stop stops the proxy server
func (p *Proxy) Stop() error {
	if p.server != nil {
		return p.server.Close()
	}
	return nil
}

// handleRequest handles incoming HTTP requests
func (p *Proxy) handleRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: %s %s", r.Method, r.URL)

	if r.Method == http.MethodConnect {
		p.handleHTTPS(w, r)
		return
	}

	p.handleHTTP(w, r)
}

// handleHTTP handles plain HTTP requests
func (p *Proxy) handleHTTP(w http.ResponseWriter, r *http.Request) {
	// Create a new request to forward
	outReq, err := http.NewRequest(r.Method, r.URL.String(), r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating request: %v", err), http.StatusInternalServerError)
		return
	}

	// Copy headers from original request
	for key, values := range r.Header {
		for _, value := range values {
			outReq.Header.Add(key, value)
		}
	}

	// Forward the request
	resp, err := p.httpClient.Do(outReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error forwarding request: %v", err), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Set response status code
	w.WriteHeader(resp.StatusCode)

	// Copy response body
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Printf("Error copying response body: %v", err)
	}

	log.Printf("HTTP request completed: %s %s -> %d", r.Method, r.URL, resp.StatusCode)
}

// handleHTTPS handles HTTPS CONNECT requests
func (p *Proxy) handleHTTPS(w http.ResponseWriter, r *http.Request) {
	// Connect to the target server
	targetConn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error connecting to target: %v", err), http.StatusBadGateway)
		return
	}

	// Hijack the client connection
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}

	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error hijacking connection: %v", err), http.StatusInternalServerError)
		return
	}

	// Send 200 OK to indicate tunnel established
	_, err = clientConn.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
	if err != nil {
		clientConn.Close()
		targetConn.Close()
		return
	}

	// Start bidirectional copy
	go func() {
		defer clientConn.Close()
		defer targetConn.Close()
		io.Copy(targetConn, clientConn)
	}()

	go func() {
		defer clientConn.Close()
		defer targetConn.Close()
		io.Copy(clientConn, targetConn)
	}()

	log.Printf("HTTPS tunnel established: %s", r.Host)
}
