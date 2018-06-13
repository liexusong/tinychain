package executor

import (
	"tinychain/core/types"
)

// Validator validate the block info
type Validator interface {
	BlockValidator
	TxValidator
	StateValidator
}

type TxValidator interface {
	ValidateTxs(txs types.Transactions) (types.Transactions, types.Transactions)
}

type StateValidator interface {
	Process(tx types.Transactions, receipts types.Receipts) error
	ValidReceipts() types.Receipts
	InvalidReceipts() types.Receipts
}

type BlockValidator interface {
	// Validate block header
	ValidateHeader(block *types.Block) error
	// Validate block body, including transactions, receipts
	ValidateBody(block *types.Block) error
}
