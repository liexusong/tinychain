package bmt

import (
	"tinychain/db/leveldb"
	"tinychain/common"
)

const (
	NodeKeyPrefix = "n" // "n" + node hash
	SlotKeyPrefix = "s" // "s" + slot hash
)

type BmtDB struct {
	db *leveldb.LDBDatabase
}

func NewBmtDB(db *leveldb.LDBDatabase) *BmtDB {
	return &BmtDB{
		db: db,
	}
}

func (bdb *BmtDB) GetNode(key common.Hash) (*MerkleNode, error) {
	data, err := bdb.db.Get([]byte(NodeKeyPrefix + key.String()))
	if err != nil {
		return nil, err
	}
	node := &MerkleNode{}
	node.Deserialize(data)
	return node, nil
}

func (bdb *BmtDB) PutNode(key common.Hash, node *MerkleNode) error {
	data, err := node.Serialize()
	if err != nil {
		return err
	}
	err = bdb.db.Put([]byte(NodeKeyPrefix+key.String()), data)
	if err != nil {
		log.Errorf("Failed to put node to BmtDB, %s", err)
		return err
	}
	return nil
}

func (bdb *BmtDB) Get