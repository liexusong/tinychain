package account

import (
	"tinychain/common"
	"tinychain/core/types"
)

type Account struct {
	Address common.Address
}

type Wallet interface {
	Accounts() []*Account
	Contains(account *Account) bool
	SignHash(account *Account, hash []byte) ([]byte, error)
	SignTx(account *Account, tx *types.Transaction) (*types.Transaction, error)
}
