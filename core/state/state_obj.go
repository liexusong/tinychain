package state

import (
	"tinychain/common"
	"math/big"
	"encoding/json"
	"tinychain/db/leveldb"
	"tinychain/bmt"
)

type Storage map[common.Hash][]byte

type stateObject struct {
	address common.Address
	data    *Account
	code    []byte     // contract code bytes
	bmt     BucketTree // bucket tree of this account

	dirtyCode bool // code is updated or not
}

type Account struct {
	Nonce    uint64      `json:"nonce"`
	Balance  *big.Int    `json:"balance"`
	Root     common.Hash `json:"root"`
	CodeHash common.Hash `json:"code_hash"`
}

func (s *Account) Serialize() ([]byte, error) {
	return json.Marshal(s)
}

func (s *Account) Deserialize(data []byte) error {
	return json.Unmarshal(data, s)
}

func newStateObject(address common.Address, data *Account, code []byte) *stateObject {
	return &stateObject{
		address: address,
		data:    data,
		code:    code,
	}
}

func (s *stateObject) Address() common.Address {
	return s.address
}

func (s *stateObject) Code() []byte {
	return s.code
}

func (s *stateObject) Balance() *big.Int {
	return s.data.Balance
}

func (s *stateObject) CodeHash() common.Hash {
	return s.data.CodeHash
}

func (s *stateObject) Root() common.Hash {
	return s.data.Root
}

func (s *stateObject) SetCode(code []byte) {
	s.code = code
	s.data.CodeHash = common.Sha256(code)
	s.dirtyCode = true
}

func (s *stateObject) AddBalance(amount *big.Int) {
	s.SetBalance(new(big.Int).Add(s.data.Balance, amount))
}

func (s *stateObject) SubBalance(amount *big.Int) {
	s.SetBalance(new(big.Int).Sub(s.data.Balance, amount))
}

func (s *stateObject) SetBalance(amount *big.Int) {
	s.data.Balance = amount
}

func (s *stateObject) Nonce() uint64 {
	return s.data.Nonce
}

func (s *stateObject) SetNonce(nonce uint64) {
	s.data.Nonce = nonce
}

func (s *stateObject) getBmt(db *leveldb.LDBDatabase) BucketTree {
	tree := bmt.NewBucketTree(db)
	tree.Init(s.data.Root.Bytes())
	return tree
}

func (s *stateObject) Commit() error {
	err := s.bmt.Commit()
	if err != nil {
		return err
	}
	s.data.Root = s.bmt.Hash()
	return nil
}
