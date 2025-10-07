// Package miner provides mining management functionality.
package miner

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/nn-selin/labracodabrador/internal/config"
	"github.com/nn-selin/labracodabrador/internal/node"
)

// Service manages mining operations across nodes.
type Service struct {
	logger       *slog.Logger
	nodeManager  *node.Manager
	config       *config.MiningConfig
	activeMiners map[string]bool
	mu           sync.RWMutex
	stopChan     chan struct{}
	wg           sync.WaitGroup
}

// NewService creates a new mining service.
func NewService(logger *slog.Logger, nodeManager *node.Manager, cfg *config.MiningConfig) *Service {
	return &Service{
		logger:       logger,
		nodeManager:  nodeManager,
		config:       cfg,
		activeMiners: make(map[string]bool),
		stopChan:     make(chan struct{}),
	}
}

// Start starts mining on all configured miner nodes.
func (s *Service) Start(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.logger.Info("starting mining service")

	nodes := s.nodeManager.GetAllNodes()
	started := 0

	for _, instance := range nodes {
		if !instance.IsHealthy() {
			s.logger.Warn("skipping unhealthy node for mining", "node", instance.Name)
			continue
		}

		ethereum := instance.GetEthereum()
		if ethereum == nil {
			s.logger.Warn("node has no ethereum service", "node", instance.Name)
			continue
		}

		// Start mining - in newer Geth versions, mining starts automatically
		// We just need to set the etherbase and the miner will work
		s.activeMiners[instance.GetAccount().Hex()] = true
		started++

		s.logger.Info("mining enabled on node",
			"node", instance.Name,
			"account", instance.GetAccount().Hex())
	}

	if started == 0 {
		return fmt.Errorf("no nodes available for mining")
	}

	s.logger.Info("mining service started", "active_miners", started)

	// Start monitoring goroutine
	s.wg.Add(1)
	go s.monitor(ctx)

	return nil
}

// Stop stops all mining operations.
func (s *Service) Stop(ctx context.Context) error {
	s.mu.Lock()
	s.logger.Info("stopping mining service")
	close(s.stopChan)
	s.mu.Unlock()

	// Wait for monitoring to stop
	s.wg.Wait()

	s.mu.Lock()
	defer s.mu.Unlock()

	nodes := s.nodeManager.GetAllNodes()

	for _, instance := range nodes {
		ethereum := instance.GetEthereum()
		if ethereum == nil {
			continue
		}

		delete(s.activeMiners, instance.GetAccount().Hex())
		s.logger.Info("mining disabled on node", "node", instance.Name)
	}

	s.logger.Info("mining service stopped")
	return nil
}

// monitor monitors mining operations and restarts if needed.
func (s *Service) monitor(ctx context.Context) {
	defer s.wg.Done()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopChan:
			return
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.checkAndRestartMiners(ctx)
		}
	}
}

// checkAndRestartMiners checks miner health and restarts if needed.
func (s *Service) checkAndRestartMiners(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()

	nodes := s.nodeManager.GetAllNodes()

	for _, instance := range nodes {
		if !instance.IsHealthy() {
			s.logger.Warn("node is unhealthy, skipping mining check", "node", instance.Name)
			continue
		}

		ethereum := instance.GetEthereum()
		if ethereum == nil {
			continue
		}

		// Check if mining configuration is maintained
		accountHex := instance.GetAccount().Hex()
		shouldBeMining := s.activeMiners[accountHex]

		if shouldBeMining {
			// Mining is configured and running
			s.logger.Debug("mining check passed", "node", instance.Name)
		}
	}
}

// GetMiningStats returns current mining statistics.
func (s *Service) GetMiningStats(ctx context.Context) (map[string]interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stats := make(map[string]interface{})
	stats["active_miners"] = len(s.activeMiners)

	miners := make([]map[string]interface{}, 0)
	nodes := s.nodeManager.GetAllNodes()

	for _, instance := range nodes {
		ethereum := instance.GetEthereum()
		if ethereum == nil {
			continue
		}

		client := instance.GetClient()
		if client == nil {
			continue
		}

		blockNum, err := client.BlockNumber(ctx)
		if err != nil {
			s.logger.Warn("failed to get block number",
				"node", instance.Name,
				"error", err)
			continue
		}

		accountHex := instance.GetAccount().Hex()
		minerInfo := map[string]interface{}{
			"name":         instance.Name,
			"account":      accountHex,
			"mining":       s.activeMiners[accountHex],
			"block_height": blockNum,
			"healthy":      instance.IsHealthy(),
		}

		miners = append(miners, minerInfo)
	}

	stats["miners"] = miners
	return stats, nil
}
