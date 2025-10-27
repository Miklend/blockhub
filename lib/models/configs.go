package models

import (
	"flag"
	"lib/utils/logging"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	ProviderRealTime Provider   `yaml:"provider_realtime"`
	Broker           Broker     `yaml:"broker"`
	Clickhouse       Clickhouse `yaml:"clickhouse"`
	Redis            Redis      `yaml:"redis"`
}

type Provider struct {
	ProviderType string `yaml:"provider_type"`
	NetworkName  string `yaml:"network_name"`
	BaseURL      string `yaml:"base_url"`
	ApiKey       string `yaml:"api_key"`
	NumClients   int    `yaml:"num_clients"`

	Limiter    int `yaml:"limiter"`
	MaxRetries int `yaml:"max_retries"`
}

type Broker struct {
	BrockerType  string        `yaml:"brocker_type"`
	Brokers      []string      `yaml:"brokers"`
	GroupID      string        `yaml:"group_id"`
	StartOffset  int64         `yaml:"start_offset"`
	BatchSize    int           `yaml:"batch_size"`
	BatchTimeout time.Duration `yaml:"batch_timeout"`
	Async        bool          `yaml:"async"`
}

type DB struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
	SSLMode  string
}

type Clickhouse struct {
	DB
}

type Redis struct {
	DB
}

// Константы, используемые для поиска конфига
const (
	flagConfigPathName = "configs"
	envConfigPathName  = "CONFIG_PATH"
	dotEnvFileName     = ".env"
)

var (
	instance *Config
	once     sync.Once
)

func GetConfig(logger *logging.Logger) *Config {
	logger.Debug("start get config")

	once.Do(func() {
		// Загружаем .env, но не падаем, если файла нет
		_ = godotenv.Load(dotEnvFileName)

		var configPath string
		// 1. Чтение пути к конфигу из флагов командной строки
		flag.StringVar(&configPath, flagConfigPathName, "", "path to config file (e.g., ./configs/config.yaml)")
		flag.Parse()

		// 2. Перезаписываем путь переменной окружения, если есть
		if path, ok := os.LookupEnv(envConfigPathName); ok && path != "" {
			configPath = path
		}

		// 3. Если путь не указан нигде, устанавливаем путь по умолчанию
		if configPath == "" {
			// Пробуем разные возможные пути
			possiblePaths := []string{
				"./configs/configs.yaml",     // для realtime-miner
				"./configs/config.yaml",      // общий путь
				"../configs/configs.yaml",    // если запускаем из cmd/
				"../../configs/configs.yaml", // если запускаем из глубоких папок
			}

			for _, path := range possiblePaths {
				if _, err := os.Stat(path); err == nil {
					configPath = path
					logger.Debugf("Found config at: %s", path)
					break
				}
			}

			// Если ни один путь не найден, используем первый по умолчанию
			if configPath == "" {
				configPath = "./configs/configs.yaml"
			}
		}

		instance = &Config{}

		// 4. Чтение и парсинг YAML-файла
		if readErr := cleanenv.ReadConfig(configPath, instance); readErr != nil {
			description, descrErr := cleanenv.GetDescription(instance, nil)
			if descrErr != nil {
				panic(descrErr)
			}

			slog.Error(
				"Failed to read config. Ensure 'config.yaml' exists or all required env variables are set.",
				slog.String("error", readErr.Error()),
				slog.String("config_path", configPath),
				slog.String("description", description),
			)
			os.Exit(1)
		}

		logger.Debug("config loaded successfully", slog.Any("config", instance))
	})

	return instance
}
