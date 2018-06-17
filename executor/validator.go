package executor

import (
	"tinychain/core"
)

type ValidateImpl struct {
	processor core.Processor // Block processor
}

func NewValidator(processor core.Processor) *ValidateImpl {
	return &ValidateImpl{
		processor: processor,
	}
}
