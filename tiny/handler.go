package tiny

import (
	"tinychain/common"
	"tinychain/core/types"
	"tinychain/event"
)

type TxPool interface {
	// AddRemotes adds remote transactions to pending tx list
	AddRemotes(txs types.Transactions) error

	// Pending returns all valid and processable transactions
	Pending() (map[common.Address]types.Transactions, error)

	SubscribeNewTxsEvent() event.Subscription
}

type HandlerManager struct {
	txPool TxPool
}
