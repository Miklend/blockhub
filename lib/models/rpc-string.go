package models

type BlockRPCDTO struct {
	Hash             string     `json:"hash"`
	Number           string     `json:"number"`
	ParentHash       string     `json:"parentHash"`
	Nonce            string     `json:"nonce" `
	Sha3Uncles       string     `json:"sha3Uncles"`
	LogsBloom        string     `json:"logsBloom"`
	TransactionsRoot string     `json:"transactionsRoot"`
	StateRoot        string     `json:"stateRoot"`
	ReceiptsRoot     string     `json:"receiptsRoot"`
	Miner            string     `json:"miner"`
	Difficulty       string     `json:"difficulty"`
	TotalDifficulty  string     `json:"totalDifficulty"`
	Size             string     `json:"size"`
	ExtraData        string     `json:"extraData"`
	GasLimit         string     `json:"gasLimit"`
	GasUsed          string     `json:"gasUsed"`
	BaseFeePerGas    *string    `json:"baseFeePerGas,omitempty"`
	Timestamp        string     `json:"timestamp"`
	MixHash          string     `json:"mixHash"`
	Transactions     []TxRpcDTO `json:"transactions"`
	Uncles           []string   `json:"uncles" `
}

type TxRpcDTO struct {
	Hash                 string                  `json:"hash"`
	BlockHash            string                  `json:"blockHash"`
	BlockNumber          string                  `json:"blockNumber"`
	TransactionIndex     string                  `json:"transactionIndex"`
	From                 string                  `json:"from"`
	To                   *string                 `json:"to,omitempty"`
	Value                string                  `json:"value"`
	Gas                  string                  `json:"gas" `
	GasPrice             string                  `json:"gasPrice"`
	Input                string                  `json:"input" `
	Nonce                string                  `json:"nonce"`
	Type                 string                  `json:"type" `
	MaxFeePerGas         *string                 `json:"maxFeePerGas,omitempty"`
	MaxPriorityFeePerGas *string                 `json:"maxPriorityFeePerGas,omitempty"`
	ChainID              string                  `json:"chainId"`
	V                    string                  `json:"v"`
	R                    string                  `json:"r"`
	S                    string                  `json:"s"`
	AccessList           []AccessListEntryRpcDTO `json:"accessList"`
}
type AccessListEntryRpcDTO struct {
	Address     string   `json:"address"`
	StorageKeys []string `json:"storageKeys"`
}

type ReceiptRpcDTO struct {
	TransactionHash   string      `json:"transactionHash"`
	TransactionIndex  string      `json:"transactionIndex"`
	BlockHash         string      `json:"blockHash"`
	BlockNumber       string      `json:"blockNumber"`
	From              string      `json:"from"`
	To                *string     `json:"to,omitempty"`
	ContractAddress   *string     `json:"contractAddress,omitempty"`
	CumulativeGasUsed string      `json:"cumulativeGasUsed"`
	GasUsed           string      `json:"gasUsed" `
	EffectiveGasPrice string      `json:"effectiveGasPrice"`
	Status            string      `json:"status" `
	LogsBloom         string      `json:"logsBloom"`
	Logs              []LogRpcDTO `json:"logs"`
}
type LogRpcDTO struct {
	BlockNumber      string   `json:"blockNumber"`
	BlockHash        string   `json:"blockHash"`
	TransactionHash  string   `json:"transactionHash"`
	TransactionIndex string   `json:"transactionIndex"`
	LogIndex         string   `json:"logIndex"`
	Address          string   `json:"address"`
	Data             string   `json:"data"`
	Topics           []string `json:"topics"`
	Removed          bool     `json:"removed"`
}
