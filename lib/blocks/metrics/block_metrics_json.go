package metrics

import (
	"encoding/json"
	"fmt"
	"lib/models"
	"math/big"
	"strconv"
	"strings"
	"time"
)

// NewBlockFromJSON парсит сырые JSON блока и квитанций в models.Block
func NewBlockFromJSON(blockRaw json.RawMessage, receiptsRaw json.RawMessage) models.Block {
	var rb map[string]interface{}
	if err := json.Unmarshal(blockRaw, &rb); err != nil {
		fmt.Printf("failed to unmarshal block JSON: %v\n", err)
		return models.Block{}
	}

	var rcs []map[string]interface{}
	if err := json.Unmarshal(receiptsRaw, &rcs); err != nil {
		fmt.Printf("failed to unmarshal receipts JSON: %v\n", err)
		rcs = []map[string]interface{}{}
	}

	blk := models.Block{
		Hash:             parseString(rb["hash"]),
		Number:           parseUint(rb["number"]),
		ParentHash:       parseString(rb["parentHash"]),
		Nonce:            parseUint(rb["nonce"]),
		Sha3Uncles:       parseString(rb["sha3Uncles"]),
		LogsBloom:        parseString(rb["logsBloom"]),
		TransactionsRoot: parseString(rb["transactionsRoot"]),
		StateRoot:        parseString(rb["stateRoot"]),
		ReceiptsRoot:     parseString(rb["receiptsRoot"]),
		Miner:            parseString(rb["miner"]),
		Difficulty:       parseString(rb["difficulty"]),
		TotalDifficulty:  parseString(rb["totalDifficulty"]),
		Size:             parseUint(rb["size"]),
		ExtraData:        parseString(rb["extraData"]),
		GasLimit:         parseUint(rb["gasLimit"]),
		GasUsed:          parseUint(rb["gasUsed"]),
		BaseFeePerGas:    parseOptionalUint(rb["baseFeePerGas"]),
		Timestamp:        parseTime(rb["timestamp"]),
		MixHash:          parseString(rb["mixHash"]),
	}

	// Uncles
	if uncles, ok := rb["uncles"].([]interface{}); ok {
		for _, u := range uncles {
			blk.Uncles = append(blk.Uncles, parseString(u))
		}
	}

	// Transactions
	if txs, ok := rb["transactions"].([]interface{}); ok {
		for i, txRaw := range txs {
			txMap, ok := txRaw.(map[string]interface{})
			if !ok {
				continue
			}

			tx := models.Tx{
				Hash:                 parseString(txMap["hash"]),
				BlockHash:            parseString(txMap["blockHash"]),
				BlockNumber:          parseUint(txMap["blockNumber"]),
				TransactionIndex:     parseUint(txMap["transactionIndex"]),
				From:                 parseString(txMap["from"]),
				To:                   parseOptionalString(txMap["to"]),
				Value:                parseString(txMap["value"]),
				Gas:                  parseUint(txMap["gas"]),
				GasPrice:             parseUint(txMap["gasPrice"]),
				Input:                parseString(txMap["input"]),
				Nonce:                parseUint(txMap["nonce"]),
				Type:                 parseUint(txMap["type"]),
				MaxFeePerGas:         parseOptionalUint(txMap["maxFeePerGas"]),
				MaxPriorityFeePerGas: parseOptionalUint(txMap["maxPriorityFeePerGas"]),
				ChainID:              parseUint(txMap["chainId"]),
				V:                    parseString(txMap["v"]),
				R:                    parseString(txMap["r"]),
				S:                    parseString(txMap["s"]),
				AccessList:           parseString(txMap["accessList"]),
				BlockTimestamp:       blk.Timestamp,
			}

			// Привязываем квитанцию
			if i < len(rcs) {
				tx.Receipt = convertReceiptJSON(rcs[i], blk.Timestamp)
			}

			blk.Transactions = append(blk.Transactions, tx.Hash)
		}
	}

	return blk
}

