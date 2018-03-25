package types

import "tinychain/common"

type Receipt struct {
	// Consensus fields
	PrevState common.Hash    `json:"prev_state"` // prev state root
	TxHash    common.Hash    `json:"tx_hash"`    // Transaction hash
	Address   common.Address `json:"address"`    // Contract address
}

func NewRecipet() *Receipt {

}
