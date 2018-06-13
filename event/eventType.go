package event

import (
	"tinychain/core/types"
	"tinychain/p2p"
	"github.com/libp2p/go-libp2p-peer"
)

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
	Tx types.Transaction
}

type NewTxsEvent struct {
	Txs types.Transactions
}

type TxBroadcastEvent struct{}

/*
	Network events
 */

// DiscvActive will be throw when peers discover run 6 rounds
type DiscvActiveEvent struct{}

// Message received from p2p network layer
type SendMsgEvent struct {
	Target peer.ID     // Target peer id
	Typ    string      // Message type
	Data   interface{} // Message data
}

// Protocol manage event
type ProtocolEvent struct {
	Typ      string // 'add' or 'del'
	Protocol *p2p.Protocol
}
