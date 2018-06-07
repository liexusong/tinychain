package consensus

type Consensus interface {
	Name() string
	Start() error
	Stop() error
}

