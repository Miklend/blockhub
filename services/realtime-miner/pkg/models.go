package pkg

import (
	"context"
	"lib/clients/node"
	"lib/models"
	"lib/utils/logging"
)

type Worker interface {
	TransferBlocks(ctx context.Context, in <-chan *models.Block) error
}

type RtCollector interface {
	SubscribeNewBlocks(ctx context.Context, maxRetries int) (<-chan *models.Block, error)
}
type BlockCollectorInterface interface {
	CollectBlockByNumber(ctx context.Context, blockNumber uint64) (*models.Block, error)
	Client() node.Provider
	Logger() *logging.Logger
}
