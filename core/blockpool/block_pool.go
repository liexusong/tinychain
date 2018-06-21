package blockpool

import (
	"sync"
	"tinychain/event"
	"tinychain/core/types"
	"sync/atomic"
	"tinychain/executor"
	"tinychain/core"
)

type BlockPool struct {
	validator executor.BlockValidator // Block validator
	valid     sync.Map                // Valid blocks pool. map[height][]*block
	all       sync.Map                // All blocks indexed by block_hash, for quick searching
	size      atomic.Value            // The size of block pool
	event     *event.TypeMux

	blockSub event.Subscription
}

func NewBlockPool(config *Config, validator executor.BlockValidator) *BlockPool {
	return &BlockPool{
		event:     event.GetEventhub(),
		validator: validator,
	}
}

func (bp *BlockPool) Start() error {
	bp.blockSub = bp.event.Subscribe(&core.NewBlockEvent{})
}

func (bp *BlockPool) listen() {
	for {
		select {
		case ev := <-bp.blockSub.Chan():
			block := ev.(*core.NewBlockEvent).Block
			go bp.add(block)
		}
	}
}

func (bp *BlockPool) Valid() {

}

func (bp *BlockPool) get() *types.Block {

}

func (bp *BlockPool) add(block *types.Block) {

}

// Clear picks up the invalid blocks in the pool and removes them.
func (bp *BlockPool) Clear() {

}

func (bp *BlockPool) Stop() {

}
