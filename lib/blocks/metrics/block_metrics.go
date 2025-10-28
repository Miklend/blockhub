package metrics

import (
	"lib/models"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// NewBlockMetrics преобразует блок и квитанции из go-ethereum в models.Block (JSON-safe)
func NewBlockMetrics(block *types.Block, receipts []*types.Receipt) models.Block {
	if block == nil {
		return models.Block{}
	}

	uncles := make([]string, len(block.Uncles()))
	for i, u := range block.Uncles() {
		uncles[i] = u.Hash().Hex()
	}

	var baseFee *uint
	if block.BaseFee() != nil {
		val := uint(block.BaseFee().Uint64())
		baseFee = &val
	}

	timestamp := time.Unix(int64(block.Time()), 0).UTC()

	blk := models.Block{
		Hash:             block.Hash().Hex(),
		Number:           uint(block.NumberU64()),
		ParentHash:       block.ParentHash().Hex(),
		Nonce:            uint(block.Nonce()),
		Sha3Uncles:       block.UncleHash().Hex(),
		LogsBloom:        "0x" + common.Bytes2Hex(block.Bloom().Bytes()),
		TransactionsRoot: block.TxHash().Hex(),
		StateRoot:        block.Root().Hex(),
		ReceiptsRoot:     block.ReceiptHash().Hex(),
		Miner:            block.Coinbase().Hex(),
		Difficulty:       block.Difficulty().String(),
		TotalDifficulty:  "", // если ты её не считаешь — можно позже добавить
		Size:             uint(block.Size()),
		ExtraData:        "0x" + common.Bytes2Hex(block.Extra()),
		GasLimit:         uint(block.GasLimit()),
		GasUsed:          uint(block.GasUsed()),
		BaseFeePerGas:    baseFee,
		Timestamp:        timestamp,
		MixHash:          block.MixDigest().Hex(),
		Transactions:     make([]string, 0, len(block.Transactions())),
		Uncles:           uncles,
	}

	// Транзакции
	for i, tx := range block.Transactions() {
		var receipt *types.Receipt
		if i < len(receipts) {
			receipt = receipts[i]
		}

		txModel := NewTx(tx, receipt, block, timestamp, i)
		blk.Transactions = append(blk.Transactions, txModel.Hash)
	}

	return blk
}

// NewTx преобразует транзакцию из go-ethereum в models.Tx
func NewTx(tx *types.Transaction, receipt *types.Receipt, block *types.Block, ts time.Time, index int) models.Tx {
	v, r, s := tx.RawSignatureValues()
	from := getTxSender(tx)

	var (
		maxFeePerGas, maxPriorityFeePerGas *uint
	)

	if tx.GasFeeCap() != nil {
		val := uint(tx.GasFeeCap().Uint64())
		maxFeePerGas = &val
	}
	if tx.GasTipCap() != nil {
		val := uint(tx.GasTipCap().Uint64())
		maxPriorityFeePerGas = &val
	}

	txModel := models.Tx{
		Hash:                 tx.Hash().Hex(),
		BlockHash:            block.Hash().Hex(),
		BlockNumber:          uint(block.NumberU64()),
		TransactionIndex:     uint(index),
		From:                 from.Hex(),
		To:                   addressToOptionalString(tx.To()),
		Value:                tx.Value().String(), // big.Int в string
		Gas:                  uint(tx.Gas()),
		GasPrice:             uint(tx.GasPrice().Uint64()),
		Input:                "0x" + common.Bytes2Hex(tx.Data()),
		Nonce:                uint(tx.Nonce()),
		Type:                 uint(tx.Type()),
		MaxFeePerGas:         maxFeePerGas,
		MaxPriorityFeePerGas: maxPriorityFeePerGas,
		ChainID:              uint(tx.ChainId().Uint64()),
		V:                    v.String(),
		R:                    r.String(),
		S:                    s.String(),
		AccessList:           "", // если нужно сериализовать — потом добавим
		BlockTimestamp:       ts,
	}

	if receipt != nil {
		txModel.Receipt = NewReceipt(receipt, tx, block, ts, index)
	}

	return txModel
}

// NewReceipt преобразует go-ethereum Receipt в models.Receipt
func NewReceipt(receipt *types.Receipt, tx *types.Transaction, block *types.Block, ts time.Time, index int) *models.Receipt {
	if receipt == nil {
		return nil
	}

	rec := &models.Receipt{
		TransactionHash:   tx.Hash().Hex(),
		TransactionIndex:  uint(index),
		BlockHash:         block.Hash().Hex(),
		BlockNumber:       uint(block.NumberU64()),
		From:              getTxSender(tx).Hex(),
		To:                addressToOptionalString(tx.To()),
		ContractAddress:   addressToOptionalString(&receipt.ContractAddress),
		CumulativeGasUsed: uint(receipt.CumulativeGasUsed),
		GasUsed:           uint(receipt.GasUsed),
		EffectiveGasPrice: uint(receipt.EffectiveGasPrice.Uint64()),
		Status:            uint(receipt.Status),
		LogsBloom:         "0x" + common.Bytes2Hex(receipt.Bloom.Bytes()),
		BlockTimestamp:    ts,
	}

	// Логи
	rec.Logs = NewLogs(receipt.Logs, ts)
	return rec
}

// NewLogs преобразует []*types.Log → []models.Log
func NewLogs(logs []*types.Log, ts time.Time) []models.Log {
	if len(logs) == 0 {
		return nil
	}
	result := make([]models.Log, len(logs))
	for i, l := range logs {
		result[i] = models.Log{
			BlockNumber:      uint(l.BlockNumber),
			BlockHash:        l.BlockHash.Hex(),
			TransactionHash:  l.TxHash.Hex(),
			TransactionIndex: uint(l.TxIndex),
			LogIndex:         uint(l.Index),
			Address:          l.Address.Hex(),
			Data:             "0x" + common.Bytes2Hex(l.Data),
			Topics:           topicsToStrings(l.Topics),
			BlockTimestamp:   ts,
			Topic0:           firstTopic(l.Topics),
		}
	}
	return result
}

func addressToOptionalString(addr *common.Address) *string {
	if addr == nil {
		return nil
	}
	s := addr.Hex()
	return &s
}

func topicsToStrings(topics []common.Hash) []string {
	res := make([]string, len(topics))
	for i, t := range topics {
		res[i] = t.Hex()
	}
	return res
}

func getTxSender(tx *types.Transaction) common.Address {
	var from common.Address
	if chainID := tx.ChainId(); chainID != nil && chainID.Sign() > 0 {
		signer := types.LatestSignerForChainID(chainID)
		if sender, err := types.Sender(signer, tx); err == nil {
			from = sender
		} else {
			log.Printf("Failed to get sender for tx %s: %v", tx.Hash().Hex(), err)
		}
	} else {
		if sender, err := types.Sender(types.HomesteadSigner{}, tx); err == nil {
			from = sender
		} else {
			log.Printf("Failed to get sender for legacy tx %s: %v", tx.Hash().Hex(), err)
		}
	}
	return from
}
