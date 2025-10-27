# ClickHouse Service

Сервис для работы с ClickHouse базой данных, реализующий CRUD операции для блоков, транзакций, квитанций и логов Ethereum.

## Структура проекта

```
services/clickhouse-service/
├── cmd/
│   └── main.go                    # Точка входа приложения
├── internal/
│   └── db/
│       ├── db.go                  # Интерфейс для работы с БД
│       └── click_house/
│           ├── client.go          # Основной клиент ClickHouse
│           ├── converters.go      # Конвертеры для преобразования данных
│           ├── block/             # Работа с блоками
│           │   ├── insert.go      # Вставка блоков
│           │   └── fetch.go       # Получение блоков
│           ├── tx/                # Работа с транзакциями
│           │   ├── insert.go      # Вставка транзакций
│           │   └── fetch.go       # Получение транзакций
│           ├── receipt/           # Работа с квитанциями
│           │   ├── insert.go      # Вставка квитанций
│           │   └── fetch.go       # Получение квитанций
│           └── log/               # Работа с логами
│               ├── insert.go      # Вставка логов
│               └── fetch.go       # Получение логов
├── db-schema/                     # Схема базы данных
│   ├── README.md                  # Документация схемы
│   ├── schema.sql                 # Основной файл схемы
│   ├── examples.sql               # Примеры запросов
│   └── tables/                    # Отдельные файлы таблиц с индексами
│       ├── blocks.sql             # Таблица блоков + индексы
│       ├── transactions.sql       # Таблица транзакций + индексы
│       ├── receipts.sql           # Таблица квитанций + индексы
│       └── logs.sql               # Таблица логов + индексы
├── Dockerfile                     # Docker образ для продакшена
├── Dockerfile.dev                 # Docker образ для разработки
├── go.mod                         # Go модули
└── README.md                      # Документация
```

## Возможности

### Блоки
- `InsertBlock` - вставка одного блока
- `InsertBlocks` - вставка массива блоков
- `FetchBlock` - получение блока по хешу
- `FetchBlocks` - получение блоков по хешам
- `FetchBlockByNumber` - получение блока по номеру
- `FetchBlocksByRange` - получение блоков в диапазоне номеров

### Транзакции
- `InsertTx` - вставка одной транзакции
- `InsertTxs` - вставка массива транзакций
- `InsertTxWithBlockData` - вставка транзакции с данными блока
- `InsertTxsWithBlockData` - вставка массива транзакций с данными блока
- `FetchTx` - получение транзакции по хешу
- `FetchTxs` - получение транзакций по хешам
- `FetchTxsByBlock` - получение транзакций по хешу блока
- `FetchTxsByBlockNumber` - получение транзакций по номеру блока
- `FetchTxsByAddress` - получение транзакций по адресу

### Квитанции
- `InsertReceipt` - вставка одной квитанции
- `InsertReceipts` - вставка массива квитанций
- `InsertReceiptsFromBlock` - вставка квитанций из блока
- `FetchReceipt` - получение квитанции по хешу транзакции
- `FetchReceipts` - получение квитанций по хешам транзакций
- `FetchReceiptsByBlock` - получение квитанций по хешу блока
- `FetchReceiptsByBlockNumber` - получение квитанций по номеру блока
- `FetchReceiptsByAddress` - получение квитанций по адресу

### Логи
- `InsertLog` - вставка одного лога
- `InsertLogs` - вставка массива логов
- `InsertLogsFromReceipt` - вставка логов из квитанции
- `InsertLogsFromBlock` - вставка логов из блока
- `FetchLogsByTransaction` - получение логов по хешу транзакции
- `FetchLogsByBlock` - получение логов по хешу блока
- `FetchLogsByBlockNumber` - получение логов по номеру блока
- `FetchLogsByAddress` - получение логов по адресу
- `FetchLogsByTopic` - получение логов по топику
- `FetchLogsByTopic0` - получение логов по первому топику
- `FetchLogsByAddressAndTopic` - получение логов по адресу и топику

## Использование

### Инициализация

```go
import (
    "context"
    clickhouseClient "lib/clients/db/clickhouse"
    "lib/models"
    "lib/utils/logging"
    "clickhouse-service/internal/db/click_house"
)

// Получаем конфигурацию
config := models.GetConfig(logger)

// Создаем ClickHouse клиент
client, err := clickhouseClient.NewClient(ctx, config.Clickhouse)
if err != nil {
    log.Fatal(err)
}

// Создаем репозиторий
repo := clickhouseRepo.NewClickhouseService(client, logger)
```

### Примеры использования

#### Вставка блока
```go
block := models.Block{
    Hash: "0x123...",
    Number: 12345,
    // ... другие поля
}

err := repo.InsertBlock("blocks", block)
if err != nil {
    log.Fatal(err)
}
```

#### Получение блока по номеру
```go
block, err := repo.FetchBlockByNumber("blocks", 12345)
if err != nil {
    log.Fatal(err)
}
```

#### Вставка транзакций с данными блока
```go
txs := []models.Tx{
    // ... транзакции
}

err := repo.InsertTxsWithBlockData("transactions", txs, blockHash, blockNumber, blockTimestamp)
if err != nil {
    log.Fatal(err)
}
```

## Схема базы данных

Схема базы данных находится в папке `db-schema/` и включает:

- **tables/** - отдельные файлы для каждой таблицы с индексами
- **schema.sql** - основной файл для создания всех таблиц
- **examples.sql** - примеры запросов

### Создание таблиц

```bash
# Создать все таблицы
clickhouse-client < db-schema/schema.sql

# Создать отдельную таблицу с индексами
clickhouse-client < db-schema/tables/blocks.sql
```

### Активация индексов

Индексы закомментированы в файлах таблиц. Для их активации:

```bash
# Раскомментировать нужные индексы в файле таблицы
# Затем выполнить:
clickhouse-client < db-schema/tables/blocks.sql
```

## Конфигурация

Сервис использует конфигурацию из `lib/models.Config`:

```yaml
clickhouse:
  host: "localhost"
  port: 9000
  username: "default"
  password: ""
  database: "ethereum"
  ssl_mode: "disable"
```

## Зависимости

- `lib/clients/db` - интерфейсы и клиенты для работы с БД
- `lib/models` - модели данных
- `lib/utils/logging` - утилиты для логирования
- `github.com/ClickHouse/clickhouse-go/v2` - ClickHouse драйвер

## Запуск

```bash
go run cmd/main.go
```

## Docker

### Сборка образа
```bash
docker build -t clickhouse-service .
```

### Запуск контейнера
```bash
docker run -d --name clickhouse-service clickhouse-service
```
