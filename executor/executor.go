package executor

import (
	"tinychain/core"
	"tinychain/event"
	"tinychain/core/state"
	"tinychain/core/types"
	batcher "github.com/yyh1102/go-batcher"
)

// Processor represents the interface of block processor
type Processor interface {
	Process(block *types.Block) (types.Receipts, error)
}

type Executor struct {
	processor Processor
	chain     *core.Blockchain // Blockchain wrapper
	batch     batcher.Batch    // Batch for creating new block
	event     *event.TypeMux
	quitCh    chan struct{}

	execblockSub event.Subscription // Subscribe new block event
	execTxsSub   event.Subscription // Execute pending txs event
}

func New(chain *core.Blockchain, statedb *state.StateDB) *Executor {
	processor := core.NewStateProcessor(chain, statedb)
	executor := &Executor{
		processor: processor,
		chain:     chain,
		event:     event.GetEventhub(),
		quitCh:    make(chan struct{}),
	}
	return executor
}

func (ex *Executor) Start() error {
	ex.execblockSub = ex.event.Subscribe(&core.ExecBlockEvent{})
	ex.execTxsSub = ex.event.Subscribe(&core.ExecPendingTxEvent{})
	go ex.listenBlock()
}

func (ex *Executor) listenBlock() {
	for {
		select {
		case ev := <-ex.execblockSub.Chan():
			block := ev.(*core.ExecBlockEvent).Block
			go ex.processBlock(block)
		case ev := <-ex.execTxsSub.Chan():
			txs := ev.(*core.ExecPendingTxEvent).Txs
			go ex.processTx(txs)
		case <-ex.quitCh:
			ex.execTxsSub.Unsubscribe()
			return
		}
	}
}

func (ex *Executor) Stop() error {
	close(ex.quitCh)
	return nil
}

func (ex *Executor) genNewBlock(txs types.Transactions, receipts types.Receipts) (*types.Block, error) {

}

func (ex *Executor) processBlock(block *types.Block) {
	receipts, err := ex.processor.Process(block)
}

// processTx execute transactions launched from tx_pool.
// 1. Simulate execute every transaction sequentially, until gasUsed reaches blocks's gasLimit
// 2. Collect valid txs and invalid txs
// 3. Collect receipts (remove invalid receipts)
func (ex *Executor) processTx(txs types.Transactions) {
}
