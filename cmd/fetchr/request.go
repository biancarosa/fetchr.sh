package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/biancarosa/fetchr.sh/internal/api"
	"github.com/biancarosa/fetchr.sh/internal/proxy"
)

func runRequest() error {
	// Parse command line flags
	method := flag.String("method", "GET", "HTTP method")
	url := flag.String("url", "", "Target URL (required)")
	port := flag.Int("port", 8080, "Proxy port")
	timeout := flag.Duration("timeout", 30*time.Second, "Request timeout")
	flag.Parse()

	if *url == "" {
		return fmt.Errorf("--url is required")
	}

	// Create proxy configuration
	config := &proxy.Config{
		Port:     *port,
		LogLevel: "info",
	}

	// Create and start proxy server
	proxyServer := proxy.New(config)

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start proxy server in a goroutine
	go func() {
		if err := proxyServer.Start(); err != nil {
			log.Printf("Proxy server error: %v", err)
		}
	}()

	// Wait a moment for the proxy to start
	time.Sleep(100 * time.Millisecond)

	// Make the request through the proxy
	proxyURL := fmt.Sprintf("http://localhost:%d", *port)
	reqConfig := api.RequestConfig{
		Method:  *method,
		URL:     *url,
		Headers: make(map[string]string),
		Timeout: *timeout,
	}

	// Add default headers
	reqConfig.Headers["User-Agent"] = "fetchr.sh/1.0"

	// Make the request
	resp, err := api.MakeRequest(proxyURL, reqConfig)
	if err != nil {
		proxyServer.Stop()
		return fmt.Errorf("request failed: %v", err)
	}

	// Print response
	fmt.Printf("Status: %d\n", resp.StatusCode)
	fmt.Println("\nHeaders:")
	for key, values := range resp.Headers {
		for _, value := range values {
			fmt.Printf("%s: %s\n", key, value)
		}
	}
	fmt.Println("\nBody:")
	fmt.Println(string(resp.Body))

	// Stop the proxy server
	return proxyServer.Stop()
}
