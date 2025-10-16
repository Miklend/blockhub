package node

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

type Provider interface {
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
	BlockReceipts(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) ([]*types.Receipt, error)
	SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error)
	BatchCallContext(ctx context.Context, batch []rpc.BatchElem) error
	Close()
	Name() string
}
