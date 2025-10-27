-- Таблица логов событий Ethereum
CREATE TABLE logs
(
    `block_number` UInt64,
    `block_hash` FixedString(66),
    `transaction_hash` FixedString(66),
    `transaction_index` UInt32,
    `log_index` UInt32,
    `address` FixedString(42),
    `data` String,
    `topics` Array(FixedString(66)), 
    `block_timestamp` DateTime64(3, 'UTC'),
    `date` Date MATERIALIZED toDate(block_timestamp),
    `topic0` FixedString(66) MATERIALIZED if(length(topics) > 0, topics[1], '')
)
ENGINE = ReplacingMergeTree
PARTITION BY toYYYYMM(block_timestamp)
ORDER BY (address, topic0, block_number, transaction_index, log_index);

-- Индексы для таблицы logs
-- CREATE INDEX idx_logs_address ON logs (address) TYPE bloom_filter GRANULARITY 1;
-- CREATE INDEX idx_logs_tx_hash ON logs (transaction_hash) TYPE bloom_filter GRANULARITY 1;
-- CREATE INDEX idx_logs_topic0 ON logs (topic0) TYPE bloom_filter GRANULARITY 1;
-- CREATE INDEX idx_logs_block_hash ON logs (block_hash) TYPE bloom_filter GRANULARITY 1;
-- CREATE INDEX idx_logs_block_number ON logs (block_number) TYPE minmax GRANULARITY 1;
-- CREATE INDEX idx_logs_timestamp ON logs (block_timestamp) TYPE minmax GRANULARITY 1;
