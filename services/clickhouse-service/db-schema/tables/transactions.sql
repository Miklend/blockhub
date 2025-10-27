-- Таблица транзакций Ethereum
CREATE TABLE transactions
(
    `hash` FixedString(66),
    `block_hash` FixedString(66),
    `block_number` UInt64,
    `transaction_index` UInt32,
    `from` FixedString(42),
    `to` Nullable(FixedString(42)), 
    `value` UInt256,
    `gas` UInt64,
    `gas_price` UInt64,
    `input` String,
    `nonce` UInt64,
    `type` UInt8,
    `max_fee_per_gas` Nullable(UInt64),
    `max_priority_fee_per_gas` Nullable(UInt64),
    `chain_id` UInt64,
    `v` String,
    `r` FixedString(66),
    `s` FixedString(66),
    `access_list` String,
    `block_timestamp` DateTime64(3, 'UTC'),
    `date` Date MATERIALIZED toDate(block_timestamp)
)
ENGINE = ReplacingMergeTree
PARTITION BY toYYYYMM(block_timestamp)
ORDER BY (block_number, transaction_index);

-- Индексы для таблицы transactions
-- CREATE INDEX idx_transactions_hash ON transactions (hash) TYPE bloom_filter GRANULARITY 1;
-- CREATE INDEX idx_transactions_from ON transactions (`from`) TYPE bloom_filter GRANULARITY 1;
-- CREATE INDEX idx_transactions_to ON transactions (`to`) TYPE bloom_filter GRANULARITY 1;
-- CREATE INDEX idx_transactions_block_hash ON transactions (block_hash) TYPE bloom_filter GRANULARITY 1;
-- CREATE INDEX idx_transactions_block_number ON transactions (block_number) TYPE minmax GRANULARITY 1;
-- CREATE INDEX idx_transactions_timestamp ON transactions (block_timestamp) TYPE minmax GRANULARITY 1;
