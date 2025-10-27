package collector

import (
	"lib/clients/node"
	"lib/utils/logging"
	"sync/atomic"
)

type BlockCollector struct {
	clients      []node.Provider
	clientsCount uint64
	logger       *logging.Logger
}

func NewBlockCollector(clients []node.Provider, logger *logging.Logger) *BlockCollector {
	blk := &BlockCollector{
		clients: clients,
		logger:  logger,
	}
	return blk
}

func (bc *BlockCollector) Client() node.Provider {
	return bc.clients[0]
}
func (bc *BlockCollector) Logger() *logging.Logger {
	return bc.logger
}

func NewHistorycalBlockCollector(clients []node.Provider, logger *logging.Logger) *BlockCollector {
	blk := &BlockCollector{
		clients: clients,
		logger:  logger,
	}
	return blk
}

func (bc *BlockCollector) GetNextClient() node.Provider {
	current := atomic.AddUint64(&bc.clientsCount, 1)
	return bc.clients[(current-1)%uint64(len(bc.clients))]
}
