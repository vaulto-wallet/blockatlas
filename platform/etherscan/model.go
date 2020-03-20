package etherscan

import (
	"github.com/trustwallet/blockatlas/pkg/blockatlas"
)

type Page struct {
	Total uint  `json:"total"`
	Docs  []Doc `json:"result"`
}

type TokenPage struct {
	Total uint       `json:"total"`
	Docs  []Contract `json:"docs"`
}

type BlockNumberPage struct {
	BlockNumber int64 `json:"result"`
}

type BlockPage struct {
	Block Block `json:"result"`
}

type BalancePage struct {
	Balance string `json:"result"`
}

type Block struct {
	Miner string `json:"miner"`
	Docs  []Doc  `json:"transactions"`
}

type Doc struct {
	Ops         []Op   `json:"operations"`
	Contract    string `json:"contract"`
	ID          string `json:"id"`
	BlockNumber string `json:"blockNumber"`
	Timestamp   string `json:"timeStamp"`
	Nonce       string `json:"nonce"`
	From        string `json:"from"`
	To          string `json:"to"`
	Value       string `json:"value"`
	Gas         string `json:"gas"`
	GasPrice    string `json:"gasPrice"`
	GasUsed     string `json:"gasUsed"`
	Input       string `json:"input"`
	Error       string `json:"error"`
	Coin        uint   `json:"coin"`
}

type Op struct {
	TxID     string                     `json:"transactionId"`
	Contract *Contract                  `json:"contract"`
	From     string                     `json:"from"`
	To       string                     `json:"to"`
	Type     blockatlas.TransactionType `json:"type"`
	Value    string                     `json:"value"`
	Coin     uint                       `json:"coin"`
}

type Contract struct {
	Address     string `json:"address"`
	Symbol      string `json:"symbol"`
	Decimals    uint   `json:"decimals"`
	TotalSupply string `json:"totalSupply,omitempty"`
	Name        string `json:"name"`
}

type NodeInfo struct {
	LatestBlock string `json:"result"`
}
