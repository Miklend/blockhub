package main

import (
	"context"
	"fmt"
	"log"

	"lib/clients/db"
	"lib/models"
	"lib/utils/logging"
)

func main() {
	// Создаем логгер
	logger := logging.NewLogger()

	// Создаем конфигурацию
	cfg := models.Config{
		PostgreSQL: models.PostgreSQL{
			Host:     "localhost",
			Port:     "5432",
			Database: "testdb",
			Username: "user",
			Password: "password",
			SSLMode:  "disable",
		},
		Clickhouse: models.Clickhouse{
			Host:     "localhost",
			Port:     "9000",
			Database: "testdb",
			Username: "default",
			Password: "",
		},
		Redis: models.Redis{
			Host:     "localhost",
			Port:     "6379",
			Password: "",
			DB:       0,
		},
	}

	// Создаем централизованное хранилище клиентов
	ctx := context.Background()
	storage, err := db.NewStorage(ctx, cfg, logger)
	if err != nil {
		log.Fatalf("Failed to create storage: %v", err)
	}
	defer storage.Close()

	// Проверяем соединение со всеми базами данных
	if err := storage.Ping(ctx); err != nil {
		log.Fatalf("Ping failed: %v", err)
	}

	fmt.Println("All database connections are healthy!")

	// Пример работы с PostgreSQL
	postgresClient, err := storage.GetSQLClient(db.PostgreSQL)
	if err != nil {
		log.Printf("Failed to get PostgreSQL client: %v", err)
	} else {
		fmt.Printf("PostgreSQL client type: %s\n", postgresClient.GetType())

		// Пример выполнения запроса
		rows, err := postgresClient.Query(ctx, "SELECT version()")
		if err != nil {
			log.Printf("PostgreSQL query failed: %v", err)
		} else {
			defer rows.Close()
			if rows.Next() {
				var version string
				if err := rows.Scan(&version); err == nil {
					fmt.Printf("PostgreSQL version: %s\n", version)
				}
			}
		}
	}

	// Пример работы с ClickHouse
	clickhouseClient, err := storage.GetSQLClient(db.ClickHouse)
	if err != nil {
		log.Printf("Failed to get ClickHouse client: %v", err)
	} else {
		fmt.Printf("ClickHouse client type: %s\n", clickhouseClient.GetType())

		// Пример выполнения запроса
		rows, err := clickhouseClient.Query(ctx, "SELECT version()")
		if err != nil {
			log.Printf("ClickHouse query failed: %v", err)
		} else {
			defer rows.Close()
			if rows.Next() {
				var version string
				if err := rows.Scan(&version); err == nil {
					fmt.Printf("ClickHouse version: %s\n", version)
				}
			}
		}
	}

	// Пример работы с Redis
	redisClient, err := storage.GetNoSQLClient(db.Redis)
	if err != nil {
		log.Printf("Failed to get Redis client: %v", err)
	} else {
		fmt.Printf("Redis client type: %s\n", redisClient.GetType())

		// Пример работы с Redis
		if err := redisClient.Set(ctx, "test_key", "test_value"); err != nil {
			log.Printf("Redis set failed: %v", err)
		} else {
			if value, err := redisClient.Get(ctx, "test_key"); err == nil {
				fmt.Printf("Redis value: %s\n", value)
			}
		}
	}

	fmt.Println("Database operations completed successfully!")
}
