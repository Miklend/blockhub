package models

import "context"

type Worker interface {
	TransferBlocks(ctx context.Context, in <-chan *Block)
}