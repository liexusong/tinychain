package executor

import (
	"tinychain/core/types"
	"tinychain/core"
)

type ValidateImpl struct {
	processor core.Processor // Block processor
}

func NewValidator(processor core.Processor) *ValidateImpl {
	return &ValidateImpl{
		processor: processor,
	}
}

// Validate block header
// 1. Validate timestamp
// 2. Validate gasUsed and gasLimit
// 3. Validate parentHash and height
// 4. Validate extra data size is within bounds
func (v *ValidateImpl) ValidateHeader(block *types.Block) error {

}

// Validate block txs
// 1. Validate txs root hash
// 2. Validate receipts root hash
func (v *ValidateImpl) ValidateBody(block *types.Block) error {
}

// Validate block state and receipts
// 1. Simulate process every transaction
// 2. Validate the result receipt is match or not
// 3.
func (v *ValidateImpl) Process(txs types.Transactions, receipts types.Receipts) error {

}