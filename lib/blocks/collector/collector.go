package collector

import (
	"lib/clients/node"
	"lib/utils/logging"
	"sync/atomic"

	"golang.org/x/time/rate"
)

type BlockCollector struct {
	clients      []node.Provider
	clientsCount uint64
	limiter      *rate.Limiter
	logger       *logging.Logger
}

func NewBlockCollector(clients []node.Provider, limiterRate float64, logger *logging.Logger) *BlockCollector {

	limiter := rate.NewLimiter(rate.Limit(limiterRate), 1)

	blk := &BlockCollector{
		clients: clients,
		logger:  logger,
		limiter: limiter,
	}
	return blk
}

func (bc *BlockCollector) Client() node.Provider {
	return bc.clients[0]
}
func (bc *BlockCollector) Logger() *logging.Logger {
	return bc.logger
}

func (bc *BlockCollector) GetNextClient() node.Provider {
	current := atomic.AddUint64(&bc.clientsCount, 1)
	return bc.clients[(current-1)%uint64(len(bc.clients))]
}
