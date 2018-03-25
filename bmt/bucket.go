package bmt

import (
	"tinychain/common"
	"encoding/binary"
	"sync"
	"encoding/json"
)

type Bucket struct {
	hash  []byte
	slots map[common.Hash][]byte
}

func NewBucket() *Bucket {
	return &Bucket{
		slots: make(map[common.Hash][]byte),
	}
}

func (bk *Bucket) Hash() common.Hash {
	if bk.hash != nil {
		var hash common.Hash
		return hash.SetBytes(bk.hash)
	}
	var bytes []byte
	for _, v := range bk.slots {
		bytes = append(bytes, v...)
	}
	hash := common.Sha256(bytes)
	bk.hash = hash.Bytes()
	return hash
}

func (bk *Bucket) Serialize() ([]byte, error) {
	return json.Marshal(bk)
}

func (bk *Bucket) Deserialize(d []byte) error {
	return json.Unmarshal(d, bk)
}

type HashTable struct {
	cap     int
	buckets []*Bucket
	dirty   []bool
	lock    sync.RWMutex
}

func NewHashTable(cap int) *HashTable {
	var buckets [cap]*Bucket
	for i := range buckets {
		buckets[i] = NewBucket()
	}
	return &HashTable{
		cap:     cap,
		buckets: buckets[:],
		dirty:   make([]bool, cap, cap),
	}
}

func (ht *HashTable) getIndex(key common.Hash) uint32 {
	val := binary.BigEndian.Uint32(key.Bytes())
	return val % uint32(ht.cap)
}

func (ht *HashTable) add(key common.Hash, value []byte) {
	ht.lock.Lock()
	defer ht.lock.Unlock()
	index := ht.getIndex(key)
	ht.buckets[index].slots[key] = value
	ht.dirty[index] = true
}

func (ht *HashTable) get(key common.Hash) []byte {
	ht.lock.RLock()
	defer ht.lock.RUnlock()
	return ht.buckets[ht.getIndex(key)].slots[key]
}
