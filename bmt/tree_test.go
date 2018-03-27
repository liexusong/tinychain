package bmt

import (
	"testing"
	"tinychain/db/leveldb"
	"github.com/stretchr/testify/assert"
	"fmt"
	"tinychain/common"
)

var (
	btree = CreateBucketTree()
)

func CreateBucketTree() *BucketTree {
	db, _ := leveldb.NewLDBDataBase(nil, "bucket_tree_test")
	return NewBucketTree(db)
}

func TestBucketTree_Process(t *testing.T) {
	writeSet := NewWriteSet()
	writeSet["test1"] = []byte("asdffsdf")
	writeSet["abcd"] = []byte("test2asd")
	writeSet["lslsl"] = []byte("test3f")
	writeSet["werw"] = []byte("test12as")
	writeSet["ffff"] = []byte("FDas")
	writeSet["asdf"] = []byte("asdfff")
	err := btree.Init(nil)
	assert.Nil(t, err)
	err = btree.Prepare(writeSet)
	assert.Nil(t, err)
	err = btree.Commit()
	assert.Nil(t, err)
}

func TestBucketTree_Read(t *testing.T) {
	//for _, bucket := range btree.hashTable.buckets {
	//	if bucket != nil && !bucket.Hash().Nil() {
	//fmt.Printf("%s\n", bucket.Hash().Hex())
	//	}
	//}
	fmt.Printf("Tree root is %s\n", btree.Hash().Hex())
	var nilHash common.Hash
	assert.NotEqual(t, nilHash, btree.Hash())
}

func TestBucketTree_Update(t *testing.T) {
	oldRoot := btree.Hash()
	newSet := NewWriteSet()
	newSet["lowesyang"] = []byte("lowesyang")
	err := btree.Prepare(newSet)
	assert.Nil(t, err)
	err = btree.Commit()
	assert.Nil(t, err)
	fmt.Printf("%s\n", btree.Hash().Hex())
	assert.NotEqual(t, oldRoot, btree.Hash())
}
