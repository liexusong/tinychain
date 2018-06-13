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
		boundGas := old.Gas * (100 + priceBump) / 100
		if boundGas < tx.Gas {
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

// Filter removes all transactions from the list with a cost or gas limit
// higher than the provided thresholds.
// Every removed transaction is returned as the second return value.
func (list *txList) Filter() (types.Transactions, types.Transactions) {

}
