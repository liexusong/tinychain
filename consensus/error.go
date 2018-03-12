package consensus

import "errors"

const (

	// A block's height doesn't equal to its parent's height plus one
	ErrInvalidHeight = errors.New("invalid block height")
)
