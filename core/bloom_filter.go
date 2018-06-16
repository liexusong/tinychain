package core

import "github.com/willf/bloom"

// BloomFilter is used to accelerate the process of checking
// a bytes is existed or not
type BloomFilter struct {
	filter bloom.BloomFilter
}
