-- Полная схема базы данных ClickHouse для Ethereum данных
-- Этот файл создает все таблицы для хранения блоков, транзакций, квитанций и логов

-- Создание базы данных (если не существует)
-- CREATE DATABASE IF NOT EXISTS ethereum;

-- Использование базы данных
-- USE ethereum;

-- Таблица блоков
-- Содержит информацию о блоках Ethereum
SOURCE tables/blocks.sql;

-- Таблица транзакций  
-- Содержит информацию о транзакциях Ethereum
SOURCE tables/transactions.sql;

-- Таблица квитанций
-- Содержит информацию о квитанциях транзакций
SOURCE tables/receipts.sql;

-- Таблица логов
-- Содержит информацию о логах событий
SOURCE tables/logs.sql;
