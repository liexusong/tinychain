// Bucket merkle tree implementaion
package bmt

import (
	"tinychain/common"
	"tinychain/db/leveldb"
)

var (
	log = common.GetLogger("bucket_tree")
)

const (
	defaultHashTableCap = 10000
	defaultAggreation   = 10
)

type BucketTree struct {
	capacity   int
	aggreation int
	db         *BmtDB
	root       Node
	dirty      bool
}

func NewBucketTree(db *leveldb.LDBDatabase) *BucketTree {
	// v1.0
	return &BucketTree{
		capacity:   defaultHashTableCap,
		aggreation: defaultAggreation,
		db:         NewBmtDB(db),
	}
}

func (bt *BucketTree) Hash() common.Hash {
	return bt.root.Hash()
}

func (bt *BucketTree) Prepare() {

}

func (bt *BucketTree) Process() (common.Hash, error) {

}

func (bt *BucketTree) Commit() {

}
