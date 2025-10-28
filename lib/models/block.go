package models

import "time"

// Block — модель блока, пригодная для обмена между микросервисами и ClickHouse
type Block struct {
	Hash             string    `json:"hash" ch:"hash"`
	Number           uint      `json:"number" ch:"number"`
	ParentHash       string    `json:"parentHash" ch:"parent_hash"`
	Nonce            uint      `json:"nonce" ch:"nonce"`
	Sha3Uncles       string    `json:"sha3Uncles" ch:"sha3_uncles"`
	LogsBloom        string    `json:"logsBloom" ch:"logs_bloom"`
	TransactionsRoot string    `json:"transactionsRoot" ch:"transactions_root"`
	StateRoot        string    `json:"stateRoot" ch:"state_root"`
	ReceiptsRoot     string    `json:"receiptsRoot" ch:"receipts_root"`
	Miner            string    `json:"miner" ch:"miner"`
	Difficulty       string    `json:"difficulty" ch:"difficulty"`
	TotalDifficulty  string    `json:"totalDifficulty" ch:"total_difficulty"`
	Size             uint      `json:"size" ch:"size"`
	ExtraData        string    `json:"extraData" ch:"extra_data"`
	GasLimit         uint      `json:"gasLimit" ch:"gas_limit"`
	GasUsed          uint      `json:"gasUsed" ch:"gas_used"`
	BaseFeePerGas    *uint     `json:"baseFeePerGas,omitempty" ch:"base_fee_per_gas"`
	Timestamp        time.Time `json:"timestamp" ch:"timestamp"`
	MixHash          string    `json:"mixHash" ch:"mix_hash"`
	Transactions     []string  `json:"transactions" ch:"transactions"`
	Uncles           []string  `json:"uncles" ch:"uncles"`
}

// Tx — модель транзакции
type Tx struct {
	Hash                 string    `json:"hash" ch:"hash"`
	BlockHash            string    `json:"blockHash" ch:"block_hash"`
	BlockNumber          uint      `json:"blockNumber" ch:"block_number"`
	TransactionIndex     uint      `json:"transactionIndex" ch:"transaction_index"`
	From                 string    `json:"from" ch:"from"`
	To                   *string   `json:"to,omitempty" ch:"to"`
	Value                string    `json:"value" ch:"value"`
	Gas                  uint      `json:"gas" ch:"gas"`
	GasPrice             uint      `json:"gasPrice" ch:"gas_price"`
	Input                string    `json:"input" ch:"input"`
	Nonce                uint      `json:"nonce" ch:"nonce"`
	Type                 uint      `json:"type" ch:"type"`
	MaxFeePerGas         *uint     `json:"maxFeePerGas,omitempty" ch:"max_fee_per_gas"`
	MaxPriorityFeePerGas *uint     `json:"maxPriorityFeePerGas,omitempty" ch:"max_priority_fee_per_gas"`
	ChainID              uint      `json:"chainId" ch:"chain_id"`
	V                    string    `json:"v" ch:"v"`
	R                    string    `json:"r" ch:"r"`
	S                    string    `json:"s" ch:"s"`
	AccessList           string    `json:"accessList" ch:"access_list"`
	BlockTimestamp       time.Time `json:"blockTimestamp" ch:"block_timestamp"`
	Receipt              *Receipt  `json:"receipts" ch:"receipts"`
}

// Receipt — модель квитанции
type Receipt struct {
	TransactionHash   string    `json:"transactionHash" ch:"transaction_hash"`
	TransactionIndex  uint      `json:"transactionIndex" ch:"transaction_index"`
	BlockHash         string    `json:"blockHash" ch:"block_hash"`
	BlockNumber       uint      `json:"blockNumber" ch:"block_number"`
	From              string    `json:"from" ch:"from"`
	To                *string   `json:"to,omitempty" ch:"to"`
	ContractAddress   *string   `json:"contractAddress,omitempty" ch:"contract_address"`
	CumulativeGasUsed uint      `json:"cumulativeGasUsed" ch:"cumulative_gas_used"`
	GasUsed           uint      `json:"gasUsed" ch:"gas_used"`
	EffectiveGasPrice uint      `json:"effectiveGasPrice" ch:"effective_gas_price"`
	Status            uint      `json:"status" ch:"status"`
	LogsBloom         string    `json:"logsBloom" ch:"logs_bloom"`
	BlockTimestamp    time.Time `json:"blockTimestamp" ch:"block_timestamp"`
	Logs              []Log     `json:"logs" ch:"logs"`
}

// Log — модель лога события
type Log struct {
	BlockNumber      uint      `json:"blockNumber" ch:"block_number"`
	BlockHash        string    `json:"blockHash" ch:"block_hash"`
	TransactionHash  string    `json:"transactionHash" ch:"transaction_hash"`
	TransactionIndex uint      `json:"transactionIndex" ch:"transaction_index"`
	LogIndex         uint      `json:"logIndex" ch:"log_index"`
	Address          string    `json:"address" ch:"address"`
	Data             string    `json:"data" ch:"data"`
	Topics           []string  `json:"topics" ch:"topics"`
	BlockTimestamp   time.Time `json:"blockTimestamp" ch:"block_timestamp"`
	Topic0           string    `json:"topic0" ch:"topic0"`
}
