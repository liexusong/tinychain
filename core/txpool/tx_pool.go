package txpool

import (
	"tinychain/common"
	"tinychain/executor"
	"tinychain/core/types"
	"tinychain/event"
	"sync"
	"tinychain/core/state"
)

var (
	log            = common.GetLogger("txpool")
	txPool *TxPool = nil
)

type TxPool struct {
	currentState *state.StateDB
	eventHub    *event.TypeMux
	txValidator executor.TxValidator
	quitCh      chan struct{}

	// all valid and processable txs.
	// map[common.Address]*txList
	pending sync.Map

	// all new-added and non-processable txs,including valid and invalid txs.
	// map[common.Address]*txList
	queue sync.Map

	newTxSub event.Subscription
}

func NewTxPool(validator executor.TxValidator) *TxPool {
	return &TxPool{
		txValidator: validator,
		eventHub:    event.GetEventhub(),
	}
}

func (tp *TxPool) Start() {
	tp.newTxSub = tp.eventHub.Subscribe(&event.NewTxEvent{})
}

func (tp *TxPool) listen() {
	for {
		select {
		case ev := <-tp.newTxSub.Chan():
			tp.Add(ev.(*event.NewTxEvent).Tx)
		case <-tp.quitCh:
			tp.newTxSub.Unsubscribe()
			break
		}
	}
}

// Pending returns all nonce-sorted transactions of every address
func (tp *TxPool) Pending() map[common.Address]types.Transactions {
	results := make(map[common.Address]types.Transactions)
	tp.pending.Range(func(key, value interface{}) bool {
		results[key.(common.Address)] = value.(*txList).All()
		return true
	})
	return results
}

func (tp *TxPool) Add(tx *types.Transaction) error {
	if err := tp.validateTx(tx); err != nil {
		log.Errorf("failed to validate tx, %s", err)
		return err
	}
	return tp.add(tx)
}

func (tp *TxPool) add(tx *types.Transaction) error {
	// check promote
}

func (tp *TxPool) addQueue(tx *types.Transaction) {

}

func (tp *TxPool) addPending(tx *types.Transaction) {

}

func (tp *TxPool) validateTx(tx *types.Transaction) error {
	return tp.txValidator.ValidateTx(tx)
}
