package types

import (
	"tinychain/common"
	"math/big"
)

// Receipt represents the results of a transaction
type Receipt struct {
	// Consensus fields
	PostState       common.Hash    `json:"root"`             // post state root
	Status          bool           `json:"status"`           // Transaction executing success or failed
	TxHash          common.Hash    `json:"tx_hash"`          // Transaction hash
	ContractAddress common.Address `json:"contract_address"` // Contract address
	GasUsed         *big.Int       `json:"gas_used"`         // gas used of transaction
}

func NewRecipet(root common.Hash, status bool, txHash common.Hash, gasUsed *big.Int) *Receipt {
	return &Receipt{
		PostState: root,
		Status:    status,
		TxHash:    txHash,
		GasUsed:   gasUsed,
	}
}

func (re *Receipt) SetContractAddress(addr common.Address) {
	re.ContractAddress = addr
}
