package core

import (
	"tinychain/consensus"
	"tinychain/core/types"
	"tinychain/core/state"
	"tinychain/common"
	"tinychain/core/vm"
)

type StateProcessor struct {
	bc      *Blockchain
	engine  consensus.Engine
	statedb *state.StateDB
}

func NewStateProcessor(bc *Blockchain, engine consensus.Engine, statedb *state.StateDB) *StateProcessor {
	return &StateProcessor{
		bc:      bc,
		engine:  engine,
		statedb: statedb,
	}
}

// Process apply transaction in state
func (sp *StateProcessor) Process(block *types.Block) ([]*types.Receipt, error) {
	var (
		receipts []*types.Receipt
		header   = block.Header
	)

	for i, tx := range block.Transactions {
		receipt, err := ApplyTransaction(sp.bc, nil, sp.statedb, header, tx)
		if err != nil {
			return nil, nil
		}
		receipts = append(receipts, receipt)
	}
	// TODO Finalize block
	return receipts, nil
}

func ApplyTransaction(bc *Blockchain, author *common.Address, statedb *state.StateDB, header *types.Header, tx *types.Transaction) (*types.Receipt, error) {
	event := tx.AsEvent()

	// Create a new context to be used in the EVM environment
	context := NewEVMContext(event, header, bc, author)
	// Create a new environment which holds all relevant information
	// about the transaction and calling mechanisms
	vm := vm.NewEVM(context, statedb, nil, nil)
	// Apply the tx to current state
	_, err := ApplyEvent(vm, event)
	if err != nil {
		return nil, err
	}
	// Get intermediate root of current state
	root, err := statedb.IntermediateRoot()
	if err != nil {
		return nil, err
	}
	receipt := types.NewRecipet(root, false, tx.Hash())
	if event.To().Nil() {
		// Create contract call
		receipt.SetContractAddress(common.CreateAddress(event.From(), event.Nonce()))
	}

	return receipt, nil
}
