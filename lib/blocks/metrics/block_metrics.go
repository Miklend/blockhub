package metrics

import (
	"lib/models"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// NewBlockMetrics преобразует блок и квитанции из go-ethereum в нашу структуру models.Block
func NewBlockMetrics(block *types.Block, receipts []*types.Receipt) models.Block {
	if block == nil {
		return models.Block{}
	}

	// Преобразуем uncle хэши в строки
	uncles := make([]string, len(block.Uncles()))
	for i, u := range block.Uncles() {
		uncles[i] = u.Hash().Hex()
	}

	return models.Block{
		BaseFeePerGas:    bigIntToString(block.BaseFee()),
		Difficulty:       bigIntToString(block.Difficulty()),
		ExtraData:        string(block.Extra()),
		GasLimit:         block.GasLimit(),
		GasUsed:          block.GasUsed(),
		Hash:             block.Hash().Hex(),
		LogsBloom:        common.Bytes2Hex(block.Bloom().Bytes()),
		Miner:            block.Coinbase().Hex(),
		MixHash:          block.MixDigest().Hex(),
		Nonce:            block.Nonce(),
		Number:           block.NumberU64(),
		ParentHash:       block.ParentHash().Hex(),
		ReceiptsRoot:     block.ReceiptHash().Hex(),
		Sha3Uncles:       block.UncleHash().Hex(),
		Size:             block.Size(),
		StateRoot:        block.Root().Hex(),
		Timestamp:        block.Time(),
		Transactions:     NewTxs(block.Transactions(), receipts),
		TransactionsRoot: block.TxHash().Hex(),
		Uncles:           uncles,
	}
}

// NewTxs преобразует транзакции и их квитанции
func NewTxs(txs types.Transactions, receipts []*types.Receipt) []models.Tx {
	if txs == nil {
		return nil
	}

	result := make([]models.Tx, 0, len(txs))
	for i, tx := range txs {
		var receipt *types.Receipt
		if i < len(receipts) {
			receipt = receipts[i]
		}

		v, r, s := tx.RawSignatureValues()

		from := getTxSender(tx)

		var chainID string
		if tx.ChainId() != nil {
			chainID = tx.ChainId().String()
		}

		result = append(result, models.Tx{
			From:                 from.Hex(),
			Gas:                  tx.Gas(),
			GasPrice:             bigIntToString(tx.GasPrice()),
			MaxFeePerGas:         bigIntToString(tx.GasFeeCap()),
			MaxPriorityFeePerGas: bigIntToString(tx.GasTipCap()),
			Hash:                 tx.Hash().Hex(),
			Input:                common.Bytes2Hex(tx.Data()),
			Nonce:                tx.Nonce(),
			To:                   addressToString(tx.To()),
			TransactionIndex:     uint64(i),
			Value:                bigIntToString(tx.Value()),
			Type:                 tx.Type(),
			ChainID:              chainID,
			V:                    v.String(),
			R:                    r.String(),
			S:                    s.String(),
			Receipt:              NewReceipt(receipt),
		})
	}

	return result
}

// NewReceipt преобразует квитанцию в models.Receipt
func NewReceipt(receipt *types.Receipt) *models.Receipt {
	if receipt == nil {
		return nil
	}

	return &models.Receipt{
		ContractAddress:   receipt.ContractAddress.Hex(),
		CumulativeGasUsed: receipt.CumulativeGasUsed,
		EffectiveGasPrice: bigIntToString(receipt.EffectiveGasPrice),
		GasUsed:           receipt.GasUsed,
		Logs:              NewLogs(receipt.Logs),
		LogsBloom:         common.Bytes2Hex(receipt.Bloom.Bytes()),
		Status:            uint64(receipt.Status),
		Type:              receipt.Type,
	}
}

// NewLogs преобразует логи в []models.Log
func NewLogs(logs []*types.Log) []models.Log {
	if logs == nil {
		return nil
	}

	result := make([]models.Log, len(logs))
	for i, l := range logs {
		result[i] = models.Log{
			Address:         l.Address.Hex(),
			Topics:          topicsToStrings(l.Topics),
			Data:            common.Bytes2Hex(l.Data),
			TransactionHash: l.TxHash.Hex(),
			LogIndex:        uint64(i),
			Removed:         l.Removed,
		}
	}
	return result
}

// Вспомогательные функции

func bigIntToString(n *big.Int) string {
	if n == nil {
		return "0"
	}
	return n.String()
}

func addressToString(addr *common.Address) string {
	if addr == nil {
		return ""
	}
	return addr.Hex()
}

func topicsToStrings(topics []common.Hash) []string {
	result := make([]string, len(topics))
	for i, t := range topics {
		result[i] = t.Hex()
	}
	return result
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
