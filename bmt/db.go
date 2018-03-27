package bmt

import (
	"tinychain/db/leveldb"
	"tinychain/common"
)

const (
	NodeKeyPrefix      = "n" // "n" + node hash
	HashTableKeyPrefix = "t" // "t" + root node hash
	BucketKeyPrefix    = "s" // "s" + slot hash
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
	node := &MerkleNode{db: bdb}
	node.deserialize(data)
	return node, nil
}

func (bdb *BmtDB) PutNode(key common.Hash, node *MerkleNode) error {
	data, err := node.serialize()
	if err != nil {
		return err
	}
	err = bdb.db.Put([]byte(NodeKeyPrefix+key.String()), data)
	if err != nil {
		log.Errorf("Failed to put node to BmtDB: %s", err)
		return err
	}
	return nil
}

func (bdb *BmtDB) GetBucket(key common.Hash) (*Bucket, error) {
	data, err := bdb.db.Get([]byte(BucketKeyPrefix + key.String()))
	if err != nil {
		return nil, err
	}
	bucket := &Bucket{}
	err = bucket.deserialize(data)
	if err != nil {
		return nil, err
	}
	return bucket, nil
}

func (bdb *BmtDB) PutBucket(key common.Hash, bucket *Bucket) error {
	data, err := bucket.serialize()
	if err != nil {
		return nil
	}
	err = bdb.db.Put([]byte(BucketKeyPrefix+key.String()), data)
	if err != nil {
		log.Errorf("Failed to put bucket to BmtDB: %s", err)
		return err
	}
	return nil
}

func (bdb *BmtDB) GetHashTable(key common.Hash) (*HashTable, error) {
	data, err := bdb.db.Get([]byte(HashTableKeyPrefix + key.String()))
	if err != nil {
		return nil, err
	}
	ht := &HashTable{}
	err = ht.deserialize(data)
	if err != nil {
		return nil, err
	}
	return ht, nil
}

func (bdb *BmtDB) PutHashTable(key common.Hash, ht *HashTable) error {
	data, err := ht.serialize()
	if err != nil {
		return err
	}
	return bdb.db.Put([]byte(HashTableKeyPrefix+key.String()), data)
}
