package models

import "time"

type Config struct {
	Provider   `json:"provider"`
	Broker     `json:"broker"`
	Clickhouse `json:"clickhouse"`
	Redis      `json:"redis"`
	PostgreSQL `json:"postgresql"`
}

type Provider struct {
	ProvaiderType string
	NetworkName   string `json:"network_name"`
	BaseURL       string `json:"base_url"`
	ApiKey        string `json:"api_key"`
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

type Clickhouse struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Redis struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

type PostgreSQL struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
	SSLMode  string `json:"ssl_mode"`
}
