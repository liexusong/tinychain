package core

import (
	"tinychain/db"
	"github.com/hashicorp/golang-lru"
	"tinychain/core/types"
	"tinychain/consensus"
	"tinychain/common"
	"sync/atomic"
)

var (
	cacheSize = 65535
	log       = common.GetLogger("blockchain")
)

// Blockchain is the canonical chain given a database with a genesis block
type Blockchain struct {
	db        *db.TinyDB       // chain db
	lastBlock atomic.Value     // last block of chain
	engine    consensus.Engine // consensus engine
	// TODO more fields

	blocksCache *lru.Cache // blocks lru cache
	headerCache *lru.Cache // headers lru cache
}

func NewBlockchain(db *db.TinyDB, engine consensus.Engine) (*Blockchain, error) {
	blocksCache, _ := lru.New(cacheSize)
	headerCache, _ := lru.New(cacheSize)
	bc := &Blockchain{
		db:          db,
		engine:      engine,
		blocksCache: blocksCache,
		headerCache: headerCache,
	}
	if err := bc.loadLastState(); err != nil {
		return nil, err
	}
	return bc, nil
}

func (bc *Blockchain) loadLastState() error {
	lastBlock, err := bc.db.GetLastBlock()
	if err != nil {
		return err
	}
	bc.lastBlock.Store(lastBlock)
	bc.blocksCache.Add(lastBlock.Hash, lastBlock)

	// TODO

	return nil
}

func (bc *Blockchain) GetLastBlock() *types.Block {
	if block := bc.lastBlock.Load(); block != nil {
		return block.(*types.Block)
	}
	return nil
}

func (bc *Blockchain) GetBlock(hash common.Hash) (*types.Block, error) {
	if block, ok := bc.blocksCache.Get(hash); ok {
		return block.(*types.Block), nil
	}
	height, err := bc.db.GetHeight(hash)
	if err != nil {
		return nil, err
	}
	block, err := bc.db.GetBlock(height, hash)
	if err != nil {
		return nil, err
	}
	bc.blocksCache.Add(hash, block)
	return block, nil
}

func (bc *Blockchain) GetHeader(hash common.Hash) (*types.Header, error) {
	if header, ok := bc.headerCache.Get(hash); ok {
		return header.(*types.Header), nil
	}
	height, err := bc.db.GetHeight(hash)
	if err != nil {
		return nil, err
	}
	header, err := bc.db.GetHeader(height, hash)
	if err != nil {
		return nil, err
	}
	bc.headerCache.Add(hash, header)
	return header, nil
}

func (bc *Blockchain) AddBlock() error {

}

// Commit the blockchain to db
func (bc *Blockchain) Commit() error {

}

func (bc *Blockchain) Engine() consensus.Engine {

}
