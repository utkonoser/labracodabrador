// Package config provides configuration management for the blockchain network.
package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the main configuration structure.
type Config struct {
	Network     NetworkConfig     `yaml:"network"`
	Nodes       []NodeConfig      `yaml:"nodes"`
	Mining      MiningConfig      `yaml:"mining"`
	HealthCheck HealthCheckConfig `yaml:"healthcheck"`
}

// NetworkConfig contains network-level settings.
type NetworkConfig struct {
	ChainID   uint64 `yaml:"chain_id"`
	NetworkID uint64 `yaml:"network_id"`
}

// NodeConfig represents configuration for a single node.
type NodeConfig struct {
	Name          string `yaml:"name"`
	Port          int    `yaml:"port"`
	RPCPort       int    `yaml:"rpc_port"`
	WSPort        int    `yaml:"ws_port"`
	DiscoveryPort int    `yaml:"discovery_port"`
	Miner         bool   `yaml:"miner"`
	Etherbase     string `yaml:"etherbase"`
}

// MiningConfig contains mining-related settings.
type MiningConfig struct {
	Threads  int    `yaml:"threads"`
	GasPrice uint64 `yaml:"gas_price"`
}

// HealthCheckConfig contains health check settings.
type HealthCheckConfig struct {
	Interval string `yaml:"interval"`
	Timeout  string `yaml:"timeout"`
	Retries  int    `yaml:"retries"`
}

// IntervalDuration returns the interval as a time.Duration.
func (h *HealthCheckConfig) IntervalDuration() (time.Duration, error) {
	return time.ParseDuration(h.Interval)
}

// TimeoutDuration returns the timeout as a time.Duration.
func (h *HealthCheckConfig) TimeoutDuration() (time.Duration, error) {
	return time.ParseDuration(h.Timeout)
}

// Load reads and parses the configuration file.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}

// Validate checks if the configuration is valid.
func (c *Config) Validate() error {
	if c.Network.ChainID == 0 {
		return fmt.Errorf("chain_id must be set")
	}

	if len(c.Nodes) == 0 {
		return fmt.Errorf("at least one node must be configured")
	}

	for i, node := range c.Nodes {
		if node.Name == "" {
			return fmt.Errorf("node %d: name must be set", i)
		}
		if node.RPCPort == 0 {
			return fmt.Errorf("node %s: rpc_port must be set", node.Name)
		}
	}

	return nil
}
