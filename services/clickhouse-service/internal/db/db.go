package db

import "lib/models"

type DB interface {
	Close() error

	InsertBlock(table string, block models.Block) error
	InsertBlocks(table string, block []models.Block) error
	FetchBlock(table string, hashBlock string) (models.Block, error)
	FetchBlocks(table string, hashBlocks []string) ([]models.Block, error)

	InsertTx(table string, block models.Tx) error
	InsertTxs(table string, block []models.Tx) error
}
