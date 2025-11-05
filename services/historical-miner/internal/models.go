package internal

import (
	"context"
	"lib/models"
)

type Worker interface {
	ProcessReceipts(ctx context.Context, in <-chan *models.BlockDTO) error
}
