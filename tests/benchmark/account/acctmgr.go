package account

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
)

type AccountManager interface {
	GetAccount(index uint32) *Account
	GetAccountCount() uint32
	GetFaucetAccount() *Account
	Increament() error
}

type accountManagerImpl struct {
	base       uint32
	count      uint32
	increment  uint32
	accounts   []*Account
	faucetAcct *Account
	client     *ethclient.Client
	chainId    *big.Int
}

func NewAccountManager(client *ethclient.Client, baseNumber, initialAccountCnt, incrementPerDay uint32, faucetPrivateKey string) (AccountManager, error) {
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return nil, err
	}

	am := &accountManagerImpl{
		base:      baseNumber,
		count:     initialAccountCnt,
		increment: incrementPerDay,
		accounts:  make([]*Account, 0, initialAccountCnt+incrementPerDay*30),
		client:    client,
		chainId:   chainID,
	}
	am.faucetAcct, err = CreateFaucetAccount(client, faucetPrivateKey, chainID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create faucet account")
	}

	println("Faucet account:", am.faucetAcct.Address.Hex())
	balance, err := am.faucetAcct.GetBalance(client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get faucet account balance")
	}
	println("Faucet account:", am.faucetAcct.Address.Hex(), "balance:", balance.String())
	if balance.Cmp(big.NewInt(10)) < 0 {
		return nil, errors.New("faucet account balance is too low")
	}

	for i := uint32(0); i < initialAccountCnt; i++ {
		key := baseNumber + i
		newAccnt, err := NewAccount(int64(key), client, chainID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create account")
		}
		am.accounts = append(am.accounts, newAccnt)
	}

	return am, nil
}

func (am accountManagerImpl) GetAccount(index uint32) *Account {
	if index < am.count {
		return am.accounts[index]
	}
	return nil
}

func (am accountManagerImpl) GetAccountCount() uint32 {
	return am.count
}

func (am accountManagerImpl) GetFaucetAccount() *Account {
	return am.faucetAcct
}

func (am *accountManagerImpl) Increament() error {
	for i := uint32(0); i < am.increment; i++ {
		key := am.base + am.count
		newAccnt, err := NewAccount(int64(key), am.client, am.chainId)
		if err != nil {
			return errors.Wrap(err, "failed to create account")
		}
		am.accounts = append(am.accounts, newAccnt)
		am.count++
	}
	return nil
}
