package metrics

import (
	"lib/models"
	"strconv"
	"strings"
)

func trim0x(s string) string {
	return strings.TrimPrefix(s, "0x")
}
func hexToUint64(hexStr string) uint64 {
	if hexStr == "" {
		return 0
	}
	val, err := strconv.ParseUint(trim0x(hexStr), 16, 64)
	if err != nil {
		return 0
	}
	return val
}
func toHashArray(list []string) []models.Hash {
	out := make([]models.Hash, len(list))
	for i, v := range list {
		out[i] = models.Hash(v)
	}
	return out
}

func ConvertBlockDTO(b models.BlockRPCDTO) models.BlockDTO {
	var txs []models.TxDTO
	for _, l := range b.Transactions {
		txs = append(txs, ConvertTxDTO(l))
	}

	return models.BlockDTO{
		Hash:             models.Hash(b.Hash),
		Number:           hexToUint64(b.Number),
		ParentHash:       models.Hash(b.ParentHash),
		Nonce:            b.Nonce,
		Sha3Uncles:       models.Hash(b.Sha3Uncles),
		LogsBloom:        b.LogsBloom,
		TransactionsRoot: models.Hash(b.TransactionsRoot),
		StateRoot:        models.Hash(b.StateRoot),
		ReceiptsRoot:     models.Hash(b.ReceiptsRoot),
		Miner:            models.Address(b.Miner),
		Difficulty:       b.Difficulty,
		TotalDifficulty:  b.TotalDifficulty,
		Size:             hexToUint64(b.Size),
		ExtraData:        b.ExtraData,
		GasLimit:         hexToUint64(b.GasLimit),
		GasUsed:          hexToUint64(b.GasUsed),
		BaseFeePerGas:    b.BaseFeePerGas,
		Timestamp:        models.Time(hexToUint64(b.Timestamp)),
		MixHash:          models.Hash(b.MixHash),
		Transactions:     txs,
		Uncles:           toHashArray(b.Uncles),
	}
}
func ConvertTxDTO(t models.TxRpcDTO) models.TxDTO {
	var to *models.Address
	if t.To != nil {
		addr := models.Address(*t.To)
		to = &addr
	}

	var accessList []models.AccessListEntryDTO
	for _, a := range t.AccessList {
		accessList = append(accessList, models.AccessListEntryDTO{
			Address:     models.Address(a.Address),
			StorageKeys: a.StorageKeys,
		})
	}

	return models.TxDTO{
		Hash:                 models.Hash(t.Hash),
		BlockHash:            models.Hash(t.BlockHash),
		BlockNumber:          hexToUint64(t.BlockNumber),
		TransactionIndex:     hexToUint64(t.TransactionIndex),
		From:                 models.Address(t.From),
		To:                   to,
		Value:                t.Value,
		Gas:                  hexToUint64(t.Gas),
		GasPrice:             t.GasPrice,
		Input:                t.Input,
		Nonce:                hexToUint64(t.Nonce),
		Type:                 t.Type,
		MaxFeePerGas:         t.MaxFeePerGas,
		MaxPriorityFeePerGas: t.MaxPriorityFeePerGas,
		ChainID:              t.ChainID,
		V:                    t.V,
		R:                    t.R,
		S:                    t.S,
		AccessList:           accessList,
	}
}
func ConvertReceiptDTO(r models.ReceiptRpcDTO) models.ReceiptDTO {
	var to *models.Address
	if r.To != nil {
		addr := models.Address(*r.To)
		to = &addr
	}
	var contract *models.Address
	if r.ContractAddress != nil {
		addr := models.Address(*r.ContractAddress)
		contract = &addr
	}

	logs := make([]models.LogDTO, len(r.Logs))
	for i, l := range r.Logs {
		logs[i] = ConvertLogDTO(l)
	}

	return models.ReceiptDTO{
		TransactionHash:   models.Hash(r.TransactionHash),
		TransactionIndex:  hexToUint64(r.TransactionIndex),
		BlockHash:         models.Hash(r.BlockHash),
		BlockNumber:       hexToUint64(r.BlockNumber),
		From:              models.Address(r.From),
		To:                to,
		ContractAddress:   contract,
		CumulativeGasUsed: hexToUint64(r.CumulativeGasUsed),
		GasUsed:           hexToUint64(r.GasUsed),
		EffectiveGasPrice: r.EffectiveGasPrice,
		Status:            r.Status,
		LogsBloom:         r.LogsBloom,
		Logs:              logs,
	}
}
func ConvertLogDTO(l models.LogRpcDTO) models.LogDTO {
	return models.LogDTO{
		BlockNumber:      hexToUint64(l.BlockNumber),
		BlockHash:        models.Hash(l.BlockHash),
		TransactionHash:  models.Hash(l.TransactionHash),
		TransactionIndex: hexToUint64(l.TransactionIndex),
		LogIndex:         hexToUint64(l.LogIndex),
		Address:          models.Address(l.Address),
		Data:             l.Data,
		Topics:           l.Topics,
		Removed:          l.Removed,
	}
}
