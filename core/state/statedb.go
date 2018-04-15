package state

import (
	"tinychain/common"
	"tinychain/bmt"
	"tinychain/db/leveldb"
	"math/big"
)

var (
	log = common.GetLogger("state")
)

// Bucket tree
type BucketTree interface {
	Hash() common.Hash
	Init(root []byte) error
	Prepare(dirty bmt.WriteSet) error
	Process() (common.Hash, error)
	Commit() error
	Get(key []byte) ([]byte, error)
	Copy() *bmt.BucketTree
}

type StateDB struct {
	db  *cacheDB
	bmt BucketTree

	stateObjects      map[common.Address]*stateObject
	stateObjectsDirty map[common.Address]struct{}
}

func New(db *leveldb.LDBDatabase, root common.Hash) *StateDB {
	tree := bmt.NewBucketTree(db)
	if err := tree.Init(root.Bytes()); err != nil {
		log.Errorf("Failed to init bucket tree when new state db, %s", err)
		return nil
	}
	return &StateDB{
		db:                newCacheDB(db),
		bmt:               tree,
		stateObjects:      make(map[common.Address]*stateObject),
		stateObjectsDirty: make(map[common.Address]struct{}),
	}
}

// Get state object from cache and bucket tree
// If not found, return nil
func (sdb *StateDB) GetStateObj(addr common.Address) (*stateObject, error) {
	if stateObj, exist := sdb.stateObjects[addr]; exist {
		return stateObj, nil
	}

	data, err := sdb.bmt.Get(addr.Bytes())
	if err != nil {
		log.Errorf("State object not found, %s", err)
		return nil, err
	}
	account := &Account{}
	err = account.Deserialize(data)
	if err != nil {
		return nil, err
	}
	code, _ := sdb.db.GetCode(account.CodeHash)
	stateObj := newStateObject(addr, account, code)
	sdb.setStateObj(stateObj)
	return stateObj, nil
}

// Create a new state object
func (sdb *StateDB) CreateStateObj(addr common.Address) *stateObject {
	account := &Account{
		Nonce:   uint64(0),
		Balance: new(big.Int),
	}
	newObj := newStateObject(addr, account, nil)
	sdb.setStateObj(newObj)
	return newObj
}

// Set "live" state object
func (sdb *StateDB) setStateObj(object *stateObject) {
	sdb.stateObjects[object.Address()] = object
	sdb.stateObjectsDirty[object.Address()] = struct{}{}
}

// Process dirty state object to state tree and get intermediate root
func (sdb *StateDB) IntermediateRoot() (common.Hash, error) {
	dirtySet := bmt.NewWriteSet()
	for addr := range sdb.stateObjectsDirty {
		stateobj := sdb.stateObjects[addr]
		data, _ := stateobj.data.Serialize()
		dirtySet[string(addr.Bytes())] = data
	}
	if err := sdb.bmt.Prepare(dirtySet); err != nil {
		return common.Hash{}, err
	}
	return sdb.bmt.Process()
}

func (sdb *StateDB) Commit() error {
	dirtySet := bmt.NewWriteSet()
	codeSet := make(map[common.Hash]*stateObject)

	for addr := range sdb.stateObjectsDirty {
		delete(sdb.stateObjectsDirty, addr)
		stateobj := sdb.stateObjects[addr]
		// Put account data to dirtySet
		data, _ := stateobj.data.Serialize()
		dirtySet[string(addr.Bytes())] = data

		// Put code bytes to codeSet
		if stateobj.dirtyCode {
			// Put stateobj as value, because we need to set `dirtyCode` to false
			// after write in batch
			codeSet[stateobj.CodeHash()] = stateobj
		}
	}

	if err := sdb.bmt.Prepare(dirtySet); err != nil {
		return err
	}
	if err := sdb.bmt.Commit(); err != nil {
		return err
	}
	if err := sdb.db.PutCodeInBatch(codeSet); err != nil {
		return err
	}
	return nil
}
