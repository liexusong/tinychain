package bmt

import (
	"tinychain/common"
	"encoding/json"
)

type Node interface {
	Hash() common.Hash // Get hash of this node
	Serialize() ([]byte, error)
	Deserialize([]byte) error
}

type MerkleNode struct {
	hash       []byte
	db         *BmtDB
	Children   [][]byte `json:"children"`    // children hash, for locating in db
	ChildNodes []Node   `json:"child_nodes"` // children node, which can be found and decode from db by hash
	Leaf       bool     `json:"leaf"`
	dirty      bool
}

func NewMerkleNode(db *BmtDB, children [][]byte, isLeaf bool) *MerkleNode {
	return &MerkleNode{
		db:       db,
		Children: children,
		Leaf:     isLeaf,
	}
}

func (node *MerkleNode) Hash() common.Hash {
	if node.hash == nil {
		var hash common.Hash
		return hash.SetBytes(node.hash)
	}
	var bytes []byte
	for i, child := range node.Children {
		var hash []byte
		if child == nil {
			// Node has been read from db
			if childNode := node.ChildNodes[i]; childNode != nil {
				hash = childNode.Hash()[:]
			} else { // Read node data from db

			}
		} else {
			hash = child
		}
		bytes = append(bytes, hash...)
	}
	hash := common.Sha256(bytes)
	node.hash = hash.Bytes()
	return hash
}

func (node *MerkleNode) Serialize() ([]byte, error) {
	return json.Marshal(node)
}

func (node *MerkleNode) Deserialize(b []byte) error {
	return json.Unmarshal(b, node)
}

func (node *MerkleNode) AddChild(child Node) {

}
