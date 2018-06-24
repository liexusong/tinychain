package txpool

import "time"

type Config struct {
	MaxTxSize     uint64 // Max size of tx pool
	PriceBump     int    // Price bump to decide whether to replace tx or not
	BatchTimeout  time.Duration
	BatchCapacity int
}
