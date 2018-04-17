package core

import (
	"tinychain/common"
	"math/big"
	"tinychain/core/vm"
	"errors"
	"github.com/op/go-logging"
)

var (
	ErrNonceTooHight = errors.New("nonce too hight")
	ErrNonceTooLow   = errors.New("nonce too low")
	MaxGas           = uint64(9999999)
	RewardPerTx      = uint64(100)
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
	log     *logging.Logger
}

func NewStateTransition(evm *vm.EVM, event Event) *StateTransition {
	return &StateTransition{
		evm:     evm,
		event:   event,
		statedb: evm.StateDB,
		log:     common.GetLogger("state"),
	}
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

func (st *StateTransition) to() common.Address {
	if st.event == nil {
		return common.Address{}
	}
	return st.event.To()
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
	if (st.to() == common.Address{}) {
		// Contract create
		ret, _, _, vmerr = st.evm.Create(vm.AccountRef(st.to()), st.data(), MaxGas, st.value())
	} else {
		ret, _, vmerr = st.evm.Call(vm.AccountRef(st.event.From()), st.to(), st.data(), MaxGas, st.value())
	}
	if vmerr != nil {
		log.Errorf("VM returned with error %s", vmerr)
		return nil, vmerr
	}
	st.statedb.AddBalance(st.evm.Coinbase, new(big.Int).Add(new(big.Int).SetUint64(RewardPerTx), new(big.Int).SetUint64(0)))
	return ret, nil
}
