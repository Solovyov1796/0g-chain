package producer

import (
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
)

type AccountMap struct {
	total      uint32
	accounts   []*Account
	faucetAcct *Account
}

func NewAccountMap(client *ethclient.Client, base, total uint32, faucetPrivateKey string, chainId *big.Int) (*AccountMap, error) {
	am := &AccountMap{
		total:    total,
		accounts: make([]*Account, 0, total),
	}
	var err error
	am.faucetAcct, err = CreateFaucetAccount(client, faucetPrivateKey, chainId)
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

	for i := uint32(0); i < total; i++ {
		key := base + i
		newAccnt, err := NewAccount(int64(key), client, chainId)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create account")
		}
		am.accounts = append(am.accounts, newAccnt)
		println("[", key, "] Created account:", newAccnt.Address.Hex())
	}

	return am, nil
}

func (am AccountMap) GetAccount(index uint32) *Account {
	if index < am.total {
		return am.accounts[index]
	}
	return nil
}

func (am AccountMap) GetAccountCount() uint32 {
	return am.total
}

func (am AccountMap) GetFaucetAccount() *Account {
	return am.faucetAcct
}
