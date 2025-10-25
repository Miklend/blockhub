package models

import "context"

type MessageBroker struct {
	Key     []byte
	Value   []byte
	Headers map[string]string
	Topic   string
}

type MessageHandlerBroker func(ctx context.Context, msg MessageBroker) error
