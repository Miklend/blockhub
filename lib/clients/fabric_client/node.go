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

func NewProviderPool(cfgList []models.Provider, logger *logging.Logger) ([]node.Provider, error) {
	if len(cfgList) == 0 {
		return nil, fmt.Errorf("No providers configured")
	}

	var providerPool []node.Provider

	for i, cfg := range cfgList {
		if cfg.NumClients <= 0 {
			return nil, fmt.Errorf("num_clients must be greater than zero")

		}

		logger.Infof("Creating %d clients for key %d...", cfg.NumClients, i+1)

		for j := 0; j < cfg.NumClients; j++ {
			provider, err := NewProvider(cfg, logger)
			if err != nil {
				return nil, fmt.Errorf("Failed to create client %d for key %d: %w", j+1, i+1, err)
			}
			providerPool = append(providerPool, provider)

		}
	}
	return providerPool, nil
}
