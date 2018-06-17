package core

import "tinychain/core/types"

/*
	Block events
 */
type NewBlockEvent struct {
	Block *types.Block
}

type BlockBroadcastEvent struct{}

/*
	Transaction events
 */
type NewTxEvent struct {
	Tx *types.Transaction
}

type NewTxsEvent struct {
	Txs types.Transactions
}

type TxBroadcastEvent struct{}
