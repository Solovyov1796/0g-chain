package utils

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/sha3"
	"golang.org/x/exp/rand"
)

func ToChecksumAddress(address string) string {
	address = strings.ToLower(strings.TrimPrefix(address, "0x"))

	hasher := sha3.NewLegacyKeccak256()
	hasher.Write([]byte(address))
	hash := hasher.Sum(nil)
	hashHex := hex.EncodeToString(hash)

	checksummedAddress := "0x"
	for i, c := range address {
		if c >= '0' && c <= '9' {
			checksummedAddress += string(c)
		} else {
			if hashHex[i] >= '8' {
				checksummedAddress += strings.ToUpper(string(c))
			} else {
				checksummedAddress += string(c)
			}
		}
	}

	return checksummedAddress
}

func LoadPrivateKey(privateKeyBytes []byte) *ecdsa.PrivateKey {
	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		panic(err)
	}
	return privateKey
}

func Shuffle(slice []uint32) []uint32 {
	rand.Seed(uint64(time.Now().UnixNano()))
	n := len(slice)
	for i := n - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}

func ToBigInt(amount any) *big.Int {
	if amount == nil {
		return big.NewInt(0)
	}
	var val *big.Int
	switch amount.(type) {
	case int:
		val = big.NewInt(int64(amount.(int)))
	case int32:
		val = big.NewInt(int64(amount.(int32)))
	case int64:
		val = big.NewInt(amount.(int64))
	case string:
		var ok bool
		val, ok = new(big.Int).SetString(amount.(string), 0)
		if !ok {
			panic(fmt.Sprintf("invalid amount string: %s", amount.(string)))
		}
	case *big.Int:
		val = amount.(*big.Int)
	case float64:
		val = decimal.NewFromFloat(amount.(float64)).BigInt()
	default:
		panic(fmt.Sprintf("invalid amount type: %T", amount))
	}

	return val
}

func EvmContractMethodId(signature string) []byte {
	transferFnSignature := []byte(signature)
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	return hash.Sum(nil)[:4]
}
