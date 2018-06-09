package types

import (
	"tinychain/common"
	"math/big"
	"sync/atomic"
	json "github.com/json-iterator/go"
	"encoding/binary"
	"encoding/hex"
	"tinychain/db"
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
	ParentHash   common.Hash    `json:"parent_hash"`  // Hash of parent block
	Height       *big.Int       `json:"height"`       // Block height
	StateRoot    common.Hash    `json:"state_root"`   // State root
	TxRoot       common.Hash    `json:"tx_root"`      // Transaction tree root
	ReceiptsHash common.Hash    `json:"receipt_hash"` // Receipts hash
	Coinbase     common.Address `json:"miner"`        // Miner address who receives reward of this block
	Extra        []byte         `json:"extra"`        // Extra data
	Time         *big.Int       `json:"time"`         // Timestamp
	GasUsed      *big.Int       `json:"gas"`          // Total gas used
	GasLimit     *big.Int       `json:"gas_limit"`    // Gas limit of this block
}

func (hd *Header) Hash() common.Hash {
	data, _ := json.Marshal(hd)
	hash := common.Sha256(data)
	return hash
}

func (hd *Header) Serialize() ([]byte, error) { return json.Marshal(hd) }

func (hd *Header) Desrialize(d []byte) error { return json.Unmarshal(d, hd) }

type Block struct {
	Header       *Header      `json:"header"`
	Transactions Transactions `json:"transactions"`
	Receipts     Receipts     `json:"receipts"`
	hash         atomic.Value // Header hash cache
	size         atomic.Value // Block size cache
}

func NewBlock(header *Header, txs Transactions) *Block {
	block := &Block{
		Header:       header,
		Transactions: txs,
	}
	return block
}

func (bl *Block) TxRoot() common.Hash       { return bl.Header.TxRoot }
func (bl *Block) ReceiptsHash() common.Hash { return bl.Header.ReceiptsHash }
func (bl *Block) ParentHash() common.Hash   { return bl.Header.ParentHash }
func (bl *Block) Height() *big.Int          { return bl.Header.Height }
func (bl *Block) StateRoot() common.Hash    { return bl.Header.StateRoot }
func (bl *Block) Coinbase() common.Address  { return bl.Header.Coinbase }
func (bl *Block) Extra() []byte             { return bl.Header.Extra }
func (bl *Block) Time() *big.Int            { return bl.Header.Time }

// Calculate hash of block
// Combine header hash and transactions hash, and sha256 it
func (bl *Block) Hash() common.Hash {
	if hash := bl.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	// Compute transaction tree hash root
	//txSet := bmt.WriteSet{}
	//for _, tx := range bl.Transactions {
	//	txSet[tx.Hash().String()] = tx.Hash().Bytes()
	//}
	root := bl.Transactions.Hash()

	bl.Header.TxRoot = root

	// Compute receipts hash root
	bl.Header.ReceiptsHash = bl.Receipts.Hash()

	hash := bl.Header.Hash()
	bl.hash.Store(hash)
	return hash
}

func (bl *Block) Size() uint64 {
	if size := bl.size.Load(); size != nil {
		return size.(uint64)
	}
	tmp, _ := bl.Serialize()
	bl.size.Store(len(tmp))
	return uint64(len(tmp))
}

// Commit stores the block to db
func (bl *Block) Commit(db *db.TinyDB) error {
	if hash := bl.hash.Load(); hash == nil {
		bl.Hash()
	}

	// Commit transactions tree
	err := bl.Transactions.Commit(db.LDB())
	if err != nil {
		return err
	}

	// Commit header
	err = db.PutHash(bl.Height(), bl.Hash())
	if err != nil {
		return err
	}
	err = db.PutHeader(bl.Header)
	if err != nil {
		return err
	}
	err = db.PutHeight(bl.Hash(), bl.Height())

	// Commit block
	return db.PutBlock(bl)
}

func (bl *Block) Serialize() ([]byte, error) { return json.Marshal(bl) }

func (bl *Block) Deserialize(d []byte) error { return json.Unmarshal(d, bl) }
