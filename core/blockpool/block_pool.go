package blockpool

import (
	"sync"
	"tinychain/event"
	"tinychain/core/types"
	"sync/atomic"
)

type BlockPool struct {
	valid sync.Map     // Valid blocks pool. map[height][]*block
	all   sync.Map     // All blocks indexed by block_hash, for quick searching
	size  atomic.Value // The size of block pool

	blockSub event.Subscription
}

func (bp *BlockPool) Start() error {

}

func (bp *BlockPool) Valid() {

}

func (bp *BlockPool) get() *types.Block {

}

func (bp *BlockPool) add() {

}

// Clear picks up the invalid blocks in the pool and removes them.
func (bp *BlockPool) Clear() {

}

func (bp *BlockPool) Stop() {

}
