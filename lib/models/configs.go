package models

import "time"

type Config struct {
	Provider   `json:"provider"`
	Broker     `json:"broker"`
	Clickhouse `json:"clickhouse"`
	Redis      `json:"redis"`
}

type Provider struct {
	ProviderType string `json:"provider_type"`
	NetworkName  string `json:"network_name"`
	BaseURL      string `json:"base_url"`
	ApiKey       string `json:"api_key"`
}

type Broker struct {
	BrockerType  string
	Brokers      []string
	GroupID      string
	StartOffset  int64
	BatchSize    int
	BatchTimeout time.Duration
	Async        bool
}

type Clickhouse struct{}

type Redis struct{}
