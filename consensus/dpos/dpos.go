package dpos

import "github.com/nebulasio/go-nebulas/consensus/dpos"

type DposEngine struct {
}

func NewDpos() *DposEngine{
	return &dpos.Dpos{}
}

func (dpos *DposEngine) Name() string {
	return "TinyDPoS"
}

func (dpos *DposEngine) Start() error {
	return nil
}

func (dpos *DposEngine) Stop() error {
	return nil
}
