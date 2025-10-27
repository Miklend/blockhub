package clickhouseRepo

import (
	"strconv"
	"time"

	"lib/models"
)

// Вспомогательные функции для парсинга

// parseHexToUint64 парсит hex строку в uint64
func parseHexToUint64(hexStr string) (uint64, error) {
	if hexStr == "" {
		return 0, nil
	}
	// Убираем префикс 0x если есть
	if len(hexStr) > 2 && hexStr[:2] == "0x" {
		hexStr = hexStr[2:]
	}
	return strconv.ParseUint(hexStr, 16, 64)
}

// parseHexToUint256 парсит hex строку в строку для UInt256
func parseHexToUint256(hexStr string) (string, error) {
	// Для UInt256 в ClickHouse используем строку
	// Убираем префикс 0x если есть
	if len(hexStr) > 2 && hexStr[:2] == "0x" {
		hexStr = hexStr[2:]
	}
	return hexStr, nil
}

// formatUint64ToHex конвертирует uint64 в hex строку
func formatUint64ToHex(n uint64) string {
	return "0x" + strconv.FormatUint(n, 16)
}

// parseHexToUint64Safe безопасно парсит hex строку в uint64
func parseHexToUint64Safe(hexStr string) uint64 {
	if hexStr == "" {
		return 0
	}
	// Убираем префикс 0x если есть
	if len(hexStr) > 2 && hexStr[:2] == "0x" {
		hexStr = hexStr[2:]
	}
	val, err := strconv.ParseUint(hexStr, 16, 64)
	if err != nil {
		return 0
	}
	return val
}

// convertBlockToClickHouseRow конвертирует Block в строку для вставки в ClickHouse
func convertBlockToClickHouseRow(block models.Block) []interface{} {
	// Конвертируем timestamp из Unix в time.Time
	timestamp := time.Unix(int64(block.Timestamp), 0)

	// Извлекаем хеши транзакций
	txHashes := make([]string, len(block.Transactions))
	for i, tx := range block.Transactions {
		txHashes[i] = tx.Hash
	}

	// Конвертируем baseFeePerGas если есть
	var baseFeePerGas *uint64
	if block.BaseFeePerGas != "" {
		if val, err := parseHexToUint64(block.BaseFeePerGas); err == nil {
			baseFeePerGas = &val
		}
	}

	// Конвертируем difficulty и totalDifficulty в UInt256 (строки)
	difficulty, _ := parseHexToUint256(block.Difficulty)
	totalDifficulty, _ := parseHexToUint256(block.Difficulty) // В реальном проекте нужно вычислять

	return []interface{}{
		block.Hash,                     // hash FixedString(66)
		block.Number,                   // number UInt64
		block.ParentHash,               // parent_hash FixedString(66)
		formatUint64ToHex(block.Nonce), // nonce String
		block.Sha3Uncles,               // sha3_uncles FixedString(66)
		block.LogsBloom,                // logs_bloom String
		block.TransactionsRoot,         // transactions_root FixedString(66)
		block.StateRoot,                // state_root FixedString(66)
		block.ReceiptsRoot,             // receipts_root FixedString(66)
		block.Miner,                    // miner FixedString(42)
		difficulty,                     // difficulty UInt256
		totalDifficulty,                // total_difficulty UInt256
		block.Size,                     // size UInt64
		block.ExtraData,                // extra_data String
		block.GasLimit,                 // gas_limit UInt64
		block.GasUsed,                  // gas_used UInt64
		baseFeePerGas,                  // base_fee_per_gas Nullable(UInt64)
		timestamp,                      // timestamp DateTime64(3, 'UTC')
		block.MixHash,                  // mix_hash FixedString(66)
		txHashes,                       // transactions Array(FixedString(66))
		block.Uncles,                   // uncles Array(FixedString(66))
		// date автоматически вычисляется из timestamp
	}
}

