package models

// Блок (eth_getBlockByNumber / eth_getBlockByHash)
type Block struct {
	BaseFeePerGas    string   `json:"baseFeePerGas,omitempty"`
	Difficulty       string   `json:"difficulty"`
	ExtraData        string   `json:"extraData"`
	GasLimit         uint64   `json:"gasLimit"`
	GasUsed          uint64   `json:"gasUsed"`
	Hash             string   `json:"hash"`
	LogsBloom        string   `json:"logsBloom"`
	Miner            string   `json:"miner"`
	MixHash          string   `json:"mixHash"`
	Nonce            uint64   `json:"nonce"`
	Number           uint64   `json:"number"`
	ParentHash       string   `json:"parentHash"`
	ReceiptsRoot     string   `json:"receiptsRoot"`
	Sha3Uncles       string   `json:"sha3Uncles"`
	Size             uint64   `json:"size"`
	StateRoot        string   `json:"stateRoot"`
	Timestamp        uint64   `json:"timestamp"`
	Transactions     []Tx     `json:"transactions"`
	TransactionsRoot string   `json:"transactionsRoot"`
	Uncles           []string `json:"uncles"`
}

// Транзакция (eth_getTransactionByHash / eth_getBlockBy*)
type Tx struct {
	BlockHash            string   `json:"blockHash,omitempty"`
	From                 string   `json:"from"`
	Gas                  uint64   `json:"gas"`
	GasPrice             string   `json:"gasPrice"`
	MaxFeePerGas         string   `json:"maxFeePerGas,omitempty"`
	MaxPriorityFeePerGas string   `json:"maxPriorityFeePerGas,omitempty"`
	Hash                 string   `json:"hash"`
	Input                string   `json:"input"`
	Nonce                uint64   `json:"nonce"`
	To                   string   `json:"to,omitempty"`
	TransactionIndex     uint64   `json:"transactionIndex"`
	Value                string   `json:"value"`
	Type                 uint8    `json:"type"`
	ChainID              string   `json:"chainId,omitempty"`
	V                    string   `json:"v"`
	R                    string   `json:"r"`
	S                    string   `json:"s"`
	Receipt              *Receipt `json:"receipt,omitempty"`
}

// Квитанция (eth_getTransactionReceipt / eth_getBlockReceipt)
type Receipt struct {
	TxHash            string `json:"transactionHash"`
	ContractAddress   string `json:"contractAddress,omitempty"`
	CumulativeGasUsed uint64 `json:"cumulativeGasUsed"`
	EffectiveGasPrice string `json:"effectiveGasPrice"`
	From              string `json:"from"`
	GasUsed           uint64 `json:"gasUsed"`
	Logs              []Log  `json:"logs"`
	LogsBloom         string `json:"logsBloom"`
	Status            uint64 `json:"status"`
	To                string `json:"to,omitempty"`
	Type              uint8  `json:"type"`
}

// Лог (event Log)
type Log struct {
	BlockHash       string   `json:"blockHash,omitempty"`
	Address         string   `json:"address"`
	Topics          []string `json:"topics"`
	Data            string   `json:"data"`
	TransactionHash string   `json:"transactionHash"`
	LogIndex        uint64   `json:"logIndex"`
	Removed         bool     `json:"removed"`
}
