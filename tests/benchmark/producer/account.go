package producer

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"math/big"
	"strings"

	"github.com/0glabs/0g-chain/tests/benchmark/utils"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
)

type Account struct {
	Nonce      uint64
	ChainId    *big.Int
	Address    common.Address
	PrivateKey *ecdsa.PrivateKey
	IsFaucet   bool

	Signer  *Signer
	ReqChan chan<- *TxSignRequest
	ResChan <-chan *TxSignResponse
}

func NewAccount(index int64, client *ethclient.Client, chainId *big.Int) (*Account, error) {
	privKey := hd.Secp256k1.Generate()(new(big.Int).SetInt64(index + 1).Bytes())
	ecdsaPrivKey, err := crypto.HexToECDSA(hex.EncodeToString(privKey.Bytes()))
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert private key to ecdsa")
	}

	addr := crypto.PubkeyToAddress(ecdsaPrivKey.PublicKey)
	nonce, err := client.PendingNonceAt(context.Background(), addr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get pending nonce")
	}

	signer, err := NewSigner(client, ecdsaPrivKey, chainId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create signer")
	}

	reqChan := make(chan *TxSignRequest)
	resChan := signer.Run(reqChan)

	return &Account{
		Nonce:      nonce,
		Address:    addr,
		PrivateKey: ecdsaPrivKey,
		ChainId:    chainId,
		Signer:     signer,
		ReqChan:    reqChan,
		ResChan:    resChan,
	}, nil
}

func CreateFaucetAccount(client *ethclient.Client, privateKey string, chainId *big.Int) (*Account, error) {
	pks := strings.TrimPrefix(privateKey, "0x")
	pkBytes, _ := hex.DecodeString(pks)
	pk := utils.LoadPrivateKey(pkBytes)
	addr := crypto.PubkeyToAddress(pk.PublicKey)

	nonce, err := client.PendingNonceAt(context.Background(), addr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get pending nonce")
	}

	signer, err := NewSigner(client, pk, chainId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create signer")
	}

	reqChan := make(chan *TxSignRequest)
	resChan := signer.Run(reqChan)

	return &Account{
		IsFaucet:   true,
		Nonce:      nonce,
		Address:    addr,
		PrivateKey: pk,
		ChainId:    chainId,
		Signer:     signer,
		ReqChan:    reqChan,
		ResChan:    resChan,
	}, nil
}

func (a *Account) GetAndIncrementNonce() uint64 {
	now := a.Nonce
	a.Nonce += 1
	return now
}

func (a *Account) GetBalance(client *ethclient.Client) (*big.Int, error) {
	b, err := client.BalanceAt(context.Background(), a.Address, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get balance")
	}
	return b, nil
}
