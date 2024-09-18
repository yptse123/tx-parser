package app

import (
	"testing"

	"tx-parser/internal/api"
	"tx-parser/internal/config"
	"tx-parser/internal/interfaces"
	"tx-parser/internal/parser"
	"tx-parser/internal/rpc"
	"tx-parser/internal/storage"
	"tx-parser/pkg/logger"

	"github.com/stretchr/testify/assert"
)

// mockConfig is a mock implementation of the config loader
func mockConfig() *config.Config {
	return &config.Config{
		Server: config.ServerConfig{
			Host:   "localhost",
			Port:   ":8088",
			Ethrpc: "https://ethereum-rpc.publicnode.com",
		},
		Logging: config.LoggingConfig{
			Level: "debug",
		},
	}
}

// mockRPCClient is a mock implementation of the RPC client
type mockRPCClient struct{}

func (m *mockRPCClient) FetchCurrentBlock() (int, error) {
	return 123456, nil
}

func (m *mockRPCClient) FetchBlockByNumber(blockNumber int) (*rpc.Block, error) {
	return &rpc.Block{
		Transactions: []interfaces.Transaction{
			{From: "0xFrom", To: "0xTo", Value: "100", Hash: "0x123"},
		},
	}, nil
}

func TestNewApp(t *testing.T) {
	// Mock the configuration
	cfg := mockConfig()

	// Initialize the logger
	log := logger.GetLogger(cfg.Logging.Level)

	// Initialize mock storage
	storage := storage.NewMemoryStorage()

	// Initialize mock Ethereum RPC client
	rpcClient := &mockRPCClient{}

	// Initialize the parser
	ethParser := parser.NewEthParser(rpcClient, storage, log)

	// Initialize API server
	apiServer := api.NewServer(ethParser, storage, log)

	// Create a new app instance
	app := &App{
		apiServer: apiServer,
		parser:    ethParser,
		config:    cfg,
		log:       log,
	}

	// Assert app initialization
	assert.NotNil(t, app, "App should be initialized")
	assert.NotNil(t, app.apiServer, "API server should be initialized")
	assert.NotNil(t, app.parser, "Parser should be initialized")
	assert.NotNil(t, app.config, "Config should be initialized")
	assert.NotNil(t, app.log, "Logger should be initialized")
}
