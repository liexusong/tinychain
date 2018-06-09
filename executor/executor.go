package executor

import (
	"tinychain/core"
	"tinychain/event"
)

type Executor interface {
	Start() error
	Stop() error
}

type ExecutorImpl struct {
	validator Validator        // Validator validate all consensus fields
	chain     *core.Blockchain // Blockchain wrapper
	event     *event.TypeMux
}

func New() Executor{

}

func (ex *ExecutorImpl) Start() {

}
