package types

import (
	"math/big"
	"sync/atomic"
	"tinychain/common"
	json "github.com/json-iterator/go"
	"github.com/libp2p/go-libp2p-crypto"
	"errors"
	"tinychain/bmt"
	"tinychain/db/leveldb"
)

type Transaction struct {
	txData

	txHash atomic.Value // hash cache

	Signature atomic.Value `json:"signature"` // Signature of tx
}

type txData struct {
	Nonce   uint64         `json:"nonce"` // Account nonce, which is used to avoid double spending
	Gas     uint64         `json:"gas"`   // Gas used
	Value   *big.Int       `json:"value"` // Transferring value
	From    common.Address `json:"from"`
	To      common.Address `json:"to"` // Recipient of this tx, nil means contract creation
	Payload []byte         `json:"payload"`
}

func NewTransaction(nonce uint64, gas uint64, value *big.Int, payload []byte, from, to common.Address) *Transaction {
	return &Transaction{txData: NewTxData(nonce, gas, value, payload, from, to)}
}

func NewTxData(nonce uint64, gas uint64, value *big.Int, payload []byte, from, to common.Address) txData {
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
	if hash := tx.txHash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	txdata := NewTxData(tx.Nonce, tx.Gas, tx.Value, tx.Payload, tx.From, tx.To)
	data, _ := txdata.Serialize()
	h := common.Sha256(data)
	tx.txHash.Store(h)
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

type Transactions []*Transaction

func (txs Transactions) Hash() common.Hash {
	txSet := bmt.WriteSet{}
	for _, tx := range txs {
		txSet[tx.Hash().String()] = tx.Hash().Bytes()
	}
	root, _ := bmt.Hash(txSet)
	return root
}

func (txs Transactions) Commit(db *leveldb.LDBDatabase) error {
	txSet := bmt.WriteSet{}
	for _, tx := range txs {
		txSet[tx.Hash().String()] = tx.Hash().Bytes()
	}
	return bmt.Commit(txSet, db)
}

// TxMeta represents the meta data of a transaction,
// contains the index of transacitons in a certain block
type TxMeta struct {
	Hash    common.Hash `json:"block_hash"`
	Height  *big.Int    `json:"height"`
	TxIndex uint64      `json:"tx_index"`
}

func (tm *TxMeta) Serialize() ([]byte, error) {
	return json.Marshal(tm)
}

func (tm *TxMeta) Deserialize(d []byte) error {
	return json.Unmarshal(d, tm)
}
