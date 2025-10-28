package rowtypes

import "time"

type BlockRow struct {
	Hash             string    `ch:"hash"`
	Number           uint64    `ch:"number"`
	ParentHash       string    `ch:"parent_hash"`
	Nonce            string    `ch:"nonce"`
	Sha3Uncles       string    `ch:"sha3_uncles"`
	LogsBloom        string    `ch:"logs_bloom"`
	TransactionsRoot string    `ch:"transactions_root"`
	StateRoot        string    `ch:"state_root"`
	ReceiptsRoot     string    `ch:"receipts_root"`
	Miner            string    `ch:"miner"`
	Difficulty       string    `ch:"difficulty"`
	TotalDifficulty  string    `ch:"total_difficulty"`
	Size             uint64    `ch:"size"`
	ExtraData        string    `ch:"extra_data"`
	GasLimit         uint64    `ch:"gas_limit"`
	GasUsed          uint64    `ch:"gas_used"`
	BaseFeePerGas    *uint64   `ch:"base_fee_per_gas"`
	Timestamp        time.Time `ch:"timestamp"`
	MixHash          string    `ch:"mix_hash"`
	Transactions     []string  `ch:"transactions"`
	Uncles           []string  `ch:"uncles"`
}

type TxRow struct {
	Hash                 string    `ch:"hash"`
	BlockHash            string    `ch:"block_hash"`
	BlockNumber          uint64    `ch:"block_number"`
	TransactionIndex     uint32    `ch:"transaction_index"`
	From                 string    `ch:"from"`
	To                   *string   `ch:"to"`
	Value                string    `ch:"value"`
	Gas                  uint64    `ch:"gas"`
	GasPrice             uint64    `ch:"gas_price"`
	Input                string    `ch:"input"`
	Nonce                uint64    `ch:"nonce"`
	Type                 uint8     `ch:"type"`
	MaxFeePerGas         *uint64   `ch:"max_fee_per_gas"`
	MaxPriorityFeePerGas *uint64   `ch:"max_priority_fee_per_gas"`
	ChainID              uint64    `ch:"chain_id"`
	V                    string    `ch:"v"`
	R                    string    `ch:"r"`
	S                    string    `ch:"s"`
	AccessList           string    `ch:"access_list"`
	BlockTimestamp       time.Time `ch:"block_timestamp"`
}

type ReceiptRow struct {
	TransactionHash   string    `ch:"transaction_hash"`
	TransactionIndex  uint32    `ch:"transaction_index"`
	BlockHash         string    `ch:"block_hash"`
	BlockNumber       uint64    `ch:"block_number"`
	From              string    `ch:"from"`
	To                *string   `ch:"to"`
	ContractAddress   *string   `ch:"contract_address"`
	CumulativeGasUsed uint64    `ch:"cumulative_gas_used"`
	GasUsed           uint64    `ch:"gas_used"`
	EffectiveGasPrice uint64    `ch:"effective_gas_price"`
	Status            uint8     `ch:"status"`
	LogsBloom         string    `ch:"logs_bloom"`
	BlockTimestamp    time.Time `ch:"block_timestamp"`
}

type LogRow struct {
	BlockNumber      uint64    `ch:"block_number"`
	BlockHash        string    `ch:"block_hash"`
	TransactionHash  string    `ch:"transaction_hash"`
	TransactionIndex uint32    `ch:"transaction_index"`
	LogIndex         uint32    `ch:"log_index"`
	Address          string    `ch:"address"`
	Data             string    `ch:"data"`
	Topics           []string  `ch:"topics"`
	BlockTimestamp   time.Time `ch:"block_timestamp"`
	Topic0           string    `ch:"topic0"`
}
