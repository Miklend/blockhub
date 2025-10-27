package fabricClient

import (
	"fmt"
	"lib/clients/node"
	"lib/clients/node/alchemy"
	"lib/models"
	"lib/utils/logging"
)

const (
	alchemyType = "alchemy"
)

func NewProvider(cfg models.Provider, logger *logging.Logger) (node.Provider, error) {
	switch cfg.ProviderType {
	case alchemyType:
		return alchemy.NewAlchemyClient(cfg, logger)
	default:
		return nil, fmt.Errorf("not found provider for client type: %s", cfg.ProviderType)
	}
}

func NewProviderPool(cfg models.Provider, logger *logging.Logger) ([]node.Provider, error) {
	if cfg.NumClients <= 0 {
		return nil, fmt.Errorf("Client Num is invalid")
	}
	pool := make([]node.Provider, cfg.NumClients)

	logger.Infof("Initializing %d Clients", cfg.NumClients)

	for i := 0; i < cfg.NumClients; i++ {
		provider, err := NewProvider(cfg, logger)
		if err != nil {
			logger.Warn("Failed to create client %d, %w", i, err)
		}
		pool[i] = provider
	}

	logger.Infof("Successfully created %d clients", cfg.NumClients)
	return pool, nil
}
