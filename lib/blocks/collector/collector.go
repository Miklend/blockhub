package collector

import (
	"lib/clients/node"
	"lib/utils/logging"

	"golang.org/x/time/rate"
)

type BlockCollector struct {
	client       node.Provider
	clientsCount uint64
	limiter      *rate.Limiter
	logger       *logging.Logger
}

func NewBlockCollector(client node.Provider, limiterRate float64, logger *logging.Logger) *BlockCollector {

	limiter := rate.NewLimiter(rate.Limit(limiterRate), 1)

	blk := &BlockCollector{
		client:  client,
		logger:  logger,
		limiter: limiter,
	}
	return blk
}

func (bc *BlockCollector) Client() node.Provider {
	return bc.client
}
func (bc *BlockCollector) Logger() *logging.Logger {
	return bc.logger
}

// func (bc *BlockCollector) GetNextClient() node.Provider {
// 	current := atomic.AddUint64(&bc.clientsCount, 1)
// 	return bc.clients[(current-1)%uint64(len(bc.clients))]
// }
