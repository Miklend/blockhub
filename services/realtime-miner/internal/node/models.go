package node

import (
	"context"
	"lib/models"
)

type Worker interface {
	TransferBlocks(ctx context.Context, in <-chan *models.BlockDTO) error
}

type RtCollector interface {
	SubscribeNewBlocks(ctx context.Context, maxRetries int) (<-chan *models.BlockDTO, error)
}
