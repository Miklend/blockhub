package fabricClient

import (
	"fmt"
	"lib/clients/node"
	"lib/clients/node/alchemy"
	"lib/models"
	"lib/utils/logging"
)

const (
	Alchemy = "alchemy"
)

func NewProvaider(cfg models.Provider, logger *logging.Logger) (node.Provider, error) {
	switch cfg.ProvaiderType {
	case Alchemy:
		return alchemy.NewAlchemyClient(cfg, logger)
	default:
		return nil, fmt.Errorf("not found provaider client")
	}
}
