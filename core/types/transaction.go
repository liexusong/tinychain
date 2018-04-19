package types

import (
	"math/big"
	"sync/atomic"
	"tinychain/common"
	"encoding/json"
	"github.com/libp2p/go-libp2p-crypto"
	"errors"
)

type Transaction struct {
	txData

	TxHash    atomic.Value `json:"hash"`
	Signature atomic.Value `json:"signature"` // Signature of tx
}

type txData struct {
	Nonce   uint64         `json:"nonce"` // Account nonce, which is used to avoid double spending
	Gas     *big.Int       `json:"gas"`   // Gas used
	Value   *big.Int       `json:"value"` // Transferring value
	From    common.Address `json:"from"`
	To      common.Address `json:"to"` // Recipient of this tx, nil means contract creation
	Payload []byte         `json:"payload"`
}

func NewTransaction(nonce uint64, gas, value *big.Int, payload []byte, from, to common.Address) *Transaction {
	return &Transaction{txData: NewTxData(nonce, gas, value, payload, from, to)}
}

func NewTxData(nonce uint64, gas, value *big.Int, payload []byte, from, to common.Address) txData {
	return txData{
		Nonce:   nonce,
		Gas:     gas,
		Value:   value,
		Payload: payload,
		From:    from,
		To:      to,
	}
}

func (txd txData) Serialize() ([]byte, error) { return json.Marshal(txd) }
func (txd txData) Deserialize(d []byte) error { return json.Unmarshal(d, txd) }

func (tx *Transaction) Serialize() ([]byte, error) { return json.Marshal(tx) }
func (tx *Transaction) Deserialize(d []byte) error { return json.Unmarshal(d, tx) }

func (tx *Transaction) Hash() common.Hash {
	if hash := tx.TxHash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	txdata := NewTxData(tx.Nonce, tx.Gas, tx.Value, tx.Payload, tx.From, tx.To)
	data, _ := txdata.Serialize()
	h := common.Sha256(data)
	tx.TxHash.Store(h)
	return h
}

// Sign the transaction with private key
func (tx *Transaction) Sign(privKey crypto.PrivKey) ([]byte, error) {
	if sign := tx.Signature.Load(); sign != nil {
		return sign.([]byte), nil
	}
	hash := tx.Hash()
	s, err := privKey.Sign(hash[:])
	if err != nil {
		return nil, err
	}
	tx.Signature.Store(s)
	return s, nil
}

// Verify transaction signature by specific public key
func (tx *Transaction) Verify(pubKey crypto.PubKey) (bool, error) {
	sign := tx.Signature.Load()
	if sign != nil {
		return false, errors.New("signature not found")
	}
	hash := tx.Hash()
	equal, err := pubKey.Verify(hash[:], sign.([]byte))
	if err != nil {
		return false, errors.New("error occur during sign verification")
	}
	return equal, nil
}

func (tx *Transaction) AsEvent() Event {
	return Event{
		from:  tx.From,
		to:    tx.To,
		nonce: tx.Nonce,
		value: tx.Value,
		data:  tx.Payload,
	}
}

// Event is a derived transaction and implements core.Event
type Event struct {
	from  common.Address
	to    common.Address
	nonce uint64
	value *big.Int
	data  []byte
}

func (ev Event) From() common.Address { return ev.from }
func (ev Event) To() common.Address   { return ev.to }
func (ev Event) Nonce() uint64        { return ev.nonce }
func (ev Event) Value() *big.Int      { return ev.value }
func (ev Event) Data() []byte         { return ev.data }
