// Package node provides functionality for managing Geth nodes.
package node

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/eth/downloader"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/nn-selin/labracodabrador/internal/config"
)

// Manager handles the lifecycle of multiple Geth nodes.
type Manager struct {
	logger      *slog.Logger
	nodes       map[string]*NodeInstance
	genesisPath string
	config      *config.Config
	mu          sync.RWMutex
}

// NodeInstance represents a single Geth node instance.
type NodeInstance struct {
	Name     string
	node     *node.Node
	ethereum *eth.Ethereum
	client   *ethclient.Client
	Config   config.NodeConfig
	dataDir  string
	account  common.Address
	healthy  bool
	mu       sync.RWMutex
}

// NewManager creates a new node manager.
func NewManager(logger *slog.Logger, cfg *config.Config, genesisPath string) *Manager {
	return &Manager{
		logger:      logger,
		nodes:       make(map[string]*NodeInstance),
		genesisPath: genesisPath,
		config:      cfg,
	}
}

// StartAll starts all configured nodes.
func (m *Manager) StartAll(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.logger.Info("starting all nodes", "count", len(m.config.Nodes))

	var wg sync.WaitGroup
	errChan := make(chan error, len(m.config.Nodes))

	for i := range m.config.Nodes {
		wg.Add(1)
		go func(nodeCfg config.NodeConfig) {
			defer wg.Done()
			if err := m.startNode(ctx, nodeCfg); err != nil {
				errChan <- fmt.Errorf("failed to start node %s: %w", nodeCfg.Name, err)
			}
		}(m.config.Nodes[i])
	}

	wg.Wait()
	close(errChan)

	// Collect any errors
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	m.logger.Info("all nodes started successfully")
	return nil
}

// startNode starts a single node.
func (m *Manager) startNode(ctx context.Context, cfg config.NodeConfig) error {
	m.logger.Info("starting node", "name", cfg.Name, "rpc_port", cfg.RPCPort)

	dataDir := filepath.Join("data", cfg.Name)
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	// Create node configuration
	nodeCfg := &node.Config{
		Name:    cfg.Name,
		DataDir: dataDir,
		P2P: p2p.Config{
			MaxPeers:    50,
			ListenAddr:  fmt.Sprintf(":%d", cfg.DiscoveryPort),
			NoDiscovery: false,
		},
		HTTPHost:         "0.0.0.0",
		HTTPPort:         cfg.RPCPort,
		HTTPModules:      []string{"eth", "net", "web3", "personal", "admin", "miner", "txpool"},
		HTTPCors:         []string{"*"},
		HTTPVirtualHosts: []string{"*"},
		WSHost:           "0.0.0.0",
		WSPort:           cfg.WSPort,
		WSModules:        []string{"eth", "net", "web3", "personal", "admin", "miner", "txpool"},
		WSOrigins:        []string{"*"},
	}

	// Create node
	stack, err := node.New(nodeCfg)
	if err != nil {
		return fmt.Errorf("failed to create node: %w", err)
	}

	// Initialize genesis if needed
	if err := m.initGenesis(dataDir); err != nil {
		return fmt.Errorf("failed to initialize genesis: %w", err)
	}

	// Create or load account
	account, err := m.getOrCreateAccount(stack, cfg.Name)
	if err != nil {
		return fmt.Errorf("failed to get account: %w", err)
	}

	// Configure Ethereum service
	ethCfg := ethconfig.Defaults
	ethCfg.Genesis = m.loadGenesis()
	ethCfg.NetworkId = m.config.Network.NetworkID
	ethCfg.SyncMode = downloader.FullSync
	ethCfg.Miner.Etherbase = account
	ethCfg.Miner.GasPrice = nil

	// Register Ethereum service
	ethereum, err := eth.New(stack, &ethCfg)
	if err != nil {
		return fmt.Errorf("failed to create ethereum service: %w", err)
	}

	// Start the node
	if err := stack.Start(); err != nil {
		return fmt.Errorf("failed to start node: %w", err)
	}

	// Connect to the node via RPC
	client, err := ethclient.Dial(fmt.Sprintf("http://localhost:%d", cfg.RPCPort))
	if err != nil {
		stack.Close()
		return fmt.Errorf("failed to connect to node: %w", err)
	}

	instance := &NodeInstance{
		Name:     cfg.Name,
		node:     stack,
		ethereum: ethereum,
		client:   client,
		Config:   cfg,
		dataDir:  dataDir,
		account:  account,
		healthy:  true,
	}

	m.nodes[cfg.Name] = instance
	m.logger.Info("node started successfully",
		"name", cfg.Name,
		"rpc", fmt.Sprintf("http://localhost:%d", cfg.RPCPort),
		"ws", fmt.Sprintf("ws://localhost:%d", cfg.WSPort),
		"account", account.Hex())

	return nil
}

// initGenesis initializes the genesis block if needed.
func (m *Manager) initGenesis(dataDir string) error {
	chainDataDir := filepath.Join(dataDir, "geth", "chaindata")
	if _, err := os.Stat(chainDataDir); err == nil {
		// Genesis already initialized
		return nil
	}

	// Genesis will be initialized automatically when the node starts
	m.logger.Info("genesis will be initialized on first run", "datadir", dataDir)
	return nil
}

// loadGenesis loads the genesis configuration.
func (m *Manager) loadGenesis() *core.Genesis {
	data, err := os.ReadFile(m.genesisPath)
	if err != nil {
		m.logger.Error("failed to read genesis file", "error", err)
		return core.DefaultGenesisBlock()
	}

	genesis := new(core.Genesis)
	if err := genesis.UnmarshalJSON(data); err != nil {
		m.logger.Error("failed to parse genesis file", "error", err)
		return core.DefaultGenesisBlock()
	}

	return genesis
}

