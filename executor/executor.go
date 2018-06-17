package executor

import (
	"tinychain/core"
	"tinychain/event"
	"tinychain/core/state"
)

// Executor handles all data modification and
// executes and validates state transition
type Executor interface {
	Start() error
	Stop() error
}

type ExecutorImpl struct {
	processor core.Processor
	validator Validator        // Validator validate all consensus fields
	chain     *core.Blockchain // Blockchain wrapper
	event     *event.TypeMux
	quitCh    chan struct{}

	blockSub event.Subscription // Subscribe new block event
	txsSub   event.Subscription // Subscribe new transactions event
}

func New(chain *core.Blockchain, statedb *state.StateDB) Executor {
	processor := core.NewStateProcessor(chain, statedb)
	executor := &ExecutorImpl{
		processor: processor,
		chain:     chain,
		event:     event.GetEventhub(),
		quitCh:    make(chan struct{}),
		validator: NewValidator(processor),
	}
	return executor
}

func (ex *ExecutorImpl) Start() error {
	ex.blockSub = ex.event.Subscribe(&core.NewBlockEvent{})
	ex.txsSub = ex.event.Subscribe(&core.NewTxsEvent{})
	go ex.listenBlock()
	go ex.listenTx()
}

func (ex *ExecutorImpl) listenBlock() {
	for {
		select {
		case ev := <-ex.blockSub.Chan():
			block := ev.(*core.NewBlockEvent).Block
			err := ex.validator.ValidateHeader(block)
			if err != nil {

			}

			err = ex.validator.ValidateBody(block)
		case ev := <-ex.txsSub.Chan():
			txs := ev.(*core.NewTxsEvent).Txs
		case <-ex.quitCh:
			ex.blockSub.Unsubscribe()
		}
	}
}

func (ex *ExecutorImpl) listenTx() {
	for {
		select {
		case ev := <-ex.txsSub.Chan():
			tx := ev.(core.NewTxEvent).Tx
			// TODO validate transaction
		case <-ex.quitCh:
			ex.txSub.Unsubscribe()
		}
	}
}

func (ex *ExecutorImpl) Stop() error {
	close(ex.quitCh)
}

func (ex *ExecutorImpl) process() {

}
