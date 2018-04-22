package types

import (
	"tinychain/common"
	"math/big"
	"sync/atomic"
	"encoding/json"
	"encoding/binary"
	"encoding/hex"
	"tinychain/bmt"
	"tinychain/core/state"
)

// BNonce is a 64-bit hash which proves that a sufficient amount of
// computation has been carried out on a block
type BNonce [8]byte

func EncodeNonce(i uint64) BNonce {
	var n BNonce
	binary.BigEndian.PutUint64(n[:], i)
	return n
}

func (n BNonce) Uint64() uint64 {
	return binary.BigEndian.Uint64(n[:])
}

func (n BNonce) Hex() []byte {
	return common.Hex(n[:])
}

func (n BNonce) DecodeHex(b []byte) error {
	dec := make([]byte, len(n))
	_, err := hex.Decode(dec, b[2:])
	if err != nil {
		return err
	}
	n.SetBytes(dec)
	return nil
}

func (n BNonce) SetBytes(b []byte) {
	if len(b) > len(n) {
		b = b[:len(n)]
	}
	copy(n[:], b)
}

type Header struct {
	ParentHash common.Hash    `json:"parent_hash"` // Hash of parent block
	Height     *big.Int       `json:"height"`      // Block height
	Difficulty *big.Int       `json:"difficulty"`  // Difficulty of miner
	StateRoot  common.Hash    `json:"stateRoot"`   // State root
	TxRoot     common.Hash    `json:"txRoot"`      // Transaction tree root
	Coinbase   common.Address `json:"miner"`       // Miner address who receives reward of this block
	Extra      []byte         `json:"extra"`       // Extra data
	Nonce      BNonce         `json:"nonce"`       // Nonce produced by pow
	Time       *big.Int       `json:"time"`        // Timestamp
}

func NewHeader(
	parentHash common.Hash,
	height *big.Int,
	difficulty *big.Int,
	stateRoot common.Hash,
	txRoot common.Hash,
	miner common.Address,
	extra []byte,
	nonce BNonce,
	tm *big.Int,
) *Header {
	header := &Header{
		parentHash,
		height,
		difficulty,
		stateRoot,
		txRoot,
		miner,
		extra,
		nonce,
		tm,
	}
	return header
}

func (hd *Header) Hash() common.Hash {
	data, _ := json.Marshal(hd)
	hash := common.Sha256(data)
	return hash
}

func (hd *Header) Serialize() ([]byte, error) { return json.Marshal(hd) }

func (hd *Header) Desrialize(d []byte) error { return json.Unmarshal(d, hd) }

type Block struct {
	Header       *Header        `json:"header"`
	Transactions []*Transaction `json:"transactions"`
	bmt          state.BucketTree             // temporary tx tree
	Hash         atomic.Value   `json:"hash"` // Header hash

	// Total difficulty, to avoid hard fork
	// Tiny will accept the block  with the largest difficulty
	// and link it to the main chain
	TD *big.Int `json:"td"`
}

func NewBlock(header *Header, txs []*Transaction, td *big.Int) *Block {
	block := &Block{
		Header:       header,
		Transactions: txs,
	}
	block.Hash.Store(header.Hash())
	return block
}

func (bl *Block) TxRoot() common.Hash      { return bl.Header.TxRoot }
func (bl *Block) ParentHash() common.Hash  { return bl.Header.ParentHash }
func (bl *Block) Height() *big.Int         { return bl.Header.Height }
func (bl *Block) StateRoot() common.Hash   { return bl.Header.StateRoot }
func (bl *Block) Coinbase() common.Address { return bl.Header.Coinbase }
func (bl *Block) Nonce() BNonce            { return bl.Header.Nonce }
func (bl *Block) Extra() []byte            { return bl.Header.Extra }
func (bl *Block) Difficulty() *big.Int     { return bl.Header.Difficulty }
func (bl *Block) Time() *big.Int           { return bl.Header.Time }

// Calculate hash of block
// Combine header hash and transactions hash, and sha256 it
func (bl *Block) SetHash() common.Hash {
	hash := bl.Header.Hash()[:]
	// Compute transaction hash root
	txSet := bmt.WriteSet{}
	for _, tx := range bl.Transactions {
		txSet[tx.Hash().String()] = tx.Hash().Bytes()
	}

	var tree state.BucketTree
	if tree = bl.bmt; tree == nil {
		tree = new(bmt.BucketTree)
		bl.bmt = tree
	}
	tree.Init(nil)
	tree.Prepare(txSet)
	root, err := tree.Process()
	if err != nil {
		return common.Hash{}
	}
	bl.Header.TxRoot = root

	hash = append(hash, root.Bytes()...)
	return common.Sha256(hash)
}

func (bl *Block) Serialize() ([]byte, error) { return json.Marshal(bl) }

func (bl *Block) Deserialize(d []byte) error { return json.Unmarshal(d, bl) }
