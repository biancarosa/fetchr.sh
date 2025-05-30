package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/biancarosa/fetchr.sh/internal/proxy"
)

func main() {
	// Parse command
	if len(os.Args) < 2 {
		log.Fatal("Please specify a command: serve or request")
	}

	command := os.Args[1]
	os.Args = os.Args[1:] // Remove command from args for flag parsing

	switch command {
	case "serve":
		runServe()
	case "request":
		if err := runRequest(); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("Unknown command: %s. Use 'serve' or 'request'", command)
	}
}

func runServe() {
	// Parse command line flags
	port := flag.Int("port", 8080, "Port to listen on")
	adminPort := flag.Int("admin-port", 0, "Admin port for health checks and metrics (0 to disable)")
	logLevel := flag.String("log-level", "info", "Logging level (debug, info, warn, error)")
	flag.Parse()

	// Create proxy configuration
	config := &proxy.Config{
		Port:      *port,
		AdminPort: *adminPort,
		LogLevel:  *logLevel,
	}

	// Create and start proxy server
	proxyServer := proxy.New(config)

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Add debug logging for when we're about to start
	if *logLevel == "debug" {
		log.Printf("Starting proxy server on port %d", config.Port)
		if config.AdminPort > 0 {
			log.Printf("Admin endpoints will be available on port %d", config.AdminPort)
		}
	}

	go func() {
		if err := proxyServer.Start(); err != nil {
			log.Fatalf("Failed to start proxy server: %v", err)
		}
	}()

	log.Printf("Proxy server started on port %d", config.Port)
	if config.AdminPort > 0 {
		log.Printf("Admin server started on port %d (health: /healthz, metrics: /metrics)", config.AdminPort)
	}

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutting down proxy server...")

	if err := proxyServer.Stop(); err != nil {
		log.Printf("Error stopping proxy server: %v", err)
	}
}
