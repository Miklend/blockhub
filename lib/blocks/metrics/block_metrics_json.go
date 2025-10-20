package metrics

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"lib/models"
	"math/big"
	"strings"
)

// NewBlockFromJSON парсит сырые JSON блока и квитанций в удобную модель models.Block
func NewBlockFromJSON(blockRaw json.RawMessage, receiptsRaw json.RawMessage) models.Block {
	var rb map[string]interface{}
	if err := json.Unmarshal(blockRaw, &rb); err != nil {
		fmt.Printf("failed to unmarshal block JSON: %v\n", err)
		return models.Block{}
	}

	var rcs []map[string]interface{}
	if err := json.Unmarshal(receiptsRaw, &rcs); err != nil {
		fmt.Printf("failed to unmarshal receipt JSON: %v\n", err)
		rcs = []map[string]interface{}{}
	}

	blk := models.Block{
		BaseFeePerGas:    parseString(rb["baseFeePerGas"]),
		Difficulty:       parseString(rb["difficulty"]),
		ExtraData:        decodeToHex(rb["extraData"]),
		GasLimit:         parseUint64(rb["gasLimit"]),
		GasUsed:          parseUint64(rb["gasUsed"]),
		Hash:             parseString(rb["hash"]),
		LogsBloom:        decodeToHex(rb["logsBloom"]),
		Miner:            parseString(rb["miner"]),
		MixHash:          parseString(rb["mixHash"]),
		Nonce:            parseUint64(rb["nonce"]),
		Number:           parseUint64(rb["number"]),
		ParentHash:       parseString(rb["parentHash"]),
		ReceiptsRoot:     parseString(rb["receiptsRoot"]),
		Sha3Uncles:       parseString(rb["sha3Uncles"]),
		Size:             parseUint64(rb["size"]),
		StateRoot:        parseString(rb["stateRoot"]),
		Timestamp:        parseUint64(rb["timestamp"]),
		TransactionsRoot: parseString(rb["transactionsRoot"]),
	}

	// Uncles
	if uncles, ok := rb["uncles"].([]interface{}); ok {
		for _, u := range uncles {
			blk.Uncles = append(blk.Uncles, parseString(u))
		}
	}

	// Транзакции
	if txs, ok := rb["transactions"].([]interface{}); ok {
		for i, txRaw := range txs {
			txMap := txRaw.(map[string]interface{})
			tx := models.Tx{
				From:                 parseString(txMap["from"]),
				Gas:                  parseUint64(txMap["gas"]),
				GasPrice:             parseString(txMap["gasPrice"]),
				MaxFeePerGas:         parseString(txMap["maxFeePerGas"]),
				MaxPriorityFeePerGas: parseString(txMap["maxPriorityFeePerGas"]),
				Hash:                 parseString(txMap["hash"]),
				Input:                decodeToHex(txMap["input"]),
				Nonce:                parseUint64(txMap["nonce"]),
				To:                   parseString(txMap["to"]),
				TransactionIndex:     parseUint64(txMap["transactionIndex"]),
				Value:                parseString(txMap["value"]),
				Type:                 uint8(parseUint64(txMap["type"])),
				ChainID:              parseString(txMap["chainId"]),
				V:                    parseString(txMap["v"]),
				R:                    parseString(txMap["r"]),
				S:                    parseString(txMap["s"]),
			}

			// Квитанция
			if i < len(rcs) {
				tx.Receipt = convertReceiptJSON(rcs[i])
			}

			blk.Transactions = append(blk.Transactions, tx)
		}
	}

	return blk
}

// вспомогательные функции
func parseString(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%v", v)
}

func parseUint64(v interface{}) uint64 {
	if v == nil {
		return 0
	}
	switch val := v.(type) {
	case string:
		n, _ := new(big.Int).SetString(strings.Replace(val, "0x", "", 1), 16)
		if n != nil {
			return n.Uint64()
		}
	case float64:
		return uint64(val)
	}
	return 0
}

func decodeToHex(v interface{}) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return hex.EncodeToString([]byte(val))
	default:
		return fmt.Sprintf("%v", val)
	}
}

func convertReceiptJSON(r map[string]interface{}) *models.Receipt {
	rec := &models.Receipt{
		ContractAddress:   parseString(r["contractAddress"]),
		CumulativeGasUsed: parseUint64(r["cumulativeGasUsed"]),
		EffectiveGasPrice: parseString(r["effectiveGasPrice"]),
		From:              parseString(r["from"]),
		GasUsed:           parseUint64(r["gasUsed"]),
		LogsBloom:         decodeToHex(r["logsBloom"]),
		Status:            parseUint64(r["status"]),
		Type:              uint8(parseUint64(r["type"])),
	}

	if to, ok := r["to"]; ok {
		rec.To = parseString(to)
	}

	// Логи
	if logs, ok := r["logs"].([]interface{}); ok {
		for _, l := range logs {
			lMap := l.(map[string]interface{})
			rec.Logs = append(rec.Logs, models.Log{
				Address:         parseString(lMap["address"]),
				Topics:          parseStringSlice(lMap["topics"]),
				Data:            decodeToHex(lMap["data"]),
				TransactionHash: parseString(lMap["transactionHash"]),
				LogIndex:        parseUint64(lMap["logIndex"]),
				Removed:         lMap["removed"].(bool),
			})
		}
	}

	return rec
}

func parseStringSlice(v interface{}) []string {
	if v == nil {
		return nil
	}
	if arr, ok := v.([]interface{}); ok {
		var res []string
		for _, a := range arr {
			res = append(res, parseString(a))
		}
		return res
	}
	return nil
}
