package node

import (
	"context"
	"lib/models"
)

type Worker interface {
	TransferBlocks(ctx context.Context, in <-chan *models.Block) error
}

type RtCollector interface {
	SubscribeNewBlocks(ctx context.Context, maxRetries int) (<-chan *models.Block, error)
}
