package main

import (
	"context"
	fabricClient "lib/clients/fabric_client"
	"lib/models"
	"lib/utils/logging"
	"log"
)

func main() {
	cfg := models.Provider{
		ProviderType: "alchemy",
		NetworkName:  "ETH_mainnet",
		BaseURL:      "https://eth-mainnet.g.alchemy.com/v2/",
		ApiKey:       "bYiMENtDz_cHRTwZIkBiV",
	}

	logger := logging.GetLogger()

	provider, err := fabricClient.NewProvider(cfg, logger)
	if err != nil {
		log.Fatalf("Fatal error connecting to provider: %v", err)
	}
	defer provider.Close()
	logger.Infof("Successfully initialized provider: %s", provider.Name())

	block, err := provider.BlockByNumber(context.Background(), nil)
	if err != nil {
		logger.Errorf("Failed to get latest block: %v", err)
		return
	}

	logger.Infof("Latest block retrieved: %d", block.Number())
}
