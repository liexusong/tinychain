package core

import (
	"tinychain/common"
	"math/big"
	"tinychain/core/vm"
	"errors"
)

var (
	ErrNonceTooHight = errors.New("nonce too hight")
	ErrNonceTooLow   = errors.New("nonce too low")
	MaxGas           = uint64(9999999) // Maximum
	RewardPerTx      = uint64(100)     // Reward token per tx for miner
)

// An event sent to a contract, which will make state transition
type Event interface {
	From() common.Address
	To() common.Address
	Value() *big.Int

	Nonce() uint64
	Data() []byte
}

type StateTransition struct {
	event   Event // state transition event
	evm     *vm.EVM
	statedb vm.StateDB
}

func NewStateTransition(evm *vm.EVM, event Event) *StateTransition {
	return &StateTransition{
		evm:     evm,
		event:   event,
		statedb: evm.StateDB,
	}
}

// Make state transition by applying a new event
func ApplyEvent(evm *vm.EVM, event Event) ([]byte, error) {
	return NewStateTransition(evm, event).Process()
}

func (st *StateTransition) value() *big.Int {
	return st.event.Value()
}

func (st *StateTransition) data() []byte {
	return st.event.Data()
}

// Check nonce is correct or not
// nonce should be equal to that of state object
func (st *StateTransition) preCheck() error {
	nonce := st.statedb.GetNonce(st.event.From())
	if nonce < st.event.Nonce() {
		return ErrNonceTooHight
	} else if nonce > st.event.Nonce() {
		return ErrNonceTooLow
	}
	return nil
}

func (st *StateTransition) from() vm.AccountRef {
	addr := st.event.From()
	if !st.statedb.Exist(addr) {
		st.statedb.CreateAccount(addr)
	}
	return vm.AccountRef(addr)
}

func (st *StateTransition) to() vm.AccountRef {
	if st.event == nil {
		return vm.AccountRef{}
	}

	if st.event.To().Nil() {
		return vm.AccountRef{}
	}
	to := st.event.To()
	if !st.statedb.Exist(to) {
		st.statedb.CreateAccount(to)
	}
	return vm.AccountRef(to)
}

// Make state transition according to transaction event
// NOTE: In tinychain we need not use gas to pay a transaction
func (st *StateTransition) Process() ([]byte, error) {
	if err := st.preCheck(); err != nil {
		return nil, err
	}

	var (
		vmerr error
		ret   []byte
	)
	if (st.to() == vm.AccountRef{}) {
		// Contract create
		ret, _, _, vmerr = st.evm.Create(st.to(), st.data(), MaxGas, st.value())
	} else {
		// Call contract
		ret, _, vmerr = st.evm.Call(st.from(), st.to().Address(), st.data(), MaxGas, st.value())
	}
	if vmerr != nil {
		log.Errorf("VM returned with error %s", vmerr)
		return nil, vmerr
	}
	st.statedb.AddBalance(st.evm.Coinbase, new(big.Int).Add(new(big.Int).SetUint64(RewardPerTx), new(big.Int).SetUint64(0)))
	return ret, nil
}
