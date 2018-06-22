package blockpool

import (
	"sync"
	"tinychain/event"
	"tinychain/core/types"
	"tinychain/executor"
	"tinychain/core"
	"math/big"
	"github.com/pkg/errors"
	"tinychain/common"
)

var (
	log = common.GetLogger("blockpool")

	ErrBlockDuplicate = errors.New("block duplicate")
	ErrPoolFull       = errors.New("block pool is full")
)

type BlockPool struct {
	config     *Config
	mu         sync.RWMutex
	blockchain *core.Blockchain          // Blockchain
	validator  executor.BlockValidator   // Block validator
	valid      map[*big.Int]*types.Block // Valid blocks pool. map[height]*block
	event      *event.TypeMux
	quitCh     chan struct{}

	blockSub  event.Subscription
	commitSub event.Subscription // Receive msg when new blocks are appended to blockchain
}

func NewBlockPool(config *Config, blockchain *core.Blockchain, validator executor.BlockValidator) *BlockPool {
	return &BlockPool{
		config:     config,
		event:      event.GetEventhub(),
		validator:  validator,
		blockchain: blockchain,
		valid:      make(map[*big.Int]*types.Block, config.MaxBlockSize),
	}
}

func (bp *BlockPool) Start() {
	bp.blockSub = bp.event.Subscribe(&core.NewBlockEvent{})
	go bp.listen()
}

func (bp *BlockPool) listen() {
	for {
		select {
		case ev := <-bp.blockSub.Chan():
			block := ev.(*core.NewBlockEvent).Block
			go bp.add(block)
		case ev := <-bp.commitSub.Chan():
			commit := ev.(*core.BlockCommitEvent)
			go bp.del(commit.Height)
		case <-bp.quitCh:
			bp.blockSub.Unsubscribe()
			return
		}
	}
}

func (bp *BlockPool) Valid() []*types.Block {
	var blocks []*types.Block
	bp.mu.RLock()
	defer bp.mu.RUnlock()
	for _, block := range bp.valid {
		blocks = append(blocks, block)
	}
	return blocks
}

func (bp *BlockPool) add(block *types.Block) error {
	bp.mu.RLock()
	// Check block duplicate
	if old := bp.valid[block.Height()]; old != nil {
		// TODO Check can be replace or not?
		return ErrBlockDuplicate
	}
	if bp.Size() >= bp.config.MaxBlockSize {
		return ErrPoolFull
	}
	bp.mu.RUnlock()

	// Validate block
	if err := bp.validate(block); err != nil {
		return err
	}

	bp.mu.Lock()
	defer bp.mu.Unlock()
	bp.valid[block.Height()] = block
	return nil
}

func (bp *BlockPool) validate(block *types.Block) error {
	err := bp.validator.ValidateHeader(block)
	if err != nil {
		return err
	}

	return bp.validator.ValidateBody(block)
}

// Clear picks up the invalid blocks in the pool and removes them.
func (bp *BlockPool) del(height *big.Int) {
	bp.mu.Lock()
	defer bp.mu.Unlock()
	delete(bp.valid, height)
}

// Size gets the size of valid blocks.
// The caller should hold the lock before invoke this func.
func (bp *BlockPool) Size() uint64 {
	return uint64(len(bp.valid))
}

func (bp *BlockPool) Stop() {
	close(bp.quitCh)
}
