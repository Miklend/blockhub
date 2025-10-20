package collector

import (
	"lib/clients/node"
	"lib/utils/logging"
)

type BlockCollector struct {
	client node.Provider
	logger *logging.Logger
}

func NewBlockCollector(client node.Provider, logger *logging.Logger) *BlockCollector {
	blk := &BlockCollector{
		client: client,
		logger: logger,
	}
	return blk
}

func (bc *BlockCollector) Client() node.Provider {
	return bc.client
}
func (bc *BlockCollector) Logger() *logging.Logger {
	return bc.logger
}
