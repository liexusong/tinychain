package core

import (
	"tinychain/core/types"
	"math/big"
)

/*
	Blockchain events
 */

type AppendBlockEvent struct {
	Blocks []*types.Block
}

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

type ExecBlockEvent struct {
	Block *types.Block
}

/*
	Transaction events
 */
type NewTxEvent struct {
	Tx *types.Transaction
}

type ExecPendingTxEvent struct {
	Txs types.Transactions
}

type TxBroadcastEvent struct{}
