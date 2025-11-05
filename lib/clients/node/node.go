package node

import (
	"context"
	"lib/models"

	"github.com/ethereum/go-ethereum/rpc"
)

type Provider interface {
	BlockByNumber(ctx context.Context, param string) (models.BlockDTO, error)
	TxByHash(ctx context.Context, param string) (models.TxDTO, error)
	ReceiptByTxHash(ctx context.Context, param string) (models.ReceiptDTO, error)
	ReceiptByBlockNumber(ctx context.Context, param string) ([]models.ReceiptDTO, error)

	BatchRequest(ctx context.Context, requests []rpc.BatchElem) ([]rpc.BatchElem, error)
	BatchBlockByNumber(ctx context.Context, numbers []string) (map[models.Hash]models.BlockDTO, error)
	BatchReceiptByBlockNumber(ctx context.Context, numbers []string) (map[models.Hash][]models.ReceiptDTO, error)
	BatchReceiptByTxHash(ctx context.Context, hashes []string) (map[models.Hash]models.ReceiptDTO, error)
	BatchBlockWithReceiptByNumber(ctx context.Context, numbers []string) (map[models.Hash]models.BlockDTO, error)

	SubscribeBlockWithReceipts(ctx context.Context, blockCh chan<- *models.BlockDTO) (*rpc.ClientSubscription, error)
	SubscribePendingTransactions(ctx context.Context, txsCh chan<- *models.TxDTO) (*rpc.ClientSubscription, error)

	Close() error
}
