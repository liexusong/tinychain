package txpool

import (
	"tinychain/core/types"
	"sync"
	"math/big"
)

type TxList struct {
	txs   sync.Map
	cache types.Transactions
}

func (list *TxList) Transactions() {

}

func (list *TxList) Get(nonce uint64) *types.Transaction {
	if tx, exist := list.txs.Load(nonce); exist {
		return tx.(*types.Transaction)
	}
	return nil
}

// Add adds a new transaction to the list,returning whether the
// transaction was accepted, and if yes, any previous transaction it replaced.
//
// PriceBump is the percent
func (list *TxList) Add(tx *types.Transaction, priceBump uint64) (bool, *types.Transaction) {
	var old *types.Transaction
	if old = list.Get(tx.Nonce); tx != nil {
		ratioGas := new(big.Int).Div(new(big.Int).Mul(old.Gas, new(big.Int).SetUint64(priceBump)), new(big.Int).SetUint64(uint64(100)))
		boundGas := new(big.Int).Add(old.Gas, ratioGas)
		if tx.Gas.Cmp(boundGas) < 0 {
			return false, nil
		}
	}
	list.txs.Store(tx.Nonce, tx)
	list.cache = nil
	return true, old
}

// Del deletes a transaction from the list
func (list *TxList) Del(nonce uint64) {
	if _, exist := list.txs.Load(nonce); !exist {
		return
	}
	list.txs.Delete(nonce)
	list.cache = nil
}

func (list *TxList) Len() int {
	var length int
	list.txs.Range(func(key, value interface{}) bool {
		length++
		return true
	})
	return length
}

// All creates a nonce-sorted slice of current transaction list,
// and the result will be cache in case any modifications are made
func (list *TxList) All() types.Transactions {
	if list.cache != nil {
		return list.cache
	}

	list.txs.Range(func(key, value interface{}) bool {
		list.cache = append(list.cache, value.(*types.Transaction))
		return true
	})

	results := make(types.Transactions, len(list.cache))
	copy(results, list.cache)
	return results
}
