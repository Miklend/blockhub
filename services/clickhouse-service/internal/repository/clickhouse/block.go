package clickhouseRepo

import (
	"context"
	"fmt"
	"lib/models"
)

func (c *ClickhouseRepo) InsertBlock(table string, block models.Block) error {
	c.logger.Debugf("InsertBlock: table=%s hash=%s", table, block.Hash)

	ctx := context.Background()

	// Формируем SQL для вставки одного блока
	query := fmt.Sprintf(`
		INSERT INTO %s (
			baseFeePerGas, difficulty, extraData, gasLimit, gasUsed, hash,
			logsBloom, miner, mixHash, nonce, number, parentHash,
			receiptsRoot, sha3Uncles, size, stateRoot, timestamp,
			transactionsRoot
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, table)

	err := c.client.Exec(ctx, query,
		block.BaseFeePerGas,
		block.Difficulty,
		block.ExtraData,
		block.GasLimit,
		block.GasUsed,
		block.Hash,
		block.LogsBloom,
		block.Miner,
		block.MixHash,
		block.Nonce,
		block.Number,
		block.ParentHash,
		block.ReceiptsRoot,
		block.Sha3Uncles,
		block.Size,
		block.StateRoot,
		block.Timestamp,
		block.TransactionsRoot,
	)

	if err != nil {
		c.logger.Errorf("InsertBlock error: %v", err)
		return err
	}

	c.logger.Debugf("InsertBlock success: %s", block.Hash)
	return nil
}

func (c *ClickhouseRepo) InsertBlocks(table string, blocks []models.Block) error {
	c.logger.Debugf("InsertBlocks: table=%s count=%d", table, len(blocks))
	// TODO: реализовать batch вставку
	return nil
}

func (c *ClickhouseRepo) FetchBlock(table string, hashBlock string) (models.Block, error) {
	c.logger.Debugf("FetchBlock: table=%s hash=%s", table, hashBlock)
	// TODO: реализовать SELECT по hash
	return models.Block{}, nil
}

func (c *ClickhouseRepo) FetchBlocks(table string, hashBlocks []string) ([]models.Block, error) {
	c.logger.Debugf("FetchBlocks: table=%s count=%d", table, len(hashBlocks))
	// TODO: реализовать SELECT по множеству hash
	return nil, nil
}
