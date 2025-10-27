-- Примеры запросов для работы с данными Ethereum

-- 1. Получить последние 10 блоков
SELECT 
    number,
    hash,
    timestamp,
    gas_used,
    gas_limit,
    miner
FROM blocks 
ORDER BY number DESC 
LIMIT 10;

-- 2. Получить транзакции определенного адреса за последний день
SELECT 
    hash,
    `from`,
    `to`,
    value,
    gas_price,
    block_timestamp
FROM transactions 
WHERE (`from` = '0x...' OR `to` = '0x...')
    AND block_timestamp >= now() - INTERVAL 1 DAY
ORDER BY block_timestamp DESC;

-- 3. Получить статистику по блокам за последний час
SELECT 
    count() as block_count,
    avg(gas_used) as avg_gas_used,
    avg(gas_limit) as avg_gas_limit,
    avg(size) as avg_size
FROM blocks 
WHERE timestamp >= now() - INTERVAL 1 HOUR;

-- 4. Получить топ-10 майнеров по количеству блоков
SELECT 
    miner,
    count() as block_count
FROM blocks 
GROUP BY miner 
ORDER BY block_count DESC 
LIMIT 10;

-- 5. Получить логи событий по определенному топику
SELECT 
    address,
    topics,
    data,
    block_timestamp
FROM logs 
WHERE topic0 = '0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef' -- Transfer event
ORDER BY block_timestamp DESC 
LIMIT 100;

-- 6. Получить статистику по транзакциям по часам
SELECT 
    toStartOfHour(block_timestamp) as hour,
    count() as tx_count,
    sum(gas_used) as total_gas_used,
    avg(gas_price) as avg_gas_price
FROM transactions 
WHERE block_timestamp >= now() - INTERVAL 24 HOUR
GROUP BY hour 
ORDER BY hour;

-- 7. Получить информацию о контрактах (созданных транзакциях)
SELECT 
    t.hash as tx_hash,
    r.contract_address,
    t.block_timestamp
FROM transactions t
JOIN receipts r ON t.hash = r.transaction_hash
WHERE r.contract_address IS NOT NULL
ORDER BY t.block_timestamp DESC 
LIMIT 100;

-- 8. Получить статистику по gas использованию
SELECT 
    quantile(0.5)(gas_used) as median_gas_used,
    quantile(0.9)(gas_used) as p90_gas_used,
    quantile(0.99)(gas_used) as p99_gas_used,
    max(gas_used) as max_gas_used
FROM transactions 
WHERE block_timestamp >= now() - INTERVAL 1 DAY;
