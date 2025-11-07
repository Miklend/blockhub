package clickhouseRepo

import (
	"strconv"
	"time"

	"lib/models"
)

func convertBlockToClickHouseRow(block models.Block) []interface{} {
	var baseFeePerGas *uint64
	if block.BaseFeePerGas != nil {
		val := uint64(*block.BaseFeePerGas)
		baseFeePerGas = &val
	}

	difficulty, _ := parseHexToUint256(block.Difficulty)
	totalDifficulty, _ := parseHexToUint256(block.TotalDifficulty)

	return []interface{}{
		block.Hash,
		uint64(block.Number),
		block.ParentHash,
		formatUint64ToHex(uint64(block.Nonce)),
		block.Sha3Uncles,
		block.LogsBloom,
		block.TransactionsRoot,
		block.StateRoot,
		block.ReceiptsRoot,
		block.Miner,
		difficulty,
		totalDifficulty,
		uint64(block.Size),
		block.ExtraData,
		uint64(block.GasLimit),
		uint64(block.GasUsed),
		baseFeePerGas,
		block.Timestamp,
		block.MixHash,
		copyStringSlice(block.Transactions),
		copyStringSlice(block.Uncles),
		block.Timestamp,
	}
}

func convertTxToClickHouseRow(tx models.Tx, blockHash string, blockNumber uint64, blockTimestamp time.Time) []interface{} {
	var maxFeePerGas *uint64
	if tx.MaxFeePerGas != nil {
		val := uint64(*tx.MaxFeePerGas)
		maxFeePerGas = &val
	}

	var maxPriorityFeePerGas *uint64
	if tx.MaxPriorityFeePerGas != nil {
		val := uint64(*tx.MaxPriorityFeePerGas)
		maxPriorityFeePerGas = &val
	}

	value, _ := parseHexToUint256(tx.Value)

	return []interface{}{
		tx.Hash,
		blockHash,
		blockNumber,
		uint32(tx.TransactionIndex),
		tx.From,
		tx.To,
		value,
		uint64(tx.Gas),
		uint64(tx.GasPrice),
		tx.Input,
		uint64(tx.Nonce),
		uint8(tx.Type),
		maxFeePerGas,
		maxPriorityFeePerGas,
		uint64(tx.ChainID),
		tx.V,
		tx.R,
		tx.S,
		tx.AccessList,
		blockTimestamp,
	}
}

func convertReceiptToClickHouseRow(receipt models.Receipt, txHash string, txIndex uint32, blockHash string, blockNumber uint64, blockTimestamp time.Time) []interface{} {
	return []interface{}{
		txHash,
		txIndex,
		blockHash,
		blockNumber,
		receipt.From,
		receipt.To,
		receipt.ContractAddress,
		uint64(receipt.CumulativeGasUsed),
		uint64(receipt.GasUsed),
		uint64(receipt.EffectiveGasPrice),
		uint8(receipt.Status),
		receipt.LogsBloom,
		blockTimestamp,
	}
}

func convertLogToClickHouseRow(log models.Log, blockHash string, blockNumber uint64, blockTimestamp time.Time, txIndex uint32) []interface{} {
	timestamp := log.BlockTimestamp
	if timestamp.IsZero() {
		timestamp = blockTimestamp
	}

	return []interface{}{
		blockNumber,
		blockHash,
		log.TransactionHash,
		txIndex,
		uint32(log.LogIndex),
		log.Address,
		log.Data,
		copyStringSlice(log.Topics),
		timestamp,
	}
}

func parseHexToUint256(hexStr string) (string, error) {
	if hexStr == "" {
		return "", nil
	}
	if len(hexStr) > 2 && hexStr[:2] == "0x" {
		hexStr = hexStr[2:]
	}
	return hexStr, nil
}

func formatUint64ToHex(n uint64) string {
	return "0x" + strconv.FormatUint(n, 16)
}

func copyStringSlice(src []string) []string {
	if len(src) == 0 {
		return nil
	}
	dst := make([]string, len(src))
	copy(dst, src)
	return dst
}
