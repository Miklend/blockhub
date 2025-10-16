package models

type Config struct {
	Provider   `json:"provider"`
	Broker     `json:"broker"`
	Clickhouse `json:"clickhouse"`
	Redis      `json:"redis"`
}

type Provider struct {
	NetworkName string `json:"network_name"`
	BaseURL     string `json:"base_url"`
	ApiKey      string `json:"api_key"`
}

type Broker struct{}

type Clickhouse struct{}

type Redis struct{}
