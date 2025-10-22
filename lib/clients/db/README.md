# Централизованная система клиентов СУБД

Эта система предоставляет единый интерфейс для работы с различными типами баз данных (PostgreSQL, ClickHouse, Redis) в одном месте.

## Особенности

- **Единый интерфейс**: Все клиенты реализуют общие интерфейсы для удобства использования
- **Централизованное управление**: Все клиенты создаются и управляются в одном месте
- **Типобезопасность**: Использование типов для различения SQL и NoSQL клиентов
- **Простота смены СУБД**: Для смены базы данных достаточно изменить конфигурацию

## Структура

```
lib/clients/db/
├── client.go          # Универсальные интерфейсы
├── factory.go         # Фабрика клиентов и централизованное хранилище
├── postgresql/
│   └── client.go      # PostgreSQL клиент
├── clickhouse/
│   └── client.go      # ClickHouse клиент
└── redis/
    └── client.go      # Redis клиент
```

## Конфигурация

```go
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
```

## Использование

### Создание централизованного хранилища

```go
storage, err := db.NewStorage(ctx, cfg, logger)
if err != nil {
    log.Fatal(err)
}
defer storage.Close()
```

### Получение SQL клиентов

```go
// PostgreSQL
postgresClient, err := storage.GetSQLClient(db.PostgreSQL)

// ClickHouse
clickhouseClient, err := storage.GetSQLClient(db.ClickHouse)
```

### Получение NoSQL клиентов

```go
// Redis
redisClient, err := storage.GetNoSQLClient(db.Redis)
```

### Работа с SQL клиентами

```go
// Выполнение запроса
rows, err := sqlClient.Query(ctx, "SELECT * FROM users WHERE id = $1", userID)
if err != nil {
    return err
}
defer rows.Close()

// Итерация по результатам
for rows.Next() {
    var user User
    if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
        return err
    }
    // обработка пользователя
}

// Выполнение команды
err := sqlClient.Exec(ctx, "INSERT INTO users (name, email) VALUES ($1, $2)", name, email)

// Транзакции
tx, err := sqlClient.Begin(ctx)
if err != nil {
    return err
}
defer tx.Rollback(ctx)

err = tx.Exec(ctx, "INSERT INTO users (name) VALUES ($1)", name)
if err != nil {
    return err
}

err = tx.Commit(ctx)
```

### Работа с NoSQL клиентами (Redis)

```go
// Установка значения
err := redisClient.Set(ctx, "key", "value")

// Получение значения
value, err := redisClient.Get(ctx, "key")

// Проверка существования
exists, err := redisClient.Exists(ctx, "key")

// Работа с хешами
err = redisClient.HSet(ctx, "user:1", "name", "John")
name, err := redisClient.HGet(ctx, "user:1", "name")
```

## Интерфейсы

### Client (базовый интерфейс)
- `Ping(ctx context.Context) error` - проверка соединения
- `Close() error` - закрытие соединения
- `GetType() DatabaseType` - получение типа базы данных

### SQLClient (для SQL баз данных)
- Наследует `Client`
- `Exec(ctx context.Context, query string, args ...any) error`
- `Query(ctx context.Context, query string, args ...any) (Rows, error)`
- `QueryRow(ctx context.Context, query string, args ...any) Row`
- `Begin(ctx context.Context) (Transaction, error)`
- `PrepareBatch(ctx context.Context, query string) (Batch, error)`

### NoSQLClient (для NoSQL баз данных)
- Наследует `Client`
- `Set(ctx context.Context, key string, value interface{}) error`
- `Get(ctx context.Context, key string) (string, error)`
- `Del(ctx context.Context, key string) error`
- `Exists(ctx context.Context, key string) (bool, error)`
- `HSet(ctx context.Context, key, field string, value interface{}) error`
- `HGet(ctx context.Context, key, field string) (string, error)`

## Преимущества

1. **Единообразие**: Все клиенты используют одинаковые интерфейсы
2. **Централизация**: Управление всеми клиентами в одном месте
3. **Гибкость**: Легко добавлять новые типы баз данных
4. **Типобезопасность**: Компилятор проверяет корректность использования
5. **Простота тестирования**: Легко создавать моки для тестов
6. **Простота смены СУБД**: Достаточно изменить конфигурацию

## Добавление новой СУБД

1. Добавить новый тип в `DatabaseType`
2. Добавить конфигурацию в `models.Config`
3. Создать клиент, реализующий соответствующий интерфейс
4. Добавить создание клиента в фабрику
5. Обновить зависимости в `go.mod`
