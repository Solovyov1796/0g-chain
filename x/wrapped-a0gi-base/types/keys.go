package types

import "encoding/hex"

const (
	// ModuleName The name that will be used throughout the module
	ModuleName = "wrapped-a0gi-base"

	// StoreKey Top level store key where all module items will be stored
	StoreKey = ModuleName

	// QuerierRoute Top level query string
	QuerierRoute = ModuleName
)

var (
	// prefix
	MinterSupplyKeyPrefix = []byte{0x01}

	// keys
	WA0GIKey = []byte{0x02}
)

func GetMinterKeyFromAccount(account string) ([]byte, error) {
	return hex.DecodeString(account)
}
