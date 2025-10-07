// Package main is the entry point for the labracodabrador blockchain application.
package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nn-selin/labracodabrador/internal/config"
	"github.com/nn-selin/labracodabrador/internal/miner"
	"github.com/nn-selin/labracodabrador/internal/node"
)

const (
	defaultConfigPath  = "config.yaml"
	defaultGenesisPath = "genesis.json"
)

var (
	configPath  = flag.String("config", defaultConfigPath, "Path to configuration file")
	genesisPath = flag.String("genesis", defaultGenesisPath, "Path to genesis file")
	logLevel    = flag.String("log-level", "info", "Log level (debug, info, warn, error)")
)

func main() {
	flag.Parse()

	// Setup logger
	logger := setupLogger(*logLevel)
	logger.Info("starting labracodabrador blockchain",
		"config", *configPath,
		"genesis", *genesisPath)

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		logger.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	if err := cfg.Validate(); err != nil {
		logger.Error("invalid configuration", "error", err)
		os.Exit(1)
	}

	logger.Info("configuration loaded successfully",
		"chain_id", cfg.Network.ChainID,
		"nodes", len(cfg.Nodes))

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create node manager
	nodeManager := node.NewManager(logger, cfg, *genesisPath)

	// Start all nodes
	if err := nodeManager.StartAll(ctx); err != nil {
		logger.Error("failed to start nodes", "error", err)
		os.Exit(1)
	}

	// Wait a bit for nodes to initialize
	logger.Info("waiting for nodes to initialize...")
	time.Sleep(5 * time.Second)

	// Connect nodes as peers
	if err := nodeManager.ConnectPeers(ctx); err != nil {
		logger.Error("failed to connect peers", "error", err)
	}

	// Wait for peer connections
	time.Sleep(3 * time.Second)

	// Create and start mining service
	miningService := miner.NewService(logger, nodeManager, &cfg.Mining)
	if err := miningService.Start(ctx); err != nil {
		logger.Error("failed to start mining service", "error", err)
		// Don't exit, continue running nodes
	}

	// Start health check routine
	go runHealthChecks(ctx, logger, nodeManager, cfg)

	// Print node information
	printNodeInfo(logger, nodeManager)

	logger.Info("labracodabrador blockchain is running. Press Ctrl+C to stop.")

	// Wait for shutdown signal
	<-sigChan
	logger.Info("shutdown signal received, stopping...")

	// Create shutdown context with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Stop mining
	if err := miningService.Stop(shutdownCtx); err != nil {
		logger.Error("error stopping mining service", "error", err)
	}

	// Stop all nodes
	if err := nodeManager.StopAll(shutdownCtx); err != nil {
		logger.Error("error stopping nodes", "error", err)
	}

	logger.Info("shutdown complete")
}

// setupLogger creates and configures the logger.
func setupLogger(level string) *slog.Logger {
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: false,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	return slog.New(handler)
}

// runHealthChecks periodically checks node health.
func runHealthChecks(ctx context.Context, logger *slog.Logger, manager *node.Manager, cfg *config.Config) {
	interval, err := cfg.HealthCheck.IntervalDuration()
	if err != nil {
		logger.Error("invalid health check interval", "error", err)
		interval = 10 * time.Second
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := manager.HealthCheck(ctx); err != nil {
				logger.Error("health check failed", "error", err)
			}
		}
	}
}

// printNodeInfo prints information about running nodes.
func printNodeInfo(logger *slog.Logger, manager *node.Manager) {
	nodes := manager.GetAllNodes()

	separator := repeatString("=", 80)

	fmt.Println("\n" + separator)
	fmt.Println("BLOCKCHAIN NODES RUNNING")
	fmt.Println(separator)

	for _, n := range nodes {
		fmt.Printf("\nNode: %s\n", n.Name)
		fmt.Printf("  RPC:       http://localhost:%d\n", n.Config.RPCPort)
		fmt.Printf("  WebSocket: ws://localhost:%d\n", n.Config.WSPort)
		fmt.Printf("  Account:   %s\n", n.GetAccount().Hex())
		fmt.Println("  Status:    ✓ Running")
	}

	fmt.Println("\n" + separator)
	fmt.Println("METAMASK CONFIGURATION")
	fmt.Println(separator)
	fmt.Println("  Network Name: Labracodabrador Network")
	fmt.Printf("  RPC URL:      http://localhost:%d\n", nodes[0].Config.RPCPort)
	fmt.Println("  Chain ID:     1337")
	fmt.Println("  Currency:     ETH")
	fmt.Println(separator + "\n")
}

// repeatString is a helper function to repeat strings.
func repeatString(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
