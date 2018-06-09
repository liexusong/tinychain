package tiny

import (
	"tinychain/event"
	"tinychain/consensus"
	"tinychain/core"
	"tinychain/db"
	"tinychain/p2p"
	"tinychain/common"
)

var (
	log = common.GetLogger("tinychain")
)

// Tinychain implements the tinychain full node service
type Tinychain struct {
	config *Config

	eventHub *event.TypeMux

	engine consensus.Engine

	chain *core.Blockchain

	db *db.TinyDB

	network *p2p.Peer
}

func New(config *Config) (*Tinychain, error) {
	eventHub := event.NewTypeMux()

	tinyDB, err := db.NewTinyDB()
	if err != nil {
		log.Error("Failed to create leveldb")
		return nil, err
	}
	peer, err := p2p.New(&config.p2p, eventHub)
	if err != nil {
		log.Error("Failed to create p2p peer")
		return nil, err
	}

	engine := consensus.New()

	bc, err := core.NewBlockchain(tinyDB, engine)
	if err != nil {
		log.Error("Failed to create blockchain")
		return nil, err
	}

	return &Tinychain{
		config:   config,
		eventHub: eventHub,
		db:       tinyDB,
		network:  peer,
		chain:    bc,
		engine:   engine,
	}, nil
}

func (self *Tinychain) Init() {

}
