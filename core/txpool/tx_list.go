package txpool

import (
	"tinychain/core/types"
	"sync"
	"sort"
)

type txList struct {
	txs   sync.Map
	cache types.Transactions
}

func (list *txList) Get(nonce uint64) *types.Transaction {
	if tx, exist := list.txs.Load(nonce); exist {
		return tx.(*types.Transaction)
	}
	return nil
}

// Add adds a new transaction to the list,returning whether the
// transaction was accepted, and if yes, any previous transaction it replaced.
//
// PriceBump is the percent
func (list *txList) Add(tx *types.Transaction, priceBump uint64) (bool, *types.Transaction) {
	var old *types.Transaction
	if old = list.Get(tx.Nonce); tx != nil {
		// Replacement strategy. Temporary design
		boundGas := old.GasLimit * (100 + priceBump) / 100
		if boundGas < tx.GasLimit {
			return false, nil
		}
	}
	list.txs.Store(tx.Nonce, tx)
	list.cache = nil
	return true, old
}

// Del deletes a transaction from the list
func (list *txList) Del(nonce uint64) {
	if _, exist := list.txs.Load(nonce); !exist {
		return
	}
	list.txs.Delete(nonce)
	list.cache = nil
}

func (list *txList) Len() int {
	var length int
	list.txs.Range(func(key, value interface{}) bool {
		length++
		return true
	})
	return length
}

// All creates a nonce-sorted slice of current transaction list,
// and the result will be cache in case any modifications are made
func (list *txList) All() types.Transactions {
	if list.cache == nil {
		list.txs.Range(func(key, value interface{}) bool {
			list.cache = append(list.cache, value.(*types.Transaction))
			return true
		})
		sort.Sort(types.NonceSortedList(list.cache))
	}

	results := make(types.Transactions, len(list.cache))
	copy(results, list.cache)

	return results
}

// Filter filters all transactions which make filter func true and false, and
// removes unmatching transactions from the list
func (list *txList) Filter(filter func(tx *types.Transaction) bool) (types.Transactions, types.Transactions) {
	var (
		match   types.Transactions
		unmatch types.Transactions
	)
	list.txs.Range(func(key, value interface{}) bool {
		tx := value.(*types.Transaction)
		if filter(tx) {
			match = append(match, tx)
		} else {
			unmatch = append(unmatch, tx)
		}
		return true
	})

	for _, tx := range unmatch {
		list.txs.Delete(tx.Nonce)
	}

	list.cache = nil
	return match, unmatch
}

// Ready retrieves a sequentially increasing list of transactions starting at the
// provided nonce that is ready for processing. The returned transactions will be
// removed from the list.
func (list *txList) Ready(start uint64) types.Transactions {
	var (
		results types.Transactions
		nonce   = start
	)
	for {
		if tx, exist := list.txs.Load(nonce); exist {
			results = append(results, tx.(*types.Transaction))
		} else {
			break
		}
	}
	list.cache = nil
	return results
}
