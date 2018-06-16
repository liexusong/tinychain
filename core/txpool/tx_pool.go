package txpool

import (
	"tinychain/common"
	"tinychain/executor"
	"tinychain/core/types"
	"tinychain/event"
	"sync"
	"tinychain/core/state"
	"errors"
)

var (
	log            = common.GetLogger("txpool")
	txPool *TxPool = nil

	ErrTxDuplicate = errors.New("transaction duplicate")
	ErrPoolFull    = errors.New("tx_pool is full")
	ErrTxDiscard   = errors.New("old transaction is better, discard the new one")
)

type TxPool struct {
	config       *Config              // Txpool config
	currentState *state.StateDB       // Current state
	txValidator  executor.TxValidator // Tx validator wrapper
	all          *txLookup            // Cache all tx hash to accelerate searching
	eventHub     *event.TypeMux
	quitCh       chan struct{}

	// all valid and processable txs.
	// map[common.Address]*txList
	pending sync.Map

	// all new-added and non-processable txs,including valid and invalid txs.
	// map[common.Address]*txList
	queue sync.Map

	newTxSub event.Subscription
}

func NewTxPool(config *Config, validator executor.TxValidator) *TxPool {
	return &TxPool{
		config:      config,
		txValidator: validator,
		eventHub:    event.GetEventhub(),
		all:         newTxLookup(),
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
	// Validate tx
	if err := tp.validateTx(tx); err != nil {
		log.Errorf("Validate tx failed, %s", err)
		return err
	}

	// check txpool queue is full or not
	if tp.all.Len() >= tp.config.MaxTxSize {
		log.Warning(ErrPoolFull.Error())
		return ErrPoolFull
	}

	// Check whether to replace a pending tx
	replace, old := tp.replacePending(tx)
	if replace {
		log.Errorf("replace an old pending tx %s", old.Hash())
		return nil
	}

	// Add queue
	err := tp.addQueue(tx)
	if err != nil {
		return err
	}

	// Check processable
	return tp.activate()
}

func (tp *TxPool) checkProcessable() {

}

func (tp *TxPool) addQueue(tx *types.Transaction) error {
	list, exist := tp.queue.Load(tx.From)
	if !exist {
		tl := newTxList()
		tl.Add(tx, tp.config.PriceBump)
		tp.queue.Store(tx.From, tl)
		return nil
	}
	tl := list.(*txList)
	inserted, _ := tl.Add(tx, tp.config.PriceBump)
	if !inserted {
		return ErrTxDiscard
	}

	// Check tx is existed in pool or not
	if !tp.all.Get(tx.Hash()) {
		tp.all.Add(tx.Hash())
	}

	return nil
}

// replacePending check whether to replace tx in pending list,
// and if yes, return true
func (tp *TxPool) replacePending(tx *types.Transaction) (bool, *types.Transaction) {
	list, exist := tp.pending.Load(tx.From)
	if !exist {
		return false, nil
	}
	tl := list.(*txList)
	canReplace, old := tl.CanInsert(tx, tp.config.PriceBump)
	if canReplace {
		tl.Put(tx)
	}

	return canReplace, old
}

// activate moves transaction that have become processable from
// the queue to the pending list. During this process, all
// invalid transactions (low nonce, low balance) are deleted.
//
// 1. drop all low-nonce transaction
// 2. drop all costly transaction
// 3. Get sequentially increasing list and activate them
func (tp *TxPool) activate() error {

}

func (tp *TxPool) validateTx(tx *types.Transaction) error {
	return tp.txValidator.ValidateTx(tx)
}
