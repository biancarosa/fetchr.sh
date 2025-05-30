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
	logLevel := flag.String("log-level", "info", "Logging level (debug, info, warn, error)")
	metrics := flag.Bool("metrics", false, "Enable metrics endpoint")
	health := flag.Bool("health", false, "Enable health check endpoint")
	flag.Parse()

	// Create proxy configuration
	config := &proxy.Config{
		Port:     *port,
		LogLevel: *logLevel,
		Metrics:  *metrics,
		Health:   *health,
	}

	// Create and start proxy server
	proxyServer := proxy.New(config)

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Add debug logging for when we're about to start
	if *logLevel == "debug" {
		log.Printf("Starting proxy server on port %d", config.Port)
	}

	go func() {
		if err := proxyServer.Start(); err != nil {
			log.Fatalf("Failed to start proxy server: %v", err)
		}
	}()

	log.Printf("Proxy server started on port %d", config.Port)

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutting down proxy server...")

	if err := proxyServer.Stop(); err != nil {
		log.Printf("Error stopping proxy server: %v", err)
	}
}
