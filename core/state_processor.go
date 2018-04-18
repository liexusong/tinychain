package core

import (
	"tinychain/consensus"
	"tinychain/core/types"
	"tinychain/core/state"
)

type StateProcessor struct {
	bc     *Blockchain
	engine consensus.Engine
}

func NewStateProcessor(bc *Blockchain, engine consensus.Engine) *StateProcessor {
	return &StateProcessor{
		bc:     bc,
		engine: engine,
	}
}

func (sp *StateProcessor) Process(block *types.Block, statedb *state.StateDB) ([]*types.Receipt, error) {

}
