package dpos

type DposEngine struct {
}

func NewDpos() *DposEngine {
	return &DposEngine{}
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
