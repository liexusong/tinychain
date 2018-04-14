package types

import "tinychain/common"

// Receipt represents the results of a transaction
type Receipt struct {
	// Consensus fields
	PostState       common.Hash    `json:"root"`             // post state root
	Status          bool           `json:"status"`           // Transaction executing success or failed
	TxHash          common.Hash    `json:"tx_hash"`          // Transaction hash
	ContractAddress common.Address `json:"contract_address"` // Contract address
}

func NewRecipet(root common.Hash, status bool, txHash common.Hash) *Receipt {
	return &Receipt{
		PostState: root,
		Status:    status,
		TxHash:    txHash,
	}
}

func (re *Receipt) SetContractAddress(addr common.Address) {
	re.ContractAddress = addr
}
