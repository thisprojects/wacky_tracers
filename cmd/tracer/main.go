package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/thisprojects/wacky_tracers/pkg/config"
	"github.com/thisprojects/wacky_tracers/pkg/tracer"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "/etc/wacky_tracers/tracer.yaml", "Path to configuration file")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadTracerConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize tracer
	tracerInstance, err := tracer.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize tracer: %v", err)
	}

	// Setup context and graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start tracer
	go func() {
		if err := tracerInstance.Start(ctx); err != nil {
			log.Fatalf("Tracer failed: %v", err)
		}
	}()

	log.Println("Wacky Tracer started")

	// Wait for signal
	<-sigChan
	log.Println("Shutting down gracefully...")

	// Shutdown
	if err := tracerInstance.Shutdown(ctx); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}
}