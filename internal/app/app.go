package app

import (
	"fmt"
	"log"
	"tx-parser/internal/api"
	"tx-parser/internal/config"
	"tx-parser/internal/interfaces"
	"tx-parser/internal/parser"
	"tx-parser/internal/rpc"
	"tx-parser/internal/storage"
	"tx-parser/pkg/logger"
)

type App struct {
	apiServer *api.Server
	parser    interfaces.Parser
	config    *config.Config
	log       *logger.Logger
}

func NewApp(configPath string) (*App, error) {
	// Load configuration
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
		return nil, err
	}

	// Initialize logger
	log := logger.GetLogger(cfg.Logging.Level)

	// Initialize storage
	storage := storage.NewMemoryStorage()

	// Initialize Ethereum RPC client
	rpcClient := rpc.NewClient(cfg.Server.Ethrpc, log)

	// Initialize parser
	ethParser := parser.NewEthParser(rpcClient, storage, log)

	// Initialize API server
	apiServer := api.NewServer(ethParser, storage, log)

	return &App{
		apiServer: apiServer,
		parser:    ethParser,
		config:    cfg,
		log:       log,
	}, nil
}

func (a *App) Run() error {
	serverAddr := fmt.Sprintf("%s%s", a.config.Server.Host, a.config.Server.Port)
	a.log.Info.Printf("Starting API server on %s...", serverAddr)
	return a.apiServer.Start(serverAddr) // Start the API server with the configured address
}
