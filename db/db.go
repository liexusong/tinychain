package db

import (
	"tinychain/common"
	"math/big"
	"tinychain/core/types"
	"tinychain/db/leveldb"
)

/*
	** Hash is block header hash


	"LastHeader" => the latest block header
	"LastBlock" => the latest block
	"WorldState" => the latest world state root hash

	"h" + height + "n" => hash   block height => block header hash
	"h" + height + hash => header
	"h" + height + hash + "t" => total difficulty
	"H" + hash => height   block hash => height
	"b" + height + hash => block body
	"r" + height + hash => block receipts
	"l" + hash => transaction
*/

const (
	KeyLastHeader = "LastHeader"
	KeyLastBlock  = "LastBlock"
	KeyWorldState = "WorldState"
)

var (
	tinydb *TinyDB
	log    = common.GetLogger("tinydb")
)

type TinyDB struct {
	db *leveldb.LDBDatabase
}

func newTinyDB() (*TinyDB, error) {
	db, err := leveldb.NewLDBDataBase("tinyDatabase")
	if err != nil {
		log.Errorf("Failed to create leveldb, %s", err)
		return nil, err
	}

	return &TinyDB{db}, nil
}

func (tdb *TinyDB) GetWorldState() (common.Hash, error) {
	data, err := tdb.db.Get([]byte(KeyWorldState))
	if err != nil {
		return common.Hash{}, err
	}
	return common.BytesToHash(data), nil
}

func (tdb *TinyDB) GetLastBlock() (*types.Block, error) {
	data, err := tdb.db.Get([]byte(KeyLastBlock))
	if err != nil {
		log.Errorf("Cannot find last block, %s", err)
		return nil, err
	}
	block := &types.Block{}
	err = block.Deserialize(data)
	if err != nil {
		log.Errorf("Failed to decode block, %s", err)
		return nil, err
	}
	return block, nil
}

func (tdb *TinyDB) PutLastBlock(block *types.Block) error {
	data, _ := block.Serialize()
	err := tdb.db.Put([]byte(KeyLastBlock), data)
	if err != nil {
		log.Errorf("Failed to put block, %s", block)
		return err
	}
	return nil
}

func (tdb *TinyDB) GetLastHeader() (*types.Header, error) {
	data, err := tdb.db.Get([]byte(KeyLastHeader))
	if err != nil {
		log.Errorf("Cannot find last header, %s", err)
		return nil, err
	}
	header := &types.Header{}
	err = header.Desrialize(data)
	if err != nil {
		log.Errorf("Failed to decode header, %s", err)
		return nil, err
	}
	return header, nil
}

func (tdb *TinyDB) PutLastHeader(header *types.Header) error {
	data, _ := header.Serialize()
	err := tdb.db.Put([]byte(KeyLastHeader), data)
	if err != nil {
		log.Errorf("Failed to put last header, %s", err)
		return err
	}
	return nil
}

func (tdb *TinyDB) GetHash(height *big.Int) (common.Hash, error) {
	var hash common.Hash
	data, err := tdb.db.Get([]byte("h" + height.String() + "n"))
	if err != nil {
		log.Errorf("Cannot find block header hash with height %s", height)
		return hash, err
	}
	hash = common.DecodeHash(data)
	return hash, nil
}

func (tdb *TinyDB) PutHash(height *big.Int, hash common.Hash) error {
	err := tdb.db.Put([]byte("h"+height.String()+"n"), hash[:])
	if err != nil {
		log.Errorf("Failed to put hash, %s", err)
		return err
	}
	return nil
}

func (tdb *TinyDB) GetHeader(height *big.Int, hash common.Hash) (*types.Header, error) {
	data, err := tdb.db.Get([]byte("h" + height.String() + hash.String()))
	if err != nil {
		log.Errorf("Cannot find header with height %s and hash %s", height, hash)
		return nil, err
	}
	header := types.Header{}
	err = header.Desrialize(data)
	if err != nil {
		log.Error("Failed to decode header")
		return nil, err
	}
	return &header, nil
}

func (tdb *TinyDB) PutHeader(header *types.Header) error {
	data, _ := header.Serialize()
	err := tdb.db.Put([]byte("h"+header.Height.String()+header.Hash().String()), data)
	if err != nil {
		log.Errorf("Failed to put header, %s", err)
		return err
	}
	return nil
}

// Total difficulty
func (tdb *TinyDB) GetTD(height *big.Int, hash common.Hash) (*big.Int, error) {
	data, err := tdb.db.Get([]byte("h" + height.String() + hash.String() + "t"))
	if err != nil {
		log.Errorf("Cannot find total difficulty with height %s and hash %s", height, hash)
		return nil, err
	}
	return new(big.Int).SetBytes(data), nil

}
func (tdb *TinyDB) PutTD(height *big.Int, hash common.Hash, td *big.Int) error {
	err := tdb.db.Put([]byte("h"+height.String()+hash.String()+"t"), td.Bytes())
	if err != nil {
		log.Errorf("Failed to put total difficulty with height %s and hash %s", height, hash)
		return err
	}
	return nil
}

func (tdb *TinyDB) GetHeight(hash common.Hash) (*big.Int, error) {
	data, err := tdb.db.Get([]byte("H" + hash.String()))
	if err != nil {
		log.Errorf("Cannot find height with hash %s", hash)
		return nil, err
	}
	return new(big.Int).SetBytes(data), nil
}

func (tdb *TinyDB) PutHeight(hash common.Hash, height *big.Int) error {
	err := tdb.db.Put([]byte("H"+hash.String()), height.Bytes())
	if err != nil {
		log.Errorf("Failed to put height with hash %s", hash)
		return err
	}
	return nil
}

func (tdb *TinyDB) GetBlock(height *big.Int, hash common.Hash) (*types.Block, error) {
	data, err := tdb.db.Get([]byte("b" + height.String() + hash.String()))
	if err != nil {
		log.Errorf("Cannot find block with height %s and hash %s", height, hash)
		return nil, err
	}
	block := types.Block{}
	err = block.Deserialize(data)
	if err != nil {
		log.Errorf("Failed to decode block with height %s and hash %s", height, hash)
		return nil, err
	}
	return &block, nil
}

func (tdb *TinyDB) PutBlock(block *types.Block) error {
	height := block.Header.Height
	hash := block.Header.Hash()
	data, _ := block.Serialize()
	err := tdb.db.Put([]byte("b"+height.String()+hash.String()), data)
	if err != nil {
		log.Errorf("Failed to put block with height %s", height)
		return err
	}
	return nil
}

func (tdb *TinyDB) GerReceipt(height *big.Int, hash common.Hash) *types.Receipt {}

func (tdb *TinyDB) PutReceipt(height *big.Int, hash common.Hash, receipt *types.Receipt) error {}

func GetTinyDB() *TinyDB {
	if tinydb != nil {
		log.Info("TinyDB has not created, create one")
		return tinydb
	}

	db, err := newTinyDB()
	if err != nil {
		return nil
	}
	tinydb = db
	return tinydb
}
