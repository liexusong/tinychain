package consensus

import "tinychain/consensus/dpos"

type Engine interface {
	Name() string
	Start() error
	Stop() error
}

func New() Engine {
	return dpos.NewDpos()
}
