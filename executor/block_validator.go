package executor

import (
	"tinychain/core/types"
	"tinychain/core"
)

type BlockValidatorImpl struct {
	config *Config
	chain  *core.Blockchain
}

func NewBlockValidator() *BlockValidatorImpl {

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
