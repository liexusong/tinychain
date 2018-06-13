package tiny

import (
	"tinychain/event"
	"tinychain/consensus"
	"tinychain/core"
	"tinychain/db"
	"tinychain/p2p"
	"tinychain/common"
	"tinychain/executor"
	"tinychain/db/leveldb"
	"tinychain/core/state"
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

	state *state.StateDB

	network *Network

	executor executor.Executor
}

func New(config *Config) (*Tinychain, error) {
	eventHub := event.GetEventhub()

	ldb, err := leveldb.NewLDBDataBase("tinychain")
	if err != nil {
		log.Errorf("Cannot create db, %s", err)
		return nil, err
	}
	// Create tiny db
	tinyDB := db.NewTinyDB(ldb)
	// Create state db
	statedb := state.New(ldb, nil)

	peer, err := p2p.New(config.p2p)
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
		state:    statedb,
	}, nil
}

func (self *Tinychain) Init() {

}
