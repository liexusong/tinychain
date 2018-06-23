package txpool

type Config struct {
	MaxTxSize    uint64 // Max size of tx pool
	PriceBump    int	// Price bump to decide whether to replace tx or not
}
