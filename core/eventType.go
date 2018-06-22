package core

import (
	"tinychain/core/types"
	"math/big"
)

/*
	Block events
 */
type NewBlockEvent struct {
	Block *types.Block
}

type BlockBroadcastEvent struct{}

type BlockCommitEvent struct {
	Height *big.Int
}

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
