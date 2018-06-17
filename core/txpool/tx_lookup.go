package txpool

import (
	"tinychain/common"
	"sync"
	"sync/atomic"
)

type txLookup struct {
	all   sync.Map
	size  atomic.Value // length cache
}

func newTxLookup() *txLookup {
	return &txLookup{}
}

func (tl *txLookup) Len() uint64 {
	if size := tl.size.Load(); size != nil {
		return size.(uint64)
	}
	var length uint64
	tl.all.Range(func(key, value interface{}) bool {
		length++
		return true
	})
	tl.size.Store(length)
	return length
}

func (tl *txLookup) Get(hash common.Hash) bool {
	_, exist := tl.all.Load(hash)
	return exist
}

func (tl *txLookup) Add(hash common.Hash) {
	tl.all.Store(hash, struct{}{})
	tl.size.Store(nil)
}

func (tl *txLookup) Del(hash common.Hash) {
	if _, exist := tl.all.Load(hash); exist {
		tl.size.Store(nil)
	}
	tl.all.Delete(hash)
}