// convertReceiptJSON преобразует JSON-квитанцию в models.Receipt
func convertReceiptJSON(r map[string]interface{}, ts time.Time) *models.Receipt {
	rec := &models.Receipt{
		TransactionHash:   parseString(r["transactionHash"]),
		TransactionIndex:  parseUint(r["transactionIndex"]),
		BlockHash:         parseString(r["blockHash"]),
		BlockNumber:       parseUint(r["blockNumber"]),
		From:              parseString(r["from"]),
		To:                parseOptionalString(r["to"]),
		ContractAddress:   parseOptionalString(r["contractAddress"]),
		CumulativeGasUsed: parseUint(r["cumulativeGasUsed"]),
		GasUsed:           parseUint(r["gasUsed"]),
		EffectiveGasPrice: parseUint(r["effectiveGasPrice"]),
		Status:            parseUint(r["status"]),
		LogsBloom:         parseString(r["logsBloom"]),
		BlockTimestamp:    ts,
	}

	// Логи
	if logs, ok := r["logs"].([]interface{}); ok {
		for _, l := range logs {
			lMap, ok := l.(map[string]interface{})
			if !ok {
				continue
			}
			rec.Logs = append(rec.Logs, models.Log{
				BlockNumber:      parseUint(lMap["blockNumber"]),
				BlockHash:        parseString(lMap["blockHash"]),
				TransactionHash:  parseString(lMap["transactionHash"]),
				TransactionIndex: parseUint(lMap["transactionIndex"]),
				LogIndex:         parseUint(lMap["logIndex"]),
				Address:          parseString(lMap["address"]),
				Data:             parseString(lMap["data"]),
				Topics:           parseStringSlice(lMap["topics"]),
				BlockTimestamp:   ts,
				Topic0:           firstTopic(lMap["topics"]),
			})
		}
	}

	return rec
}

func parseString(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%v", v)
}

func parseOptionalString(v interface{}) *string {
	if v == nil {
		return nil
	}
	s := parseString(v)
	if s == "" {
		return nil
	}
	return &s
}

func parseUint(v interface{}) uint {
	if v == nil {
		return 0
	}
	switch val := v.(type) {
	case float64:
		return uint(val)
	case json.Number:
		n, _ := val.Int64()
		return uint(n)
	case string:
		// поддерживаем 0x-формат
		s := strings.TrimPrefix(val, "0x")
		if s == "" {
			return 0
		}
		n := new(big.Int)
		n.SetString(s, 16)
		return uint(n.Uint64())
	default:
		return 0
	}
}

func parseOptionalUint(v interface{}) *uint {
	u := parseUint(v)
	if u == 0 {
		return nil
	}
	return &u
}

func parseTime(v interface{}) time.Time {
	switch val := v.(type) {
	case float64:
		return time.Unix(int64(val), 0).UTC()
	case string:
		// timestamp может быть в hex, decimal или RFC3339
		if strings.HasPrefix(val, "0x") {
			n := new(big.Int)
			n.SetString(strings.TrimPrefix(val, "0x"), 16)
			return time.Unix(n.Int64(), 0).UTC()
		}
		if t, err := strconv.ParseInt(val, 10, 64); err == nil {
			return time.Unix(t, 0).UTC()
		}
		if t, err := time.Parse(time.RFC3339, val); err == nil {
			return t.UTC()
		}
	}
	return time.Time{}
}

func parseStringSlice(v interface{}) []string {
	if v == nil {
		return nil
	}
	arr, ok := v.([]interface{})
	if !ok {
		return nil
	}
	out := make([]string, 0, len(arr))
	for _, a := range arr {
		out = append(out, parseString(a))
	}
	return out
}

func firstTopic(v interface{}) string {
	topics := parseStringSlice(v)
	if len(topics) > 0 {
		return topics[0]
	}
	return ""
}
