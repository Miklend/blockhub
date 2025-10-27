package worker

import (
	"context"

	"lib/utils/logging"
)

// структура, распределяющая сырые номера блоков по каналу заданий
type Dispatcher struct {
	logger         *logging.Logger
	incomingBlocks <-chan uint64
	jobsChan       chan<- uint64
}

// Конструктор диспатчера
func NewDispatcher(
	logger *logging.Logger,
	incomingBlocks <-chan uint64,
	jobsChan chan<- uint64,
) *Dispatcher {
	return &Dispatcher{
		logger:         logger,
		incomingBlocks: incomingBlocks,
		jobsChan:       jobsChan,
	}
}

// функция мгновенно передающая приходящие номера блоков в очередь
func (d *Dispatcher) Dispatch(ctx context.Context) {
	d.logger.Info("Dispatcher started, putting block numbers to job queue")

	for {
		//Ожидание нового номера блока либо остановки
		select {
		case blockNumber, ok := <-d.incomingBlocks:
			if !ok {
				d.logger.Warn("Incoming blocks channel closed. Stopping dispatcher")
				return
			}
			// Пробует записать номер в канал очереди задач
			select {
			case d.jobsChan <- blockNumber:
				d.logger.Debugf("Dispatcher block number #%d to job queue.", blockNumber)
			case <-ctx.Done():
				d.logger.Warn("Context done while dispatching, stopping.")
				return
			}
		case <-ctx.Done():
			d.logger.Info("Dispatcher recieved shutdown signal, stopping")
			return
		}
	}
}
