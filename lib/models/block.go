package models

type Hash string

type Address string

type Time uint64

type BlockDTO struct {
	Hash             Hash    `json:"hash"`
	Number           uint64  `json:"number"`
	ParentHash       Hash    `json:"parentHash"`
	Nonce            string  `json:"nonce"`
	Sha3Uncles       Hash    `json:"sha3Uncles"`
	LogsBloom        string  `json:"logsBloom"`
	TransactionsRoot Hash    `json:"transactionsRoot"`
	StateRoot        Hash    `json:"stateRoot"`
	ReceiptsRoot     Hash    `json:"receiptsRoot"`
	Miner            Address `json:"miner"`
	Difficulty       string  `json:"difficulty"`
	TotalDifficulty  string  `json:"totalDifficulty"`
	Size             uint64  `json:"size"`
	ExtraData        string  `json:"extraData"`
	GasLimit         uint64  `json:"gasLimit"`
	GasUsed          uint64  `json:"gasUsed"`
	BaseFeePerGas    *string `json:"baseFeePerGas,omitempty"`
	Timestamp        Time    `json:"timestamp"`
	MixHash          Hash    `json:"mixHash"`
	Transactions     []TxDTO `json:"transactions"`
	Uncles           []Hash  `json:"uncles"`
}

type TxDTO struct {
	Hash                 Hash                 `json:"hash"`
	BlockHash            Hash                 `json:"blockHash"`
	BlockNumber          uint64               `json:"blockNumber"`
	TransactionIndex     uint64               `json:"transactionIndex"`
	From                 Address              `json:"from"`
	To                   *Address             `json:"to,omitempty"`
	Value                string               `json:"value"`
	Gas                  uint64               `json:"gas"`
	GasPrice             string               `json:"gasPrice"`
	Input                string               `json:"input"`
	Nonce                uint64               `json:"nonce"`
	Type                 string               `json:"type"`
	MaxFeePerGas         *string              `json:"maxFeePerGas,omitempty"`
	MaxPriorityFeePerGas *string              `json:"maxPriorityFeePerGas,omitempty"`
	ChainID              string               `json:"chainId"`
	V                    string               `json:"v"`
	R                    string               `json:"r"`
	S                    string               `json:"s"`
	AccessList           []AccessListEntryDTO `json:"accessList"`
	Receipt              ReceiptDTO           `json:"receipt"`
}

type AccessListEntryDTO struct {
	Address     Address  `json:"address"`
	StorageKeys []string `json:"storageKeys"`
}

type ReceiptDTO struct {
	TransactionHash   Hash     `json:"transactionHash"`
	TransactionIndex  uint64   `json:"transactionIndex"`
	BlockHash         Hash     `json:"blockHash"`
	BlockNumber       uint64   `json:"blockNumber"`
	From              Address  `json:"from"`
	To                *Address `json:"to,omitempty"`
	ContractAddress   *Address `json:"contractAddress,omitempty"`
	CumulativeGasUsed uint64   `json:"cumulativeGasUsed"`
	GasUsed           uint64   `json:"gasUsed"`
	EffectiveGasPrice string   `json:"effectiveGasPrice"`
	Status            string   `json:"status"`
	LogsBloom         string   `json:"logsBloom"`
	Logs              []LogDTO `json:"logs"`
}

type LogDTO struct {
	BlockNumber      uint64   `json:"blockNumber"`
	BlockHash        Hash     `json:"blockHash"`
	TransactionHash  Hash     `json:"transactionHash"`
	TransactionIndex uint64   `json:"transactionIndex"`
	LogIndex         uint64   `json:"logIndex"`
	Address          Address  `json:"address"`
	Data             string   `json:"data"`
	Topics           []string `json:"topics"`
	Removed          bool     `json:"removed"`
}