// getOrCreateAccount gets or creates an account for the node.
func (m *Manager) getOrCreateAccount(stack *node.Node, nodeName string) (common.Address, error) {
	ks := keystore.NewKeyStore(
		filepath.Join(stack.DataDir(), "keystore"),
		keystore.StandardScryptN,
		keystore.StandardScryptP,
	)

	if len(ks.Accounts()) > 0 {
		return ks.Accounts()[0].Address, nil
	}

	account, err := ks.NewAccount("password")
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to create account: %w", err)
	}

	// Unlock account
	if err := ks.Unlock(account, "password"); err != nil {
		return common.Address{}, fmt.Errorf("failed to unlock account: %w", err)
	}

	m.logger.Info("created new account", "node", nodeName, "address", account.Address.Hex())
	return account.Address, nil
}

// StopAll stops all nodes gracefully.
func (m *Manager) StopAll(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.logger.Info("stopping all nodes", "count", len(m.nodes))

	var wg sync.WaitGroup
	for name, instance := range m.nodes {
		wg.Add(1)
		go func(n string, inst *NodeInstance) {
			defer wg.Done()
			if err := inst.Stop(); err != nil {
				m.logger.Error("failed to stop node", "name", n, "error", err)
			} else {
				m.logger.Info("node stopped", "name", n)
			}
		}(name, instance)
	}

	wg.Wait()
	m.nodes = make(map[string]*NodeInstance)
	m.logger.Info("all nodes stopped")
	return nil
}

// GetNode returns a node instance by name.
func (m *Manager) GetNode(name string) (*NodeInstance, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	node, ok := m.nodes[name]
	return node, ok
}

// GetAllNodes returns all node instances.
func (m *Manager) GetAllNodes() []*NodeInstance {
	m.mu.RLock()
	defer m.mu.RUnlock()

	nodes := make([]*NodeInstance, 0, len(m.nodes))
	for _, node := range m.nodes {
		nodes = append(nodes, node)
	}
	return nodes
}

// HealthCheck performs health checks on all nodes.
func (m *Manager) HealthCheck(ctx context.Context) error {
	m.mu.RLock()
	nodes := make([]*NodeInstance, 0, len(m.nodes))
	for _, node := range m.nodes {
		nodes = append(nodes, node)
	}
	m.mu.RUnlock()

	for _, instance := range nodes {
		if err := instance.CheckHealth(ctx); err != nil {
			m.logger.Warn("node health check failed",
				"name", instance.Name,
				"error", err)
		}
	}

	return nil
}

// ConnectPeers connects all nodes to each other.
func (m *Manager) ConnectPeers(ctx context.Context) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	nodes := make([]*NodeInstance, 0, len(m.nodes))
	for _, node := range m.nodes {
		nodes = append(nodes, node)
	}

	if len(nodes) < 2 {
		return nil
	}

	// Connect each node to all other nodes
	for i := 0; i < len(nodes); i++ {
		for j := i + 1; j < len(nodes); j++ {
			if err := nodes[i].ConnectTo(nodes[j]); err != nil {
				m.logger.Warn("failed to connect nodes",
					"from", nodes[i].Name,
					"to", nodes[j].Name,
					"error", err)
			} else {
				m.logger.Info("nodes connected",
					"from", nodes[i].Name,
					"to", nodes[j].Name)
			}
		}
	}

	return nil
}

// Stop stops the node.
func (n *NodeInstance) Stop() error {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.client != nil {
		n.client.Close()
	}

	if n.node != nil {
		return n.node.Close()
	}

	return nil
}

// CheckHealth checks if the node is healthy.
func (n *NodeInstance) CheckHealth(ctx context.Context) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.client == nil {
		n.healthy = false
		return fmt.Errorf("client is nil")
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := n.client.BlockNumber(timeoutCtx)
	if err != nil {
		n.healthy = false
		return fmt.Errorf("failed to get block number: %w", err)
	}

	n.healthy = true
	return nil
}

// IsHealthy returns whether the node is healthy.
func (n *NodeInstance) IsHealthy() bool {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.healthy
}

// GetClient returns the Ethereum client.
func (n *NodeInstance) GetClient() *ethclient.Client {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.client
}

// GetEthereum returns the Ethereum service.
func (n *NodeInstance) GetEthereum() *eth.Ethereum {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.ethereum
}

// GetAccount returns the node's account address.
func (n *NodeInstance) GetAccount() common.Address {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.account
}

// ConnectTo connects this node to another node.
func (n *NodeInstance) ConnectTo(peer *NodeInstance) error {
	n.mu.RLock()
	myNode := n.node
	n.mu.RUnlock()

	peer.mu.RLock()
	peerNode := peer.node
	peer.mu.RUnlock()

	if myNode == nil || peerNode == nil {
		return fmt.Errorf("one of the nodes is not initialized")
	}

	// Get peer's enode URL
	var enodeURL string
	if peerNode.Server() != nil {
		enodeURL = peerNode.Server().Self().URLv4()
	} else {
		return fmt.Errorf("peer server not available")
	}

	// Parse enode
	node, err := enode.Parse(enode.ValidSchemes, enodeURL)
	if err != nil {
		return fmt.Errorf("failed to parse enode: %w", err)
	}

	// Add peer
	myNode.Server().AddPeer(node)
	return nil
}
