package executor

import (
	"tinychain/core/types"
	"tinychain/core/state"
	"tinychain/core"
)

type StateValidatorImpl struct {
	config    *Config
	processor core.Processor
	state     *state.StateDB
}

func NewStateValidator(config *Config, state *state.StateDB, processor core.Processor) *StateValidatorImpl {
	return &StateValidatorImpl{
		config:    config,
		state:     state,
		processor: processor,
	}
}

// Validate block state and receipts
// 1. Simulate process every transaction
// 2. Validate every tx result match the given receipts or not
// 3. Collect all valid and invalid receipts. if len(invalid receipts) > 0,
//    the whole block rollback and is dropped
func (sv *StateValidatorImpl) Process(txs types.Transactions, receipts types.Receipts) (valid types.Receipts, invalid types.Receipts) {

}
