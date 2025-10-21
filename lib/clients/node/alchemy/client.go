package alchemy

import (
	"context"
	"fmt"
	"lib/clients/node"
	"lib/models"
	"lib/utils/logging"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type alchemyClient struct {
	networkName string
	apiKey      string
	baseURL     string
	ethClient   *ethclient.Client
	rpcClient   *rpc.Client
	logger      *logging.Logger
}

// NewAlchemyClient создаёт новый клиент и возвращает интерфейс node.Provider
func NewAlchemyClient(cfg models.Provider, logger *logging.Logger) (node.Provider, error) {
	logger.Infof("Initializing Alchemy client for network: %s", cfg.NetworkName)

	fullURL := fmt.Sprintf("%s%s", cfg.BaseURL, cfg.ApiKey)
	logger.Debugf("Connecting to Alchemy endpoint: %s...", fullURL)

	rpcClient, err := rpc.Dial(fullURL)
	if err != nil {
		logger.Errorf("Failed to connect to Alchemy network %s: %v", cfg.NetworkName, err)
		return nil, fmt.Errorf("failed connect to %s: %w", cfg.NetworkName, err)
	}

	ethClient := ethclient.NewClient(rpcClient)

	logger.Infof("Successfully connected to Alchemy network: %s", cfg.NetworkName)

	return &alchemyClient{
		networkName: cfg.NetworkName,
		apiKey:      cfg.ApiKey,
		baseURL:     cfg.BaseURL,
		ethClient:   ethClient,
		rpcClient:   rpcClient,
		logger:      logger,
	}, nil
}

func (a *alchemyClient) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	return a.ethClient.BlockByNumber(ctx, number)
}

func (a *alchemyClient) BlockReceipts(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) ([]*types.Receipt, error) {
	return a.ethClient.BlockReceipts(ctx, blockNrOrHash)
}

func (a *alchemyClient) SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error) {
	return a.ethClient.SubscribeNewHead(ctx, ch)
}

func (a *alchemyClient) BatchCallContext(ctx context.Context, batch []rpc.BatchElem) error {
	return a.rpcClient.BatchCallContext(ctx, batch)
}

func (a *alchemyClient) Close() {
	a.logger.Debugf("Closing Alchemy client connection for network: %s", a.networkName)
	a.ethClient.Close()
	a.rpcClient.Close()
}

func (a *alchemyClient) Name() string {
	return "alchemy"
}
