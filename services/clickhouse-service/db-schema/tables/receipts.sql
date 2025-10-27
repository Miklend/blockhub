-- Таблица квитанций транзакций Ethereum
CREATE TABLE receipts
(
    `transaction_hash` FixedString(66),
    `transaction_index` UInt32,
    `block_hash` FixedString(66),
    `block_number` UInt64,
    `from` FixedString(42),
    `to` Nullable(FixedString(42)),
    `contract_address` Nullable(FixedString(42)),
    `cumulative_gas_used` UInt64,
    `gas_used` UInt64,
    `effective_gas_price` UInt64,
    `status` UInt8, -- 1 (успех) или 0 (реверт)
    `logs_bloom` String,
    `block_timestamp` DateTime64(3, 'UTC'),
    `date` Date MATERIALIZED toDate(block_timestamp)
)
ENGINE = ReplacingMergeTree
PARTITION BY toYYYYMM(block_timestamp)
ORDER BY (block_number, transaction_index);

-- Индексы для таблицы receipts
-- CREATE INDEX idx_receipts_tx_hash ON receipts (transaction_hash) TYPE bloom_filter GRANULARITY 1;
-- CREATE INDEX idx_receipts_from ON receipts (`from`) TYPE bloom_filter GRANULARITY 1;
-- CREATE INDEX idx_receipts_to ON receipts (`to`) TYPE bloom_filter GRANULARITY 1;
-- CREATE INDEX idx_receipts_contract ON receipts (contract_address) TYPE bloom_filter GRANULARITY 1;
-- CREATE INDEX idx_receipts_block_hash ON receipts (block_hash) TYPE bloom_filter GRANULARITY 1;
-- CREATE INDEX idx_receipts_block_number ON receipts (block_number) TYPE minmax GRANULARITY 1;
-- CREATE INDEX idx_receipts_status ON receipts (status) TYPE set(0) GRANULARITY 1;
