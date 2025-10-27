-- Таблица блоков Ethereum
CREATE TABLE blocks
(
    `hash` FixedString(66),
    `number` UInt64,
    `parent_hash` FixedString(66),
    `nonce` String,
    `sha3_uncles` FixedString(66),
    `logs_bloom` String,
    `transactions_root` FixedString(66),
    `state_root` FixedString(66),
    `receipts_root` FixedString(66),
    `miner` FixedString(42),
    `difficulty` UInt256,
    `total_difficulty` UInt256,
    `size` UInt64,
    `extra_data` String,
    `gas_limit` UInt64,
    `gas_used` UInt64,
    `base_fee_per_gas` Nullable(UInt64), 
    `timestamp` DateTime64(3, 'UTC'),
    `mix_hash` FixedString(66),
    `transactions` Array(FixedString(66)),
    `uncles` Array(FixedString(66)),
    `date` Date MATERIALIZED toDate(timestamp)
)
ENGINE = ReplacingMergeTree
PARTITION BY toYYYYMM(timestamp)
ORDER BY (number, timestamp);

-- Индексы для таблицы blocks
-- CREATE INDEX idx_blocks_hash ON blocks (hash) TYPE bloom_filter GRANULARITY 1;
-- CREATE INDEX idx_blocks_miner ON blocks (miner) TYPE bloom_filter GRANULARITY 1;
-- CREATE INDEX idx_blocks_timestamp ON blocks (timestamp) TYPE minmax GRANULARITY 1;
-- CREATE INDEX idx_blocks_parent_hash ON blocks (parent_hash) TYPE bloom_filter GRANULARITY 1;
