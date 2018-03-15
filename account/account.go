package account

import (
	"tinychain/common"
	"tinychain/core/types"
	"github.com/libp2p/go-libp2p-crypto"
	"sync"
	"errors"
)

var (
	log              = common.GetLogger("account")
	ErrGenAddress    = errors.New("failed to generate address")
	ErrCreateKeyPair = errors.New("failed to create key pair")
	ErrCreateAcc     = errors.New("failed to create account")
	ErrUnlockAcc     = errors.New("failed to unlock account in wallet")
	ErrNotFoundAcc   = errors.New("account is not in wallet. plz unlock it first")
	ErrNotUnlock     = errors.New("account has not been unlocked")
	ErrSignTx        = errors.New("failed to sign transaction")
)

type Account struct {
	Address common.Address
	key     *Key
}

// New account by private key
func NewAccountWithKey(privKey crypto.PrivKey) (*Account, error) {
	key := &Key{privKey, privKey.GetPublic()}
	addr, err := GenAddrByPrivkey(privKey)
	if err != nil {
		return nil, ErrGenAddress
	}
	return &Account{addr, key}, nil
}

func NewAccount() (*Account, error) {
	key, err := NewKeyPairs()
	if err != nil {
		log.Errorf(ErrCreateKeyPair.Error())
		return nil, ErrCreateKeyPair
	}
	addr, err := GenAddrByPubkey(key.pubKey)
	if err != nil {

		return nil, err
	}
	return &Account{addr, key}, nil
}

func (ac *Account) PrivKey() crypto.PrivKey {
	return ac.key.priKey
}

func (ac *Account) PubKey() crypto.PubKey {
	return ac.key.pubKey
}

// Generate address by public key
func GenAddrByPubkey(key crypto.PubKey) (common.Address, error) {
	var addr common.Address
	pubkey, err := key.Bytes()
	if err != nil {
		log.Errorf("Failed to decode pubkey to bytes, %s", err)
		return addr, ErrGenAddress
	}
	pubkey = pubkey[1:]
	h := common.Sha256(pubkey)
	hash := h[len(h)-common.AddressLength:]
	addr = common.HashToAddr(common.Sha256(hash))
	return addr, nil
}

// Generate address by private key
func GenAddrByPrivkey(key crypto.PrivKey) (common.Address, error) {
	pubkey := key.GetPublic()
	return GenAddrByPubkey(pubkey)
}

type Wallet interface {
	Accounts() []*Account
	CreateAccount() (*Account, error)
	Contains(account *Account) bool
	SignHash(account *Account, hash []byte) ([]byte, error)
	SignTx(account *Account, tx *types.Transaction) (*types.Transaction, error)
}

type TinyWallet struct {
	mu       sync.RWMutex
	accounts []*Account
	unlocked map[common.Address]*Key
}

func NewTinyWallet(accounts []*Account) *TinyWallet {
	return &TinyWallet{
		accounts: accounts,
		unlocked: make(map[common.Address]*Key),
	}
}

func (tw *TinyWallet) Accounts() []*Account {
	tw.mu.RLock()
	defer tw.mu.RUnlock()
	var account []*Account
	for _, acc := range tw.accounts {
		account = append(account, acc)
	}
	return account
}

func (tw *TinyWallet) Contains(account *Account) bool {
	for _, acc := range tw.accounts {
		if acc.Address == account.Address {
			return true
		}
	}
	return false
}

func (tw *TinyWallet) CreateAccount() (*Account, error) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	acc, err := NewAccount()
	if err != nil {
		return nil, err
	}
	tw.accounts = append(tw.accounts, acc)
	return acc, nil
}

// Unlock the account, register key in wallet
func (tw *TinyWallet) Unlock(address common.Address, key crypto.PrivKey) error {
	acc, err := NewAccountWithKey(key)
	if err != nil {
		log.Errorf("Failed to generate address %s with key %s", address, key.Bytes())
		return ErrUnlockAcc
	}
	if acc.Address != address {
		log.Errorf("address gen by private key is not equal to the target address")
		return ErrUnlockAcc
	}
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if !tw.Contains(acc) {
		tw.accounts = append(tw.accounts, acc)
	}
	tw.unlocked[address] = acc.key
	return nil
}

func (tw *TinyWallet) SignHash(account *Account, hash []byte) ([]byte, error) {
	if !tw.Contains(account) {
		return nil, ErrNotFoundAcc
	}
	key, ok := tw.unlocked[account.Address]
	if !ok {
		log.Error("account has not been unlocked")
		return nil, ErrNotUnlock
	}
	// TODO Sign hash

}

func (tw *TinyWallet) SignTx(account *Account, tx *types.Transaction) (*types.Transaction, error) {
	if !tw.Contains(account) {
		return nil, ErrNotFoundAcc
	}
	key, ok := tw.unlocked[account.Address]
	if !ok {
		return nil, ErrNotUnlock
	}
	// TODO Sign transaction
}