// convertTxToClickHouseRow конвертирует Tx в строку для вставки в ClickHouse
func convertTxToClickHouseRow(tx models.Tx, blockHash string, blockNumber uint64, blockTimestamp uint64) []interface{} {
	timestamp := time.Unix(int64(blockTimestamp), 0)

	// Конвертируем to в указатель
	var to *string
	if tx.To != "" {
		to = &tx.To
	}

	// Конвертируем maxFeePerGas если есть
	var maxFeePerGas *uint64
	if tx.MaxFeePerGas != "" {
		if val, err := parseHexToUint64(tx.MaxFeePerGas); err == nil {
			maxFeePerGas = &val
		}
	}

	// Конвертируем maxPriorityFeePerGas если есть
	var maxPriorityFeePerGas *uint64
	if tx.MaxPriorityFeePerGas != "" {
		if val, err := parseHexToUint64(tx.MaxPriorityFeePerGas); err == nil {
			maxPriorityFeePerGas = &val
		}
	}

	// Конвертируем chainID
	chainID, _ := parseHexToUint64(tx.ChainID)

	// Конвертируем gasPrice
	gasPrice, _ := parseHexToUint64(tx.GasPrice)

	// Конвертируем value в UInt256 (строка)
	value, _ := parseHexToUint256(tx.Value)

	return []interface{}{
		tx.Hash,                     // hash FixedString(66)
		blockHash,                   // block_hash FixedString(66)
		blockNumber,                 // block_number UInt64
		uint32(tx.TransactionIndex), // transaction_index UInt32
		tx.From,                     // from FixedString(42)
		to,                          // to Nullable(FixedString(42))
		value,                       // value UInt256
		tx.Gas,                      // gas UInt64
		gasPrice,                    // gas_price UInt64
		tx.Input,                    // input String
		tx.Nonce,                    // nonce UInt64
		tx.Type,                     // type UInt8
		maxFeePerGas,                // max_fee_per_gas Nullable(UInt64)
		maxPriorityFeePerGas,        // max_priority_fee_per_gas Nullable(UInt64)
		chainID,                     // chain_id UInt64
		tx.V,                        // v String
		tx.R,                        // r FixedString(66)
		tx.S,                        // s FixedString(66)
		"",                          // access_list String
		timestamp,                   // block_timestamp DateTime64(3, 'UTC')
		// date автоматически вычисляется из block_timestamp
	}
}

// convertReceiptToClickHouseRow конвертирует Receipt в строку для вставки в ClickHouse
func convertReceiptToClickHouseRow(receipt models.Receipt, txHash string, txIndex uint32, blockHash string, blockNumber uint64, blockTimestamp uint64) []interface{} {
	timestamp := time.Unix(int64(blockTimestamp), 0)

	// Конвертируем to в указатель
	var to *string
	if receipt.To != "" {
		to = &receipt.To
	}

	// Конвертируем contractAddress в указатель
	var contractAddress *string
	if receipt.ContractAddress != "" {
		contractAddress = &receipt.ContractAddress
	}

	// Конвертируем effectiveGasPrice
	effectiveGasPrice, _ := parseHexToUint64(receipt.EffectiveGasPrice)

	return []interface{}{
		txHash,                    // transaction_hash FixedString(66)
		txIndex,                   // transaction_index UInt32
		blockHash,                 // block_hash FixedString(66)
		blockNumber,               // block_number UInt64
		receipt.From,              // from FixedString(42)
		to,                        // to Nullable(FixedString(42))
		contractAddress,           // contract_address Nullable(FixedString(42))
		receipt.CumulativeGasUsed, // cumulative_gas_used UInt64
		receipt.GasUsed,           // gas_used UInt64
		effectiveGasPrice,         // effective_gas_price UInt64
		uint8(receipt.Status),     // status UInt8
		receipt.LogsBloom,         // logs_bloom String
		timestamp,                 // block_timestamp DateTime64(3, 'UTC')
		// date автоматически вычисляется из block_timestamp
	}
}

// convertLogToClickHouseRow конвертирует Log в строку для вставки в ClickHouse
func convertLogToClickHouseRow(log models.Log, blockHash string, blockNumber uint64, blockTimestamp uint64, txIndex uint32) []interface{} {
	timestamp := time.Unix(int64(blockTimestamp), 0)

	return []interface{}{
		blockNumber,          // block_number UInt64
		blockHash,            // block_hash FixedString(66)
		log.TransactionHash,  // transaction_hash FixedString(66)
		txIndex,              // transaction_index UInt32
		uint32(log.LogIndex), // log_index UInt32
		log.Address,          // address FixedString(42)
		log.Data,             // data String
		log.Topics,           // topics Array(FixedString(66))
		timestamp,            // block_timestamp DateTime64(3, 'UTC')
		// date и topic0 автоматически вычисляются
	}
}
