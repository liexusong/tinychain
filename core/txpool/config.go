package txpool

type Config struct {
	MaxTxSize    uint64 // Max size of tx pool
	MaxBloomSize uint	// Max bloom size
	PriceBump    int	// Price bump to decide whether to replace tx or not
}
