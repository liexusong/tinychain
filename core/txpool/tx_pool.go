package txpool

import (
	"tinychain/common"
	"tinychain/executor"
	"tinychain/core/types"
	"tinychain/event"
)

type TxPool struct {
	txValidator executor.TxValidator
	pending     map[common.Address]*txList // all valid and processable txs
	queue       map[common.Address]*txList // all new-added and non-processable txs,including valid and invalid txs.
}

func (tp *TxPool) Pending() (map[common.Address]types.Transactions, error) {

}

func (tp *TxPool) AddLocal(txs types.Transactions) error {

}

func (tp *TxPool) AddRemotes(txs types.Transactions) error {

}

func (tp *TxPool) SubscribeNewTxsEvent() event.Subscription {

}
