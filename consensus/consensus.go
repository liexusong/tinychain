package consensus

import (
	"tinychain/consensus/dpos"
	"tinychain/core/types"
	"tinychain/common"
)

type Engine interface {
	Name() string
	Start() error
	Stop() error
}

type TxPool interface {
	// AddRemotes adds remote transactions to queue tx list
	AddRemotes(txs types.Transactions) error

	// Pending returns all valid and processable transactions
	Pending() map[common.Address]types.Transactions
}

func New() Engine {
	return dpos.NewDpos()
}
